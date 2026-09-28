package main

import (
	"archive/zip"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
)

// ─── FILE MANAGER ─────────────────────────────────────

// fileRequest validates access to the connection referenced by the "id" query parameter.
func fileRequest(w http.ResponseWriter, r *http.Request) (Connection, bool) {
	connID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if !isConnectionAccessible(r, connID) {
		jsonError(w, "Not found", 404)
		return Connection{}, false
	}
	c, err := loadConnection(connID)
	if err != nil {
		jsonError(w, "Connection not found", 404)
		return Connection{}, false
	}
	return c, true
}

// dialFTP connects and logs in; FTPS uses explicit TLS (AUTH TLS).
func dialFTP(c Connection) (*ftp.ServerConn, error) {
	opts := []ftp.DialOption{ftp.DialWithTimeout(15 * time.Second)}
	if strings.ToUpper(c.Protocol) == "FTPS" {
		opts = append(opts, ftp.DialWithExplicitTLS(&tls.Config{InsecureSkipVerify: true, ServerName: hostOnly(c.Host)}))
	}
	fc, err := ftp.Dial(c.Host, opts...)
	if err != nil {
		return nil, fmt.Errorf("FTP: %v", err)
	}
	if err := fc.Login(c.Username, c.Password); err != nil {
		fc.Quit()
		return nil, fmt.Errorf("FTP login: %v", err)
	}
	return fc, nil
}

func withFTP(c Connection, fn func(*ftp.ServerConn) error) error {
	fc, err := dialFTP(c)
	if err != nil {
		return err
	}
	defer fc.Quit()
	return fn(fc)
}

// fileOp runs fn against either an FTP or an SFTP connection.
func fileOp(c Connection, retry bool, sftpFn func(*sftp.Client) error, ftpFn func(*ftp.ServerConn) error) error {
	if isFTP(c) {
		return withFTP(c, ftpFn)
	}
	return withSFTP(c, retry, sftpFn)
}

func listRemoteFilesHandler(w http.ResponseWriter, r *http.Request) {
	c, ok := fileRequest(w, r)
	if !ok {
		return
	}
	p := r.URL.Query().Get("path")
	if p == "" {
		p = "."
	}
	list := make([]FileItem, 0)
	err := fileOp(c, true, func(sc *sftp.Client) error {
		list = list[:0]
		if p == "." || p == "~" {
			if rp, err := sc.RealPath("."); err == nil {
				p = rp
			} else if wd, err := sc.Getwd(); err == nil {
				p = wd
			}
		} else if rp, err := sc.RealPath(p); err == nil {
			p = rp // resolves "..", "." and duplicate slashes
		}
		files, err := sc.ReadDir(p)
		if err != nil {
			return fmt.Errorf("cannot read %s: %v", p, err)
		}
		for _, f := range files {
			item := FileItem{Name: f.Name(), IsDir: f.IsDir(), Size: f.Size(), Mode: f.Mode().String(), ModTime: f.ModTime().Unix()}
			if f.Mode()&os.ModeSymlink != 0 {
				item.IsLink = true
				if st, err := sc.Stat(joinRemote(p, f.Name())); err == nil {
					item.IsDir = st.IsDir()
					if !st.IsDir() {
						item.Size = st.Size()
					}
				}
			}
			list = append(list, item)
		}
		return nil
	}, func(fc *ftp.ServerConn) error {
		if p != "." && p != "" {
			if err := fc.ChangeDir(p); err != nil {
				return fmt.Errorf("cannot open %s: %v", p, err)
			}
		}
		p, _ = fc.CurrentDir()
		entries, err := fc.List(p)
		if err != nil {
			return fmt.Errorf("cannot list %s: %v", p, err)
		}
		for _, e := range entries {
			if e.Name == "." || e.Name == ".." {
				continue
			}
			list = append(list, FileItem{Name: e.Name, IsDir: e.Type == ftp.EntryTypeFolder, IsLink: e.Type == ftp.EntryTypeLink, Size: int64(e.Size), ModTime: e.Time.Unix()})
		}
		return nil
	})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]interface{}{"path": p, "files": list})
}

