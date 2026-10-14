package main

import (
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestNestedFoldersOwnershipCycleAndDelete(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()
	u := newTestUser(t, srv, "nested-folders", false)
	other := newTestUser(t, srv, "nested-other", false)
	create := func(name string, parent *int) Folder {
		var f Folder
		if code := u.jsonDo("POST", "/api/folders", map[string]interface{}{"name": name, "parent_id": parent}, &f); code != 201 {
			t.Fatalf("create: %d", code)
		}
		return f
	}
	a := create("Client-A", nil)
	b := create("Site-1", &a.ID)
	c := create("Servers", &b.ID)
	if c.Path != "Client-A / Site-1 / Servers" {
		t.Fatalf("path: %q", c.Path)
	}
	if code := u.jsonDo("PUT", "/api/folders/"+strconv.Itoa(a.ID), map[string]interface{}{"name": a.Name, "parent_id": c.ID}, nil); code != 400 {
		t.Fatalf("cycle accepted: %d", code)
	}
	if code := other.jsonDo("POST", "/api/folders", map[string]interface{}{"name": "foreign", "parent_id": a.ID}, nil); code != 400 {
		t.Fatalf("foreign parent accepted: %d", code)
	}
	var updated Folder
	if code := u.jsonDo("PUT", "/api/folders/"+strconv.Itoa(c.ID), map[string]interface{}{"name": "Renamed"}, &updated); code != 200 || updated.ParentID == nil || *updated.ParentID != b.ID {
		t.Fatalf("old client lost parent: %d %+v", code, updated)
	}
	if code := u.jsonDo("DELETE", "/api/folders/"+strconv.Itoa(b.ID), nil, nil); code != 200 {
		t.Fatalf("delete: %d", code)
	}
	for _, f := range loadFolders(u.userID) {
		if f.ID == c.ID && (f.ParentID == nil || *f.ParentID != a.ID || f.Path != "Client-A / Renamed") {
			t.Fatalf("child not reparented: %+v", f)
		}
	}
}

func TestNestedFolderImportRoundTripAndDepth(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()
	u := newTestUser(t, srv, "folder-roundtrip", false)
	list := []Folder{{ID: 4, Name: "Site-1", ParentID: folderTestPtr(2)}, {ID: 3, Name: "Site-1", ParentID: folderTestPtr(1)}, {ID: 2, Name: "Client-B"}, {ID: 1, Name: "Client-A"}}
	conns := []Connection{{ID: 10, Name: "Backup", Protocol: "SSH", Host: "example.com", FolderID: folderTestPtr(3)}, {ID: 11, Name: "Backup", Protocol: "SSH", Host: "example.net", FolderID: folderTestPtr(4)}}
	if code := u.jsonDo("POST", "/api/config/import", map[string]interface{}{"folders": list, "connections": conns}, nil); code != 200 {
		t.Fatalf("import: %d", code)
	}
	if fs := loadFolders(u.userID); len(fs) != 4 {
		t.Fatalf("same-named children merged: %+v", fs)
	}
	var exported struct {
		Folders     []Folder     `json:"folders"`
		Connections []Connection `json:"connections"`
	}
	if code := u.jsonDo("GET", "/api/config/export", nil, &exported); code != 200 {
		t.Fatal(code)
	}
	v := newTestUser(t, srv, "folder-roundtrip-copy", false)
	if code := v.jsonDo("POST", "/api/config/import", exported, nil); code != 200 {
		t.Fatal(code)
	}
	paths := map[string]bool{}
	for _, c := range loadUserConnections(v.userID) {
		paths[folderPathOf(v.userID, c.FolderID)] = true
	}
	if !paths["Client-A / Site-1"] || !paths["Client-B / Site-1"] {
		t.Fatalf("hierarchy lost: %v", paths)
	}
	list = []Folder{}
	for i := 1; i <= folderMaxDepth+1; i++ {
		f := Folder{ID: i, Name: "folder"}
		if i > 1 {
			f.ParentID = folderTestPtr(i - 1)
		}
		list = append(list, f)
	}
	if _, err := folderImportOrder(list); err == nil {
		t.Fatal("excess depth accepted")
	}
	list = []Folder{{ID: 1, Name: "a", ParentID: folderTestPtr(2)}, {ID: 2, Name: "b", ParentID: folderTestPtr(1)}}
	if _, err := folderImportOrder(list); err == nil {
		t.Fatal("import cycle accepted")
	}
}

func folderTestPtr(n int) *int { return &n }
