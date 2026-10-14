package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
)

const remotePreviewTextMax = 256 << 10

type remotePreviewInfo struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	MIME      string `json:"mime"`
	Size      int64  `json:"size"`
	Text      string `json:"text,omitempty"`
	Truncated bool   `json:"truncated,omitempty"`
}

// The MIME allowlist is independent of OS MIME databases. Active documents are shown
// as escaped text; SVG/HTML/JavaScript never execute in WRM's authenticated origin.
func remotePreviewType(name string) (string, string) {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".mp4", ".m4v":
		return "video", "video/mp4"
	case ".webm":
		return "video", "video/webm"
	case ".ogv":
		return "video", "video/ogg"
	case ".mov":
		return "video", "video/quicktime"
	case ".mkv":
		return "video", "video/x-matroska"
	case ".mp3":
		return "audio", "audio/mpeg"
	case ".m4a", ".aac":
		if ext == ".aac" {
			return "audio", "audio/aac"
		}
		return "audio", "audio/mp4"
	case ".ogg", ".oga":
		return "audio", "audio/ogg"
	case ".wav":
		return "audio", "audio/wav"
	case ".flac":
		return "audio", "audio/flac"
	case ".png":
		return "image", "image/png"
	case ".jpg", ".jpeg":
		return "image", "image/jpeg"
	case ".gif":
		return "image", "image/gif"
	case ".webp":
		return "image", "image/webp"
	case ".avif":
		return "image", "image/avif"
	case ".bmp":
		return "image", "image/bmp"
	case ".pdf":
		return "pdf", "application/pdf"
	case ".txt", ".log", ".md", ".ini", ".conf", ".cfg", ".json", ".yaml", ".yml", ".xml", ".csv", ".tsv", ".py", ".go", ".sh", ".sql", ".js", ".ts", ".css", ".html", ".htm", ".svg", ".env":
		return "text", "text/plain; charset=utf-8"
	}
	base := strings.ToLower(path.Base(name))
	if base == ".env" || strings.HasPrefix(base, ".env.") || base == "dockerfile" || base == "makefile" || base == "readme" {
		return "text", "text/plain; charset=utf-8"
	}
	return "", ""
}

type previewWriter struct {
	http.ResponseWriter
	bytes  int64
	status int
}

func (w *previewWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}
func (w *previewWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, e := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, e
}

func previewRemoteFileHandler(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet, http.MethodHead) {
		return
	}
	c, acc, ok := fileRequest(w, r, PermFilesRead)
	if !ok {
		return
	}
	rel := r.URL.Query().Get("path")
	if rel == "" || strings.ContainsAny(rel, "\x00\n\r") {
		jsonError(w, "Invalid file path", 400)
		return
	}
	kind, mime := remotePreviewType(rel)
	if kind == "" {
		jsonError(w, "This file type has no embedded preview; download it to open it", 415)
		return
	}
	meta := r.URL.Query().Get("metadata") == "1"
	text := r.URL.Query().Get("text") == "1"
	if text && kind != "text" {
		jsonError(w, "Not a text preview", 400)
		return
	}
	w.Header().Set("Cache-Control", "private, no-store")
	info := remotePreviewInfo{Name: path.Base(rel), Kind: kind, MIME: mime}
	pw := &previewWriter{ResponseWriter: w}
	started := false
	serve := func(reader io.ReadSeeker, size int64, mtime time.Time) error {
		info.Size = size
		if meta {
			jsonOK(pw, info)
			started = true
			return nil
		}
		if text {
			b, err := io.ReadAll(io.LimitReader(reader, remotePreviewTextMax+1))
			if err != nil {
				return err
			}
			if len(b) > remotePreviewTextMax {
				b = b[:remotePreviewTextMax]
				info.Truncated = true
			}
			info.Text = strings.ToValidUTF8(string(b), "�")
			jsonOK(pw, info)
			started = true
			return nil
		}
		pw.Header().Set("Content-Type", mime)
		pw.Header().Set("Content-Disposition", "inline")
		pw.Header().Set("X-Content-Type-Options", "nosniff")
		pw.Header().Set("Cache-Control", "private, no-store")
		pw.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
		pw.Header().Set("Accept-Ranges", "bytes")
		started = true
		http.ServeContent(pw, r, info.Name, mtime, reader)
		return nil
	}
	err := fileOp(c, true, func(sc *sftp.Client) error {
		st, err := sc.Stat(rel)
		if err != nil {
			return err
		}
		if !st.Mode().IsRegular() {
			return fmt.Errorf("preview requires a regular file")
		}
		if meta {
			return serve(nil, st.Size(), st.ModTime())
		}
		file, err := sc.Open(rel)
		if err != nil {
			return err
		}
		defer file.Close()
		stop := context.AfterFunc(r.Context(), func() { file.Close() })
		defer stop()
		return serve(file, st.Size(), st.ModTime())
	}, func(fc *ftp.ServerConn) error {
		size, err := fc.FileSize(rel)
		if err != nil || size < 0 {
			return fmt.Errorf("FTP server must provide file size for streaming previews")
		}
		if meta {
			return serve(nil, size, time.Time{})
		}
		reader := &ftpPreviewReader{client: fc, name: rel, size: size}
		defer reader.Close()
		return serve(reader, size, time.Time{})
	})
	if err != nil && !started {
		jsonError(w, "Cannot preview file: "+redactText(err.Error()), 502)
		return
	}
	if started {
		status := "ok"
		if err != nil || r.Context().Err() != nil || pw.status >= 400 {
			status = "failed"
		}
		auditFileOp(r, acc, c, "file.preview", rel, map[string]interface{}{"bytes": pw.bytes, "status": status, "range": truncateStr(r.Header.Get("Range"), 100), "metadata": meta, "text": text})
	}
}

// Seek starts the next FTP transfer at REST offset; no discarded download prefix is
// needed to seek a large file. Used only by the request goroutine.
type ftpPreviewReader struct {
	client    *ftp.ServerConn
	name      string
	size, pos int64
	response  *ftp.Response
}

func (f *ftpPreviewReader) Read(b []byte) (int, error) {
	if f.pos >= f.size {
		return 0, io.EOF
	}
	if f.response == nil {
		resp, err := f.client.RetrFrom(f.name, uint64(f.pos))
		if err != nil {
			return 0, err
		}
		f.response = resp
	}
	if int64(len(b)) > f.size-f.pos {
		b = b[:f.size-f.pos]
	}
	n, err := f.response.Read(b)
	f.pos += int64(n)
	return n, err
}
func (f *ftpPreviewReader) Seek(offset int64, whence int) (int64, error) {
	p := offset
	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		p += f.pos
	case io.SeekEnd:
		p += f.size
	default:
		return 0, fmt.Errorf("invalid seek")
	}
	if p < 0 {
		return 0, fmt.Errorf("negative seek")
	}
	if p != f.pos {
		if err := f.Close(); err != nil {
			return 0, err
		}
		f.pos = p
	}
	return p, nil
}
func (f *ftpPreviewReader) Close() error {
	if f.response != nil {
		r := f.response
		f.response = nil
		return r.Close()
	}
	return nil
}