func listRemoteRecursiveHandler(w http.ResponseWriter, r *http.Request) {
	c, ok := fileRequest(w, r)
	if !ok {
		return
	}
	p := r.URL.Query().Get("path")
	if p == "" {
		p = "."
	}
	if isFTP(c) {
		jsonError(w, "Not supported for FTP", 400)
		return
	}
	type recursiveFile struct {
		Path string `json:"path"`
		Name string `json:"name"`
		Size int64  `json:"size"`
	}
	const maxFiles = 200000
	var files []recursiveFile
	notFound := false
	err := withSFTP(c, true, func(sc *sftp.Client) error {
		files = files[:0]
		info, err := sc.Stat(p)
		if err != nil {
			notFound = true
			return fmt.Errorf("path not found: %v", err)
		}
		if !info.IsDir() {
			files = append(files, recursiveFile{Path: p, Name: path.Base(p), Size: info.Size()})
			return nil
		}
		var walk func(dir string) error
		walk = func(dir string) error {
			entries, err := sc.ReadDir(dir)
			if err != nil {
				if isConnLostErr(err) {
					return err
				}
				return nil
			}
			for _, e := range entries {
				if len(files) >= maxFiles || r.Context().Err() != nil {
					return nil
				}
				full := joinRemote(dir, e.Name())
				if e.IsDir() {
					if err := walk(full); err != nil {
						return err
					}
				} else {
					files = append(files, recursiveFile{Path: full, Name: e.Name(), Size: e.Size()})
				}
			}
			return nil
		}
		return walk(p)
	})
	if err != nil {
		code := 500
		if notFound {
			code = 404
		}
		jsonError(w, err.Error(), code)
		return
	}
	var totalSize int64
	for _, f := range files {
		totalSize += f.Size
	}
	jsonOK(w, map[string]interface{}{"files": files, "count": len(files), "totalSize": totalSize})
}

// errStreamStarted wraps errors that happen after bytes were sent to the browser, so the
// pool never retries (which would corrupt the download).
type errStreamStarted struct{ err error }

func (e errStreamStarted) Error() string { return e.err.Error() }

func downloadRemoteFileHandler(w http.ResponseWriter, r *http.Request) {
	c, ok := fileRequest(w, r)
	if !ok {
		return
	}
	remotePath := r.URL.Query().Get("path")
	disposition := attachmentHeader(path.Base(remotePath))
	if r.URL.Query().Get("inline") == "1" {
		disposition = "inline"
	}
	started := false
	err := fileOp(c, true, func(sc *sftp.Client) error {
		info, err := sc.Stat(remotePath)
		if err != nil {
			return fmt.Errorf("file not found: %v", err)
		}
		if info.IsDir() {
			return errors.New("cannot download a directory (use ZIP download)")
		}
		file, err := sc.Open(remotePath)
		if err != nil {
			return err
		}
		defer file.Close()
		w.Header().Set("Content-Disposition", disposition)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
		started = true
		if _, err := io.Copy(w, io.LimitReader(file, info.Size())); err != nil {
			return errStreamStarted{err}
		}
		return nil
	}, func(fc *ftp.ServerConn) error {
		if size, err := fc.FileSize(remotePath); err == nil && size >= 0 {
			w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
		}
		resp, err := fc.Retr(remotePath)
		if err != nil {
			w.Header().Del("Content-Length")
			return err
		}
		defer resp.Close()
		w.Header().Set("Content-Disposition", disposition)
		w.Header().Set("Content-Type", "application/octet-stream")
		started = true
		_, err = io.Copy(w, resp)
		return err
	})
	if err != nil {
		if started {
			log.Printf("download %s aborted: %v", remotePath, err)
			return
		}
		jsonError(w, err.Error(), 500)
	}
}

func downloadRemoteDirHandler(w http.ResponseWriter, r *http.Request) {
	c, ok := fileRequest(w, r)
	if !ok {
		return
	}
	if isFTP(c) {
		jsonError(w, "ZIP download is not supported for FTP", 400)
		return
	}
	remotePath := r.URL.Query().Get("path")
	dirName := path.Base(strings.TrimSuffix(remotePath, "/"))
	if dirName == "" || dirName == "/" || dirName == "." {
		dirName = "root"
	}
	started := false
	err := withSFTP(c, true, func(sc *sftp.Client) error {
		info, err := sc.Stat(remotePath)
		if err != nil {
			return fmt.Errorf("folder not found: %v", err)
		}
		if !info.IsDir() {
			return errors.New("not a directory")
		}
		w.Header().Set("Content-Disposition", attachmentHeader(dirName+".zip"))
		w.Header().Set("Content-Type", "application/zip")
		started = true
		zw := zip.NewWriter(w)
		if err := zipRemoteDir(r, sc, zw, remotePath, dirName); err != nil {
			zw.Close()
			return errStreamStarted{err}
		}
		return zw.Close()
	})
	if err != nil {
		if started {
			log.Printf("zip %s aborted: %v", remotePath, err)
			return
		}
		jsonError(w, err.Error(), 500)
	}
}

