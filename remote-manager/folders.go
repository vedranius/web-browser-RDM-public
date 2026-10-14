package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const folderMaxDepth = 32

// Serialize hierarchy changes so concurrent moves cannot introduce a cycle.
var folderHierarchyMu sync.Mutex

func folderPath(list []Folder, id int) string {
	byID := map[int]Folder{}
	for _, f := range list {
		byID[f.ID] = f
	}
	parts := []string{}
	seen := map[int]bool{}
	for id > 0 && !seen[id] && len(parts) < folderMaxDepth {
		f, ok := byID[id]
		if !ok {
			break
		}
		seen[id] = true
		parts = append([]string{f.Name}, parts...)
		id = 0
		if f.ParentID != nil {
			id = *f.ParentID
		}
	}
	return strings.Join(parts, " / ")
}

func folderPathOf(userID int, id *int) string {
	if id == nil {
		return ""
	}
	return folderPath(loadFolders(userID), *id)
}

func validateFolderParent(userID int, f *Folder) error {
	if f.ParentID != nil && *f.ParentID <= 0 {
		f.ParentID = nil
	}
	list := loadFolders(userID)
	byID := map[int]Folder{}
	for _, x := range list {
		byID[x.ID] = x
	}
	if f.ParentID != nil {
		if _, ok := byID[*f.ParentID]; !ok {
			return fmt.Errorf("parent folder not found")
		}
	}
	// Apply the proposed parent before checking every descendant's depth.
	byID[f.ID] = *f
	for start := range byID {
		seen := map[int]bool{}
		id := start
		for depth := 0; ; depth++ {
			if depth >= folderMaxDepth {
				return fmt.Errorf("at most %d nested folders", folderMaxDepth)
			}
			if seen[id] {
				return fmt.Errorf("a folder cannot contain itself or one of its ancestors")
			}
			seen[id] = true
			x, ok := byID[id]
			if !ok || x.ParentID == nil {
				break
			}
			id = *x.ParentID
		}
	}
	return nil
}

func folderImportKey(f Folder) string {
	if f.ID > 0 {
		return fmt.Sprintf("id:%d", f.ID)
	}
	return "name:" + strings.TrimSpace(f.Name)
}

// Order parents before children, regardless of the export's alphabetical order.
func folderImportOrder(list []Folder) ([]int, error) {
	byID := map[int]int{}
	for i, f := range list {
		if f.ID > 0 {
			if _, dup := byID[f.ID]; dup {
				return nil, fmt.Errorf("duplicate folder id")
			}
			byID[f.ID] = i
		}
	}
	state := map[int]int{}
	order := []int{}
	var visit func(int, int) error
	visit = func(i, depth int) error {
		if depth >= folderMaxDepth {
			return fmt.Errorf("at most %d nested folders", folderMaxDepth)
		}
		if state[i] == 1 {
			return fmt.Errorf("folder hierarchy contains a cycle")
		}
		if state[i] == 2 {
			return nil
		}
		state[i] = 1
		if p := list[i].ParentID; p != nil && *p > 0 {
			j, ok := byID[*p]
			if !ok || strings.TrimSpace(list[j].Name) == "" {
				return fmt.Errorf("parent folder not found in import")
			}
			if err := visit(j, depth+1); err != nil {
				return err
			}
		}
		state[i] = 2
		order = append(order, i)
		return nil
	}
	for i := range list {
		if err := visit(i, 0); err != nil {
			return nil, err
		}
	}
	// Cached visits alone cannot measure a child of a previously visited deep parent.
	for i := range list {
		depth := 1
		for p := list[i].ParentID; p != nil && *p > 0; {
			depth++
			if depth > folderMaxDepth {
				return nil, fmt.Errorf("at most %d nested folders", folderMaxDepth)
			}
			p = list[byID[*p]].ParentID
		}
	}
	return order, nil
}

// Reorder siblings or move a folder under a new parent in one transaction.
func apiFolderReorderHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireAuth(w, r)
	if !ok {
		return
	}
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var in struct {
		ID       int  `json:"id"`
		ParentID *int `json:"parent_id"`
		BeforeID int  `json:"before_id"`
		AfterID  int  `json:"after_id"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil || in.ID <= 0 || (in.BeforeID != 0 && in.AfterID != 0) {
		jsonError(w, "Invalid folder order", 400)
		return
	}
	folderHierarchyMu.Lock()
	defer folderHierarchyMu.Unlock()
	list := loadFolders(uid)
	var moved Folder
	found := false
	for _, f := range list {
		if f.ID == in.ID {
			moved = f
			found = true
			break
		}
	}
	if !found {
		jsonError(w, "Not found", 404)
		return
	}
	moved.ParentID = in.ParentID
	if e := validateFolderParent(uid, &moved); e != nil {
		jsonError(w, e.Error(), 400)
		return
	}
	siblings := []int{}
	for _, f := range list {
		if f.ID != in.ID && intPtrEq(f.ParentID, moved.ParentID) {
			siblings = append(siblings, f.ID)
		}
	}
	index := len(siblings)
	anchor := in.BeforeID
	if anchor == 0 {
		anchor = in.AfterID
	}
	if anchor != 0 {
		index = -1
		for i, id := range siblings {
			if id == anchor {
				index = i
				if in.AfterID != 0 {
					index++
				}
				break
			}
		}
		if index < 0 {
			jsonError(w, "Order anchor is not a sibling owned by you", 400)
			return
		}
	}
	siblings = append(siblings, 0)
	copy(siblings[index+1:], siblings[index:])
	siblings[index] = in.ID
	tx, e := db.Begin()
	if e != nil {
		jsonError(w, "Cannot change folder order", 500)
		return
	}
	defer tx.Rollback()
	if _, e = tx.Exec("UPDATE folders SET parent_id=? WHERE id=? AND user_id=?", moved.ParentID, in.ID, uid); e == nil {
		for i, id := range siblings {
			if _, e = tx.Exec("UPDATE folders SET sort_order=? WHERE id=? AND user_id=?", i+1, id, uid); e != nil {
				break
			}
		}
	}
	if e != nil {
		jsonError(w, "Cannot change folder order", 500)
		return
	}
	if e = tx.Commit(); e != nil {
		jsonError(w, "Cannot change folder order", 500)
		return
	}
	auditLog(r, uid, "folder.reordered", moved.Name, map[string]interface{}{"folder_id": in.ID, "parent_id": moved.ParentID})
	jsonOK(w, loadFolders(uid))
}