func zipRemoteDir(r *http.Request, sc *sftp.Client, zw *zip.Writer, remotePath, baseName string) error {
	entries, err := sc.ReadDir(remotePath)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		zw.Create(baseName + "/")
	}
	for _, entry := range entries {
		if err := r.Context().Err(); err != nil {
			return err
		}
		fullRemote := joinRemote(remotePath, entry.Name())
		zipPath := baseName + "/" + entry.Name()
		if entry.IsDir() {
			if err := zipRemoteDir(r, sc, zw, fullRemote, zipPath); err != nil {
				if r.Context().Err() != nil || isConnLostErr(err) {
					return err
				}
				log.Printf("zip subdir %s: %v", fullRemote, err)
			}
			continue
		}
		if !entry.Mode().IsRegular() {
			continue
		}
		f, err := sc.Open(fullRemote)
		if err != nil {
			continue
		}
		hdr := &zip.FileHeader{Name: zipPath, Method: zip.Deflate}
		hdr.Modified = entry.ModTime()
		fw, err := zw.CreateHeader(hdr)
		if err != nil {
			f.Close()
			return err
		}
		_, err = io.Copy(fw, f)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// cleanRelPath validates a client-supplied relative path (for folder uploads).
func cleanRelPath(rel string) (string, error) {
	rel = strings.ReplaceAll(rel, "\\", "/")
	rel = strings.Trim(rel, "/")
	if rel == "" {
		return "", nil
	}
	for _, part := range strings.Split(rel, "/") {
		if part == ".." || part == "." || part == "" {
			return "", fmt.Errorf("invalid relative path %q", rel)
		}
	}
	return rel, nil
}

type uploadResult struct {
	Name   string `json:"name"`
	Status string `json:"status"` // ok | conflict | error
	Size   int64  `json:"size"`
	Error  string `json:"error,omitempty"`
}

// uploadRemoteFileHandler streams multipart uploads straight to the remote server
// (no buffering of the whole file in RAM). Parameters can come from the query string or
// from form fields that precede the file parts. A "relpath" field before a file part
// uploads that file into a sub-directory (folder upload / drag & drop of folders).
func uploadRemoteFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	mr, err := r.MultipartReader()
	if err != nil {
		jsonError(w, "Upload error: "+err.Error(), 400)
		return
	}
	q := r.URL.Query()
	params := map[string]string{"id": q.Get("id"), "path": q.Get("path"), "overwrite": q.Get("overwrite"), "share_token": q.Get("share_token")}
	var (
		c        Connection
		haveConn bool
		results  []uploadResult
		relDir   string
	)
	ensureConn := func() bool {
		if haveConn {
			return true
		}
		connID, _ := strconv.Atoi(params["id"])
		if params["share_token"] != "" && q.Get("share_token") == "" {
			qq := r.URL.Query()
			qq.Set("share_token", params["share_token"])
			r.URL.RawQuery = qq.Encode()
		}
		if !isConnectionAccessible(r, connID) {
			jsonError(w, "Not found", 404)
			return false
		}
		var err error
		if c, err = loadConnection(connID); err != nil {
			jsonError(w, "Connection not found", 404)
			return false
		}
		haveConn = true
		return true
	}
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			if !haveConn {
				jsonError(w, "Upload error: "+err.Error(), 400)
				return
			}
			results = append(results, uploadResult{Status: "error", Error: err.Error()})
			break
		}
		name := part.FormName()
		if part.FileName() == "" {
			val, _ := io.ReadAll(io.LimitReader(part, 64*1024))
			if name == "relpath" {
				relDir = string(val)
			} else {
				params[name] = string(val)
			}
			continue
		}
		if !ensureConn() {
			return
		}
		fileName := path.Base(strings.ReplaceAll(part.FileName(), "\\", "/"))
		rel, err := cleanRelPath(relDir)
		relDir = ""
		res := uploadResult{Name: fileName}
		if err != nil {
			res.Status, res.Error = "error", err.Error()
			results = append(results, res)
			io.Copy(io.Discard, part)
			continue
		}
		baseDir := params["path"]
		if baseDir == "" {
			baseDir = "."
		}
		targetDir := baseDir
		if rel != "" {
			dir := path.Dir(rel)
			if dir != "." {
				targetDir = joinRemote(baseDir, dir)
			}
			res.Name = rel
		}
		fullPath := joinRemote(targetDir, fileName)
		overwrite := params["overwrite"] == "true" || params["overwrite"] == "1"
		cr := &countingReader{r: part}
		opErr := fileOp(c, false, func(sc *sftp.Client) error {
			if targetDir != baseDir {
				if err := sc.MkdirAll(targetDir); err != nil {
					return fmt.Errorf("mkdir %s: %v", targetDir, err)
				}
			}
			if !overwrite {
				if _, err := sc.Stat(fullPath); err == nil {
					return errConflict
				}
			}
			dst, err := sc.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
			if err != nil {
				return err
			}
			if _, err := dst.ReadFrom(cr); err != nil {
				dst.Close()
				return err
			}
			return dst.Close()
		}, func(fc *ftp.ServerConn) error {
			if targetDir != baseDir {
				ftpMkdirAll(fc, targetDir)
			}
			if !overwrite {
				if _, err := fc.FileSize(fullPath); err == nil {
					return errConflict
				}
			}
			return fc.Stor(fullPath, cr)
		})
		res.Size = cr.n
		switch {
		case opErr == nil:
			res.Status = "ok"
		case errors.Is(opErr, errConflict):
			res.Status = "conflict"
		default:
			res.Status, res.Error = "error", opErr.Error()
		}
		io.Copy(io.Discard, part) // drain whatever was not consumed (conflict/error)
		results = append(results, res)
	}
	if !haveConn {
		jsonError(w, "No file received", 400)
		return
	}
	status := http.StatusOK
	allConflict := len(results) > 0
	for _, res := range results {
		if res.Status != "conflict" {
			allConflict = false
		}
	}
	if allConflict {
		status = http.StatusConflict
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body := map[string]interface{}{"results": results}
	if len(results) == 1 && results[0].Status == "conflict" {
		body["conflict"] = results[0].Name
	}
	w.Write(jsonMarshal(body))
}

var errConflict = errors.New("file exists")

type countingReader struct {
	r io.Reader
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}

func ftpMkdirAll(fc *ftp.ServerConn, dir string) {
	cur := ""
	if strings.HasPrefix(dir, "/") {
		cur = "/"
	}
	for _, part := range strings.Split(strings.Trim(dir, "/"), "/") {
		if part == "" {
			continue
		}
		if cur == "" || cur == "/" {
			cur += part
		} else {
			cur += "/" + part
		}
		fc.MakeDir(cur) // errors ignored: the directory usually already exists
	}
}

func mkdirRemoteHandler(w http.ResponseWriter, r *http.Request) {
	c, ok := fileRequest(w, r)
	if !ok {
		return
	}
	remotePath := r.URL.Query().Get("path")
	err := fileOp(c, true, func(sc *sftp.Client) error {
		return sc.MkdirAll(remotePath)
	}, func(fc *ftp.ServerConn) error {
		return fc.MakeDir(remotePath)
	})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]interface{}{"ok": true, "status": "ok"})
}

func deleteRemoteHandler(w http.ResponseWriter, r *http.Request) {
	c, ok := fileRequest(w, r)
	if !ok {
		return
	}
	remotePath := r.URL.Query().Get("path")
	isDir := r.URL.Query().Get("dir") == "1"
	if p := strings.TrimRight(remotePath, "/"); p == "" || p == "." || p == ".." {
		jsonError(w, "Refusing to delete this path", 400)
		return
	}
	err := fileOp(c, true, func(sc *sftp.Client) error {
		if isDir {
			return removeRemoteDirSFTP(sc, remotePath)
		}
		return sc.Remove(remotePath)
	}, func(fc *ftp.ServerConn) error {
		if isDir {
			return fc.RemoveDirRecur(remotePath)
		}
		return fc.Delete(remotePath)
	})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]interface{}{"ok": true, "status": "ok"})
}

func renameRemoteHandler(w http.ResponseWriter, r *http.Request) {
	oldPath := r.URL.Query().Get("old")
	newPath := r.URL.Query().Get("new")
	if oldPath == "" || newPath == "" {
		jsonError(w, "Missing old/new path", 400)
		return
	}
	c, ok := fileRequest(w, r)
	if !ok {
		return
	}
	err := fileOp(c, true, func(sc *sftp.Client) error {
		if _, err := sc.Lstat(newPath); err == nil {
			return fmt.Errorf("%s already exists", path.Base(newPath))
		}
		if err := sc.PosixRename(oldPath, newPath); err != nil {
			return sc.Rename(oldPath, newPath)
		}
		return nil
	}, func(fc *ftp.ServerConn) error {
		return fc.Rename(oldPath, newPath)
	})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]interface{}{"ok": true, "status": "ok"})
}

func removeRemoteDirSFTP(sc *sftp.Client, p string) error {
	entries, err := sc.ReadDir(p)
	if err != nil {
		return err
	}
	for _, e := range entries {
		full := joinRemote(p, e.Name())
		if e.IsDir() {
			if err := removeRemoteDirSFTP(sc, full); err != nil {
				return err
			}
		} else if err := sc.Remove(full); err != nil {
			return err
		}
	}
	return sc.RemoveDirectory(p)
}
