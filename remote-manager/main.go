package main

import (
        "crypto/aes"
        "crypto/cipher"
        "crypto/rand"
        "database/sql"
        "embed"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "os"
        "path"
        "path/filepath"
        "strconv"
        "strings"
        "sync"
        "sync/atomic"
        "time"

        "github.com/gorilla/websocket"
        "github.com/pkg/sftp"
        "golang.org/x/crypto/bcrypt"
        "golang.org/x/crypto/ssh"
        _ "modernc.org/sqlite"
)

//go:embed static/*
var staticFiles embed.FS

// AppVersion can be overridden at build time with -ldflags "-X main.AppVersion=..."
var AppVersion = "v9.10.2-mimo"
const sessionCookieName = "wrm_session"
const defaultEncryptionKey = "change-this-to-a-secure-32-char!"

var upgrader = websocket.Upgrader{CheckOrigin: checkWSOrigin, ReadBufferSize: 32 * 1024, WriteBufferSize: 64 * 1024}
var db *sql.DB
var encryptionKey []byte

// ─── TYPES ───────────────────────────────────────────

type Connection struct {
        ID         int    `json:"id"`
        Name       string `json:"name"`
        Protocol   string `json:"protocol"`
        Host       string `json:"host"`
        Username   string `json:"username"`
        AuthMethod string `json:"auth_method"`
        Password   string `json:"password"`
        PrivateKey string `json:"private_key"`
        KeyPath    string `json:"key_path"`
        FolderID   *int   `json:"folder_id"`
}

type Folder struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
}

type FileItem struct {
        Name    string `json:"name"`
        Size    int64  `json:"size"`
        IsDir   bool   `json:"is_dir"`
        IsLink  bool   `json:"is_link,omitempty"`
        Mode    string `json:"mode"`
        ModTime int64  `json:"mtime"`
}

type ResizeMsg struct {
        Type string `json:"type"`
        Cols uint32 `json:"cols"`
        Rows uint32 `json:"rows"`
}

type Session struct {
        ID            int    `json:"id"`
        Name          string `json:"name"`
        CreatedAt     string `json:"created_at"`
        Locked        bool   `json:"locked"`
        Pane1ConnID   *int   `json:"pane1_conn_id"`
        Pane1Mode     string `json:"pane1_mode"`
        Pane2ConnID   *int   `json:"pane2_conn_id"`
        Pane2Mode     string `json:"pane2_mode"`
        Pane1ConnName string `json:"pane1_conn_name,omitempty"`
        Pane2ConnName string `json:"pane2_conn_name,omitempty"`
        // Layout is a JSON document describing the open windows (connection, mode, path,
        // position/snap) so a saved session can be restored with one click.
        Layout string `json:"layout"`
}

const maxSessionLayout = 256 * 1024

type User struct {
        ID        int    `json:"id"`
        Username  string `json:"username"`
        CreatedAt string `json:"created_at"`
}

type ShareLink struct {
        ID           int    `json:"id"`
        OwnerID      int    `json:"owner_id"`
        Token        string `json:"token"`
        Name         string `json:"name"`
        PasswordHash string `json:"-"`
        Public       bool   `json:"public"`
        Active       bool   `json:"active"`
        CreatedAt    string `json:"created_at"`
}

// ─── BROADCAST HUB ───────────────────────────────────

type Hub struct {
        mu      sync.Mutex
        clients map[chan []byte]struct{}
}

var hub = &Hub{clients: make(map[chan []byte]struct{})}

func (h *Hub) subscribe() chan []byte {
        ch := make(chan []byte, 32)
        h.mu.Lock()
        h.clients[ch] = struct{}{}
        h.mu.Unlock()
        return ch
}
func (h *Hub) unsubscribe(ch chan []byte) {
        h.mu.Lock()
        delete(h.clients, ch)
        close(ch)
        h.mu.Unlock()
}
func (h *Hub) broadcast(msg []byte) {
        h.mu.Lock()
        defer h.mu.Unlock()
        for ch := range h.clients {
                select {
                case ch <- msg:
                default:
                }
        }
}
func broadcastSessionUpdate() { hub.broadcast([]byte(`{"type":"sessions_changed"}`)) }

// ─── SHARE COLLAB HUB ─────────────────────────────────

type ShareRoom struct {
        mu         sync.Mutex
        clients    map[*websocket.Conn]string
        nameToConn map[string]*websocket.Conn
        // sharers are clients currently sharing a terminal; only they may grant control.
        sharers map[*websocket.Conn]bool
        // controllers maps a client with keyboard control to the sharer who granted it.
        // Remote input is delivered to that sharer only.
        controllers map[*websocket.Conn]*websocket.Conn
}

var (
        shareRooms   = make(map[string]*ShareRoom)
        shareRoomsMu sync.Mutex
)

func getShareRoom(token string) *ShareRoom {
        shareRoomsMu.Lock()
        defer shareRoomsMu.Unlock()
        room, ok := shareRooms[token]
        if !ok {
                room = &ShareRoom{
                        clients:     make(map[*websocket.Conn]string),
                        nameToConn:  make(map[string]*websocket.Conn),
                        sharers:     make(map[*websocket.Conn]bool),
                        controllers: make(map[*websocket.Conn]*websocket.Conn),
                }
                shareRooms[token] = room
        }
        return room
}

// addClient registers a participant and returns its unique display name (a suffix is
// added when the name is already taken, e.g. the same user in two browser tabs).
func (r *ShareRoom) addClient(conn *websocket.Conn, name string) string {
        r.mu.Lock()
        defer r.mu.Unlock()
        base := name
        for i := 2; r.nameToConn[name] != nil; i++ {
                name = fmt.Sprintf("%s (%d)", base, i)
        }
        r.clients[conn] = name
        r.nameToConn[name] = conn
        r.broadcastPresenceLocked()
        return name
}

func (r *ShareRoom) removeClient(conn *websocket.Conn) {
        r.mu.Lock()
        defer r.mu.Unlock()
        if n, ok := r.clients[conn]; ok {
                delete(r.nameToConn, n)
        }
        delete(r.clients, conn)
        delete(r.controllers, conn)
        delete(r.sharers, conn)
        r.revokeGrantedByLocked(conn)
        r.broadcastPresenceLocked()
}

// revokeGrantedByLocked removes every keyboard control granted by the given sharer
// (it stopped sharing or left). Must be called with r.mu held.
func (r *ShareRoom) revokeGrantedByLocked(granter *websocket.Conn) {
        by := r.clients[granter]
        for c, g := range r.controllers {
                if g != granter {
                        continue
                }
                delete(r.controllers, c)
                writeLocked(c, jsonMarshal(map[string]interface{}{"type": "control-status", "granted": false}))
                msg := jsonMarshal(map[string]interface{}{"type": "control-revoked", "user": r.clients[c], "by": by})
                for other := range r.clients {
                        writeLocked(other, msg)
                }
        }
}

func removeShareRoomIfEmpty(token string) {
        shareRoomsMu.Lock()
        defer shareRoomsMu.Unlock()
        if room, ok := shareRooms[token]; ok {
                room.mu.Lock()
                empty := len(room.clients) == 0
                room.mu.Unlock()
                if empty {
                        delete(shareRooms, token)
                }
        }
}

func (r *ShareRoom) listUsers() []string {
        r.mu.Lock()
        defer r.mu.Unlock()
        users := make([]string, 0, len(r.clients))
        for _, name := range r.clients {
                users = append(users, name)
        }
        return users
}

// writeLocked must be called with r.mu held: gorilla/websocket connections allow only
// one concurrent writer, and every write to a room member goes through the room lock.
func writeLocked(conn *websocket.Conn, msg []byte) {
        conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
        conn.WriteMessage(websocket.TextMessage, msg)
}

func (r *ShareRoom) broadcast(msg []byte) {
        r.mu.Lock()
        defer r.mu.Unlock()
        for conn := range r.clients {
                writeLocked(conn, msg)
        }
}

func (r *ShareRoom) broadcastExcept(skip *websocket.Conn, msg []byte) {
        r.mu.Lock()
        defer r.mu.Unlock()
        for conn := range r.clients {
                if conn != skip {
                        writeLocked(conn, msg)
                }
        }
}

func (r *ShareRoom) sendTo(conn *websocket.Conn, msg []byte) {
        r.mu.Lock()
        defer r.mu.Unlock()
        if _, ok := r.clients[conn]; ok {
                writeLocked(conn, msg)
        }
}

func (r *ShareRoom) broadcastPresenceLocked() {
        users := make([]string, 0, len(r.clients))
        for _, name := range r.clients {
                users = append(users, name)
        }
        payload, _ := json.Marshal(map[string]interface{}{"type": "presence", "users": users})
        for conn := range r.clients {
                writeLocked(conn, payload)
        }
}

func shareWSHandler(w http.ResponseWriter, r *http.Request) {
        tokenPath := strings.TrimPrefix(r.URL.Path, "/ws/share/")
        parts := strings.SplitN(strings.TrimLeft(tokenPath, "/"), "/", 2)
        token := strings.TrimSpace(parts[0])
        if token == "" {
                http.Error(w, "Bad share token", 400)
                return
        }
        share, _, err := getShareByToken(token)
        if err != nil || !share.Active {
                http.Error(w, "Share not found", 404)
                return
        }
        if !hasShareAccess(r, share) {
                http.Error(w, "Unauthorized", 401)
                return
        }
        ws, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                return
        }
        defer ws.Close()
        ws.SetReadLimit(8 << 20)
        name := strings.TrimSpace(r.URL.Query().Get("name"))
        if len(name) > 40 {
                name = name[:40]
        }
        if name == "" {
                if userID, err := currentUserID(r); err == nil {
                        db.QueryRow("SELECT username FROM users WHERE id=?", userID).Scan(&name)
                }
        }
        if name == "" {
                b := make([]byte, 3)
                rand.Read(b)
                name = fmt.Sprintf("Guest-%x", b)
        }
        room := getShareRoom(token)
        name = room.addClient(ws, name)
        room.sendTo(ws, jsonMarshal(map[string]interface{}{"type": "welcome", "name": name}))
        // Tell late joiners who is already sharing a terminal, so they can watch it.
        room.mu.Lock()
        for c := range room.sharers {
                writeLocked(ws, jsonMarshal(map[string]interface{}{"type": "screen", "user": room.clients[c], "action": "start"}))
        }
        room.mu.Unlock()
        defer removeShareRoomIfEmpty(token)
        defer room.removeClient(ws)
        room.broadcast(jsonMarshal(map[string]interface{}{
                "type":      "system",
                "user":      "SYSTEM",
                "text":      name + " joined the collaboration",
                "timestamp": time.Now().Format(time.RFC3339),
        }))
        for {
                _, msg, err := ws.ReadMessage()
                if err != nil {
                        break
                }
                var payload map[string]interface{}
                if err := json.Unmarshal(msg, &payload); err != nil {
                        continue
                }
                switch payload["type"] {
                case "chat":
                        text, _ := payload["text"].(string)
                        if strings.TrimSpace(text) == "" {
                                continue
                        }
                        room.broadcast(jsonMarshal(map[string]interface{}{
                                "type":      "chat",
                                "user":      name,
                                "text":      text,
                                "timestamp": time.Now().Format(time.RFC3339),
                        }))
                case "file":
                        dataStr, _ := payload["data"].(string)
                        if len(dataStr) > 5*1024*1024 {
                                continue
                        }
                        room.broadcastExcept(ws, jsonMarshal(map[string]interface{}{
                                "type":      "file",
                                "user":      name,
                                "name":      payload["name"],
                                "data":      dataStr,
                                "mime":      payload["mime"],
                                "timestamp": time.Now().Format(time.RFC3339),
                        }))
                case "screen":
                        action, _ := payload["action"].(string)
                        room.mu.Lock()
                        allowed := true
                        switch action {
                        case "start":
                                room.sharers[ws] = true
                        case "stop":
                                delete(room.sharers, ws)
                                room.revokeGrantedByLocked(ws)
                        default:
                                allowed = room.sharers[ws] // only active sharers may stream data
                        }
                        room.mu.Unlock()
                        if !allowed {
                                continue
                        }
                        room.broadcastExcept(ws, jsonMarshal(map[string]interface{}{
                                "type":   "screen",
                                "user":   name,
                                "action": payload["action"],
                                "data":   payload["data"],
                        }))
                case "grant-control":
                        // Only a participant who is sharing a terminal can hand out control of it.
                        targetUser, _ := payload["user"].(string)
                        room.mu.Lock()
                        targetConn := room.nameToConn[targetUser]
                        if !room.sharers[ws] || targetConn == ws {
                                targetConn = nil
                        }
                        if targetConn != nil {
                                room.controllers[targetConn] = ws
                        }
                        room.mu.Unlock()
                        if targetConn != nil {
                                room.sendTo(targetConn, jsonMarshal(map[string]interface{}{
                                        "type":    "control-status",
                                        "granted": true,
                                        "by":      name,
                                }))
                                room.broadcast(jsonMarshal(map[string]interface{}{
                                        "type": "control-granted",
                                        "user": targetUser,
                                        "by":   name,
                                }))
                        }
                case "revoke-control":
                        targetUser, _ := payload["user"].(string)
                        room.mu.Lock()
                        targetConn := room.nameToConn[targetUser]
                        if targetConn != nil && room.controllers[targetConn] == ws {
                                delete(room.controllers, targetConn)
                        } else {
                                targetConn = nil // only the granting sharer can revoke
                        }
                        room.mu.Unlock()
                        if targetConn != nil {
                                room.sendTo(targetConn, jsonMarshal(map[string]interface{}{
                                        "type":    "control-status",
                                        "granted": false,
                                }))
                                room.broadcast(jsonMarshal(map[string]interface{}{
                                        "type": "control-revoked",
                                        "user": targetUser,
                                        "by":   name,
                                }))
                        }
                case "remote-input":
                        // Keystrokes go only to the sharer who granted control, and only while
                        // that sharer is still sharing.
                        data, _ := payload["data"].(string)
                        room.mu.Lock()
                        if granter := room.controllers[ws]; granter != nil && room.sharers[granter] && len(data) <= 64*1024 {
                                writeLocked(granter, jsonMarshal(map[string]interface{}{
                                        "type": "remote-input",
                                        "user": name,
                                        "data": data,
                                }))
                        }
                        room.mu.Unlock()
                }
        }
        room.broadcast(jsonMarshal(map[string]interface{}{
                "type":      "system",
                "user":      "SYSTEM",
                "text":      name + " left the collaboration",
                "timestamp": time.Now().Format(time.RFC3339),
        }))
}

func jsonMarshal(v interface{}) []byte {
        b, _ := json.Marshal(v)
        return b
}

// ─── MAIN ────────────────────────────────────────────

func main() {
        initDB()
        log.Printf("Web Remote Manager %s starting", AppVersion)

        mux := http.NewServeMux()
        mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))
        mux.HandleFunc("/", rootHandler)
        mux.HandleFunc("/share/", sharePageHandler)
        mux.HandleFunc("/api/version", versionHandler)
        mux.HandleFunc("/api/auth/register", apiAuthRegisterHandler)
        mux.HandleFunc("/api/auth/login", apiAuthLoginHandler)
        mux.HandleFunc("/api/auth/logout", apiAuthLogoutHandler)
        mux.HandleFunc("/api/auth/me", apiAuthMeHandler)
        mux.HandleFunc("/api/auth/keepalive", apiAuthKeepaliveHandler)
        mux.HandleFunc("/api/connections", apiConnectionsHandler)
        mux.HandleFunc("/api/connections/bulk", apiConnectionsBulkHandler)
        mux.HandleFunc("/api/connections/test", apiConnectionTestHandler)
        mux.HandleFunc("/api/connections/", apiConnectionByIDHandler)
        mux.HandleFunc("/api/folders", apiFoldersHandler)
        mux.HandleFunc("/api/folders/", apiFolderByIDHandler)
        mux.HandleFunc("/api/sessions", apiSessionsHandler)
        mux.HandleFunc("/api/sessions/", apiSessionByIDHandler)
        mux.HandleFunc("/api/shares", apiSharesHandler)
        mux.HandleFunc("/api/shares/", apiSharesByIDHandler)
        mux.HandleFunc("/api/share/", apiShareByTokenHandler)
        mux.HandleFunc("/api/shared", apiSharedHandler)
        mux.HandleFunc("/api/users", apiUsersHandler)
        mux.HandleFunc("/api/admin/users", apiAdminUsersHandler)
        mux.HandleFunc("/api/admin/users/", apiAdminUserByIDHandler)
        mux.HandleFunc("/api/config/export", apiExportHandler)
        mux.HandleFunc("/api/config/import", apiImportHandler)
        mux.HandleFunc("/api/remote/list", listRemoteFilesHandler)
        mux.HandleFunc("/api/remote/list-recursive", listRemoteRecursiveHandler)
        mux.HandleFunc("/api/remote/download", downloadRemoteFileHandler)
        mux.HandleFunc("/api/remote/download-dir", downloadRemoteDirHandler)
        mux.HandleFunc("/api/remote/upload", uploadRemoteFileHandler)
        mux.HandleFunc("/api/remote/mkdir", mkdirRemoteHandler)
        mux.HandleFunc("/api/remote/delete", deleteRemoteHandler)
        mux.HandleFunc("/api/remote/rename", renameRemoteHandler)
        mux.HandleFunc("/api/remote/transfer", transferRemoteHandler)
        mux.HandleFunc("/api/remote/search", searchRemoteHandler)
        mux.HandleFunc("/ws/ssh", sshHandler)
        mux.HandleFunc("/ws/share/", shareWSHandler)
        mux.HandleFunc("/ws/events", eventsHandler)

        addr := os.Getenv("LISTEN_ADDR")
        if addr == "" {
                port := os.Getenv("PORT")
                if port == "" {
                        port = "8080"
                }
                addr = ":" + port
        }
        log.Printf("Listening on %s", addr)

        // Background session cleanup: remove stale sessions every hour
        go func() {
                for {
                        time.Sleep(1 * time.Hour)
                        // Delete sessions not active for 7 days
                        cutoff := time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339)
                        res, err := db.Exec(`DELETE FROM auth_sessions WHERE last_active_at != '' AND last_active_at < ?`, cutoff)
                        if err == nil {
                                n, _ := res.RowsAffected()
                                if n > 0 {
                                        log.Printf("Cleaned %d stale auth sessions", n)
                                }
                        }
                        // Also clean sessions with no last_active_at and expired expires_at
                        db.Exec(`DELETE FROM auth_sessions WHERE (last_active_at = '' OR last_active_at IS NULL) AND expires_at < ?`, time.Now().Format(time.RFC3339))
                }
        }()

        srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 30 * time.Second}
        certFile, keyFile := os.Getenv("HTTPS_CERT_FILE"), os.Getenv("HTTPS_KEY_FILE")
        if certFile != "" && keyFile != "" {
                log.Printf("HTTPS enabled (cert %s)", certFile)
                log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
        }
        log.Fatal(srv.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" || r.URL.Path == "" {
                data, err := staticFiles.ReadFile("static/index.html")
                if err != nil {
                        http.Error(w, "Not found", 404)
                        return
                }
                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                w.Header().Set("Cache-Control", "no-cache")
                w.Write(data)
                return
        }
        http.NotFound(w, r)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"version": AppVersion})
}

// ─── DATABASE ────────────────────────────────────────

func resolveDBPath() string {
        if p := os.Getenv("DB_PATH"); p != "" {
                return p
        }
        wd, err := os.Getwd()
        if err != nil {
                return "remote_manager.db"
        }
        return filepath.Join(wd, "remote_manager.db")
}

func initDB() {
        dbPath := resolveDBPath()
        var err error
        db, err = sql.Open("sqlite", dbPath)
        if err != nil {
                log.Fatalf("Open DB: %v", err)
        }
        if err = db.Ping(); err != nil {
                log.Fatalf("Ping DB: %v", err)
        }
        db.Exec(`PRAGMA journal_mode=WAL`)
        db.Exec(`PRAGMA foreign_keys=ON`)
        db.Exec(`PRAGMA busy_timeout=5000`)

        mustExec(`CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT NOT NULL UNIQUE,
                password_hash TEXT NOT NULL,
                is_admin INTEGER NOT NULL DEFAULT 0,
                created_at TEXT NOT NULL)`)
        mustExec(`CREATE TABLE IF NOT EXISTS auth_sessions (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                token TEXT NOT NULL UNIQUE,
                expires_at TEXT NOT NULL)`)
        mustExec(`CREATE TABLE IF NOT EXISTS folders (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                user_id INTEGER DEFAULT NULL)`)
        mustExec(`CREATE TABLE IF NOT EXISTS connections (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                protocol TEXT NOT NULL DEFAULT 'SSH',
                host TEXT NOT NULL,
                username TEXT NOT NULL DEFAULT '',
                auth_method TEXT NOT NULL DEFAULT 'PASSWORD',
                password TEXT DEFAULT '',
                private_key TEXT DEFAULT '',
                key_path TEXT DEFAULT '',
                folder_id INTEGER REFERENCES folders(id) ON DELETE SET NULL,
                user_id INTEGER DEFAULT NULL)`)
        mustExec(`CREATE TABLE IF NOT EXISTS sessions (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                created_at TEXT NOT NULL,
                locked INTEGER NOT NULL DEFAULT 0,
                pane1_conn_id INTEGER REFERENCES connections(id) ON DELETE SET NULL,
                pane1_mode TEXT NOT NULL DEFAULT '',
                pane2_conn_id INTEGER REFERENCES connections(id) ON DELETE SET NULL,
                pane2_mode TEXT NOT NULL DEFAULT '',
                user_id INTEGER DEFAULT NULL)`)
        mustExec(`CREATE TABLE IF NOT EXISTS share_links (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                token TEXT NOT NULL UNIQUE,
                name TEXT NOT NULL,
                password_hash TEXT DEFAULT NULL,
                public INTEGER NOT NULL DEFAULT 1,
                active INTEGER NOT NULL DEFAULT 1,
                created_at TEXT NOT NULL)`)
        mustExec(`CREATE TABLE IF NOT EXISTS share_items (
                share_id INTEGER NOT NULL REFERENCES share_links(id) ON DELETE CASCADE,
                connection_id INTEGER REFERENCES connections(id) ON DELETE CASCADE,
                folder_id INTEGER REFERENCES folders(id) ON DELETE CASCADE)`)
        mustExec(`CREATE TABLE IF NOT EXISTS share_members (
                share_id INTEGER NOT NULL REFERENCES share_links(id) ON DELETE CASCADE,
                user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE)`)

        // Safe migrations
        for _, m := range []string{
                `ALTER TABLE connections ADD COLUMN user_id INTEGER DEFAULT NULL`,
                `ALTER TABLE sessions    ADD COLUMN user_id INTEGER DEFAULT NULL`,
                `ALTER TABLE folders     ADD COLUMN user_id INTEGER DEFAULT NULL`,
                `ALTER TABLE connections ADD COLUMN private_key TEXT DEFAULT ''`,
                `ALTER TABLE connections ADD COLUMN key_path TEXT DEFAULT ''`,
                `ALTER TABLE users       ADD COLUMN is_admin INTEGER NOT NULL DEFAULT 0`,
                `ALTER TABLE auth_sessions ADD COLUMN last_active_at TEXT DEFAULT ''`,
                `ALTER TABLE sessions    ADD COLUMN layout TEXT NOT NULL DEFAULT ''`,
        } {
                if _, err := db.Exec(m); err != nil && !strings.Contains(err.Error(), "duplicate column") {
                        log.Printf("Migration warning: %v", err)
                }
        }

        // Session keepalive: update last_active_at for existing sessions that lack it
        db.Exec(`UPDATE auth_sessions SET last_active_at = expires_at WHERE last_active_at = ''`)

        // Ensure the very first user is always admin
        var userCount int
        db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
        if userCount > 0 {
                db.Exec("UPDATE users SET is_admin=1 WHERE id=(SELECT MIN(id) FROM users)")
        }

        var nC, nF, nS int
        db.QueryRow("SELECT COUNT(*) FROM connections").Scan(&nC)
        db.QueryRow("SELECT COUNT(*) FROM folders").Scan(&nF)
        db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&nS)
        log.Printf("DB ready: %d connections, %d folders, %d sessions", nC, nF, nS)
}

func mustExec(q string) {
        if _, err := db.Exec(q); err != nil {
                log.Fatalf("Schema error: %v\n%s", err, q)
        }
}

// ─── CRYPTO & AUTH HELPERS ───────────────────────────

func getEncryptionKey() []byte {
        if len(encryptionKey) > 0 {
                return encryptionKey
        }
        key := os.Getenv("ENCRYPTION_KEY")
        if key == "" {
                key = defaultEncryptionKey
        }
        if len(key) < 32 {
                key = key + strings.Repeat("0", 32-len(key))
        }
        encryptionKey = []byte(key)[:32]
        return encryptionKey
}

func encryptValue(plaintext string) string {
        if plaintext == "" {
                return ""
        }
        key := getEncryptionKey()
        block, err := aes.NewCipher(key)
        if err != nil {
                return plaintext
        }
        aesgcm, err := cipher.NewGCM(block)
        if err != nil {
                return plaintext
        }
        nonce := make([]byte, aesgcm.NonceSize())
        if _, err := rand.Read(nonce); err != nil {
                return plaintext
        }
        ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
        return "ENC:" + base64.StdEncoding.EncodeToString(ciphertext)
}

func decryptValue(value string) string {
        if value == "" || !strings.HasPrefix(value, "ENC:") {
                return value
        }
        decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, "ENC:"))
        if err != nil {
                return value
        }
        key := getEncryptionKey()
        block, err := aes.NewCipher(key)
        if err != nil {
                return value
        }
        aesgcm, err := cipher.NewGCM(block)
        if err != nil {
                return value
        }
        nonceSize := aesgcm.NonceSize()
        if len(decoded) < nonceSize {
                return value
        }
        plain, err := aesgcm.Open(nil, decoded[:nonceSize], decoded[nonceSize:], nil)
        if err != nil {
                return value
        }
        return string(plain)
}

func hashPassword(password string) (string, error) {
        hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        return string(hash), err
}

func checkPassword(hash, password string) bool {
        return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateToken(length int) string {
        const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
        b := make([]byte, length)
        rand.Read(b)
        for i := range b {
                b[i] = charset[int(b[i])%len(charset)]
        }
        return string(b)
}

func setAuthCookie(w http.ResponseWriter, token string) {
        http.SetCookie(w, &http.Cookie{
                Name:     sessionCookieName,
                Value:    token,
                Path:     "/",
                HttpOnly: true,
                SameSite: http.SameSiteLaxMode,
                MaxAge:   86400 * 30,
        })
}

func clearAuthCookie(w http.ResponseWriter) {
        http.SetCookie(w, &http.Cookie{
                Name:     sessionCookieName,
                Value:    "",
                Path:     "/",
                HttpOnly: true,
                MaxAge:   -1,
        })
}

func currentUserID(r *http.Request) (int, error) {
        cookie, err := r.Cookie(sessionCookieName)
        if err != nil {
                return 0, fmt.Errorf("unauthenticated")
        }
        var userID int
        var expiresAt string
        if err := db.QueryRow(`SELECT user_id, expires_at FROM auth_sessions WHERE token=?`, cookie.Value).Scan(&userID, &expiresAt); err != nil {
                return 0, fmt.Errorf("invalid session")
        }
        if expiresAt != "" {
                exp, err := time.Parse(time.RFC3339, expiresAt)
                if err == nil && time.Now().After(exp) {
                        // Delete expired session
                        db.Exec(`DELETE FROM auth_sessions WHERE token=?`, cookie.Value)
                        return 0, fmt.Errorf("session expired")
                }
        }
        // Update last_active_at to keep session alive
        now := time.Now().Format(time.RFC3339)
        db.Exec(`UPDATE auth_sessions SET last_active_at=? WHERE token=?`, now, cookie.Value)
        return userID, nil
}

func requireAuth(w http.ResponseWriter, r *http.Request) (int, bool) {
        userID, err := currentUserID(r)
        if err != nil {
                jsonError(w, "Unauthorized", 401)
                return 0, false
        }
        return userID, true
}

func createAuthSession(userID int) (string, error) {
        token := generateToken(40)
        expiresAt := time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339)
        _, err := db.Exec(`INSERT INTO auth_sessions (user_id, token, expires_at) VALUES (?,?,?)`, userID, token, expiresAt)
        return token, err
}

func jsonError(w http.ResponseWriter, msg string, code int) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonOK(w http.ResponseWriter, v interface{}) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(v)
}

func boolToInt(v bool) int {
        if v {
                return 1
        }
        return 0
}

func decryptConnectionSecrets(c *Connection) {
        c.Password = decryptValue(c.Password)
        c.PrivateKey = decryptValue(c.PrivateKey)
}

func encryptConnectionSecrets(c *Connection) {
        c.Password = encryptValue(c.Password)
        c.PrivateKey = encryptValue(c.PrivateKey)
}

func userOwnsConnection(connectionID, userID int) bool {
        var count int
        db.QueryRow("SELECT COUNT(1) FROM connections WHERE id=? AND user_id=?", connectionID, userID).Scan(&count)
        return count == 1
}

func userOwnsFolder(folderID, userID int) bool {
        var count int
        db.QueryRow("SELECT COUNT(1) FROM folders WHERE id=? AND user_id=?", folderID, userID).Scan(&count)
        return count == 1
}

func userOwnsSession(sessionID, userID int) bool {
        var count int
        db.QueryRow("SELECT COUNT(1) FROM sessions WHERE id=? AND user_id=?", sessionID, userID).Scan(&count)
        return count == 1
}

// ─── AUTH HANDLERS ───────────────────────────────────

func apiAuthRegisterHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                jsonError(w, "Method not allowed", 405)
                return
        }
        var payload struct {
                Username string `json:"username"`
                Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                jsonError(w, "Bad JSON", 400)
                return
        }
        payload.Username = strings.TrimSpace(payload.Username)
        payload.Password = strings.TrimSpace(payload.Password)
        if payload.Username == "" || payload.Password == "" {
                jsonError(w, "Username and password are required", 400)
                return
        }
        if len(payload.Password) < 6 {
                jsonError(w, "Password must be at least 6 characters", 400)
                return
        }
        hash, err := hashPassword(payload.Password)
        if err != nil {
                jsonError(w, "Server error", 500)
                return
        }
        var existingCount int
        db.QueryRow("SELECT COUNT(*) FROM users").Scan(&existingCount)
        isFirstUser := existingCount == 0
        adminVal := 0
        if isFirstUser {
                adminVal = 1
        }
        res, err := db.Exec(`INSERT INTO users (username, password_hash, is_admin, created_at) VALUES (?,?,?,?)`, payload.Username, hash, adminVal, time.Now().Format(time.RFC3339))
        if err != nil {
                jsonError(w, "Username already exists", 400)
                return
        }
        id, _ := res.LastInsertId()
        token, err := createAuthSession(int(id))
        if err != nil {
                jsonError(w, "Server error", 500)
                return
        }
        setAuthCookie(w, token)
        jsonOK(w, map[string]interface{}{"id": int(id), "username": payload.Username, "is_admin": isFirstUser})
}

func apiAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                jsonError(w, "Method not allowed", 405)
                return
        }
        var payload struct {
                Username string `json:"username"`
                Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                jsonError(w, "Bad JSON", 400)
                return
        }
        ip := clientIP(r)
        if blocked, wait := loginBlocked(ip); blocked {
                jsonError(w, fmt.Sprintf("Too many failed attempts. Try again in %d s.", int(wait.Seconds())+1), 429)
                return
        }
        var userID int
        var passwordHash string
        if err := db.QueryRow(`SELECT id, password_hash FROM users WHERE username=?`, strings.TrimSpace(payload.Username)).Scan(&userID, &passwordHash); err != nil {
                recordLoginResult(ip, false)
                jsonError(w, "Invalid credentials", 401)
                return
        }
        if !checkPassword(passwordHash, payload.Password) {
                recordLoginResult(ip, false)
                jsonError(w, "Invalid credentials", 401)
                return
        }
        recordLoginResult(ip, true)
        token, err := createAuthSession(userID)
        if err != nil {
                jsonError(w, "Server error", 500)
                return
        }
        setAuthCookie(w, token)
        jsonOK(w, map[string]interface{}{"id": userID, "username": strings.TrimSpace(payload.Username)})
}

func apiAuthLogoutHandler(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie(sessionCookieName)
        if err == nil {
                db.Exec("DELETE FROM auth_sessions WHERE token=?", cookie.Value)
        }
        clearAuthCookie(w)
        jsonOK(w, map[string]string{"status": "ok"})
}

func apiAuthMeHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        var user User
        var isAdmin int
        if err := db.QueryRow(`SELECT id, username, created_at, is_admin FROM users WHERE id=?`, userID).Scan(&user.ID, &user.Username, &user.CreatedAt, &isAdmin); err != nil {
                jsonError(w, "Unauthorized", 401)
                return
        }
        jsonOK(w, map[string]interface{}{
                "id":         user.ID,
                "username":   user.Username,
                "created_at": user.CreatedAt,
                "is_admin":   isAdmin == 1,
        })
}

func apiAuthKeepaliveHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        jsonOK(w, map[string]interface{}{
                "user_id": userID,
                "status":  "alive",
        })
}

// ─── ADMIN & USERS ────────────────────────────────────

func requireAdmin(w http.ResponseWriter, r *http.Request) (int, bool) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return 0, false
        }
        var isAdmin int
        if err := db.QueryRow("SELECT is_admin FROM users WHERE id=?", userID).Scan(&isAdmin); err != nil || isAdmin == 0 {
                jsonError(w, "Admin access required", 403)
                return 0, false
        }
        return userID, true
}

func apiUsersHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
                jsonError(w, "Method not allowed", 405)
                return
        }
        _, ok := requireAuth(w, r)
        if !ok {
                return
        }
        rows, err := db.Query("SELECT id, username FROM users ORDER BY username")
        if err != nil {
                jsonError(w, err.Error(), 500)
                return
        }
        defer rows.Close()
        users := []map[string]interface{}{}
        for rows.Next() {
                var id int
                var username string
                rows.Scan(&id, &username)
                users = append(users, map[string]interface{}{"id": id, "username": username})
        }
        jsonOK(w, users)
}

func apiAdminUsersHandler(w http.ResponseWriter, r *http.Request) {
        _, ok := requireAdmin(w, r)
        if !ok {
                return
        }
        if r.Method != http.MethodGet {
                jsonError(w, "Method not allowed", 405)
                return
        }
        rows, err := db.Query("SELECT id, username, is_admin, created_at FROM users ORDER BY username")
        if err != nil {
                jsonError(w, err.Error(), 500)
                return
        }
        defer rows.Close()
        users := []map[string]interface{}{}
        for rows.Next() {
                var id, isAdmin int
                var username, createdAt string
                rows.Scan(&id, &username, &isAdmin, &createdAt)
                users = append(users, map[string]interface{}{
                        "id": id, "username": username,
                        "is_admin": isAdmin == 1, "created_at": createdAt,
                })
        }
        jsonOK(w, users)
}

func apiAdminUserByIDHandler(w http.ResponseWriter, r *http.Request) {
        adminID, ok := requireAdmin(w, r)
        if !ok {
                return
        }
        idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
        id, err := strconv.Atoi(strings.TrimSpace(idStr))
        if err != nil {
                jsonError(w, "Bad ID", 400)
                return
        }
        switch r.Method {
        case http.MethodPut:
                var payload struct {
                        IsAdmin bool `json:"is_admin"`
                }
                json.NewDecoder(r.Body).Decode(&payload)
                if id == adminID && !payload.IsAdmin {
                        jsonError(w, "Cannot remove your own admin status", 400)
                        return
                }
                adminVal := 0
                if payload.IsAdmin {
                        adminVal = 1
                }
                db.Exec("UPDATE users SET is_admin=? WHERE id=?", adminVal, id)
                w.WriteHeader(http.StatusOK)
        case http.MethodDelete:
                if id == adminID {
                        jsonError(w, "Cannot delete your own account", 400)
                        return
                }
                db.Exec("DELETE FROM users WHERE id=?", id)
                w.WriteHeader(http.StatusOK)
        default:
                jsonError(w, "Method not allowed", 405)
        }
}

// ─── SHARE HANDLERS ──────────────────────────────────

func sharePageHandler(w http.ResponseWriter, r *http.Request) {
        data, err := staticFiles.ReadFile("static/index.html")
        if err != nil {
                http.Error(w, "Not found", 404)
                return
        }
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Header().Set("Cache-Control", "no-cache")
        w.Write(data)
}

func getShareByToken(token string) (ShareLink, sql.NullString, error) {
        var share ShareLink
        var ph sql.NullString
        err := db.QueryRow(`SELECT id, owner_id, token, name, password_hash, public, active, created_at FROM share_links WHERE token=?`, token).
                Scan(&share.ID, &share.OwnerID, &share.Token, &share.Name, &ph, &share.Public, &share.Active, &share.CreatedAt)
        return share, ph, err
}

func shareAuthCookieName(token string) string {
        return "wrm_share_" + token[:min(len(token), 8)]
}

func min(a, b int) int {
        if a < b {
                return a
        }
        return b
}

func setShareAccessCookie(w http.ResponseWriter, token string) {
        http.SetCookie(w, &http.Cookie{
                Name:     shareAuthCookieName(token),
                Value:    "1",
                Path:     "/",
                HttpOnly: true,
                SameSite: http.SameSiteLaxMode,
                MaxAge:   86400 * 7,
        })
}

func hasShareAccess(r *http.Request, share ShareLink) bool {
        if !share.Active {
                return false
        }
        // No password required
        var ph sql.NullString
        db.QueryRow(`SELECT password_hash FROM share_links WHERE id=?`, share.ID).Scan(&ph)
        if !ph.Valid || ph.String == "" {
                return true
        }
        // Check share access cookie
        if cookie, err := r.Cookie(shareAuthCookieName(share.Token)); err == nil && cookie.Value == "1" {
                return true
        }
        // Check if logged-in user is a member
        if userID, err := currentUserID(r); err == nil {
                var count int
                db.QueryRow(`SELECT COUNT(1) FROM share_links WHERE id=? AND active=1 AND (owner_id=? OR public=1 OR id IN (SELECT share_id FROM share_members WHERE user_id=?))`, share.ID, userID, userID).Scan(&count)
                if count > 0 {
                        return true
                }
        }
        return false
}

func getSharedConnectionsForShare(share ShareLink) []Connection {
        rows, err := db.Query(`SELECT connection_id, folder_id FROM share_items WHERE share_id=?`, share.ID)
        if err != nil {
                return nil
        }
        defer rows.Close()
        connIDs, folderIDs := []int{}, []int{}
        for rows.Next() {
                var cid, fid sql.NullInt64
                rows.Scan(&cid, &fid)
                if cid.Valid {
                        connIDs = append(connIDs, int(cid.Int64))
                }
                if fid.Valid {
                        folderIDs = append(folderIDs, int(fid.Int64))
                }
        }
        if len(connIDs) == 0 && len(folderIDs) == 0 {
                return nil
        }
        query := `SELECT id,name,protocol,host,username,auth_method,password,private_key,key_path,folder_id FROM connections WHERE user_id=?`
        args := []interface{}{share.OwnerID}
        clauses := []string{}
        if len(connIDs) > 0 {
                ph := strings.TrimRight(strings.Repeat("?,", len(connIDs)), ",")
                clauses = append(clauses, "id IN ("+ph+")")
                for _, id := range connIDs {
                        args = append(args, id)
                }
        }
        if len(folderIDs) > 0 {
                ph := strings.TrimRight(strings.Repeat("?,", len(folderIDs)), ",")
                clauses = append(clauses, "folder_id IN ("+ph+")")
                for _, id := range folderIDs {
                        args = append(args, id)
                }
        }
        query += " AND (" + strings.Join(clauses, " OR ") + ") ORDER BY name"
        connRows, err := db.Query(query, args...)
        if err != nil {
                return nil
        }
        defer connRows.Close()
        connections := []Connection{}
        for connRows.Next() {
                var c Connection
                connRows.Scan(&c.ID, &c.Name, &c.Protocol, &c.Host, &c.Username, &c.AuthMethod, &c.Password, &c.PrivateKey, &c.KeyPath, &c.FolderID)
                decryptConnectionSecrets(&c)
                connections = append(connections, c)
        }
        return connections
}

func apiSharesHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
                userID, ok := requireAuth(w, r)
                if !ok {
                        return
                }
                rows, err := db.Query(`SELECT id, owner_id, token, name, password_hash, public, active, created_at FROM share_links WHERE owner_id=? ORDER BY created_at DESC`, userID)
                if err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                defer rows.Close()
                shares := []map[string]interface{}{}
                for rows.Next() {
                        var s ShareLink
                        var ph sql.NullString
                        rows.Scan(&s.ID, &s.OwnerID, &s.Token, &s.Name, &ph, &s.Public, &s.Active, &s.CreatedAt)
                        shares = append(shares, map[string]interface{}{
                                "id": s.ID, "token": s.Token, "name": s.Name,
                                "public": s.Public, "active": s.Active, "created_at": s.CreatedAt,
                                "password_protected": ph.Valid && ph.String != "",
                        })
                }
                jsonOK(w, shares)
        case http.MethodPost:
                userID, ok := requireAuth(w, r)
                if !ok {
                        return
                }
                var payload struct {
                        Name               string   `json:"name"`
                        ConnectionIDs      []int    `json:"connection_ids"`
                        FolderIDs          []int    `json:"folder_ids"`
                        RecipientUsernames []string `json:"recipient_usernames"`
                        Password           string   `json:"password"`
                        Public             bool     `json:"public"`
                }
                if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                        jsonError(w, "Bad JSON", 400)
                        return
                }
                if payload.Name == "" {
                        payload.Name = "Shared Connections"
                }
                var ph sql.NullString
                if strings.TrimSpace(payload.Password) != "" {
                        hash, err := hashPassword(payload.Password)
                        if err != nil {
                                jsonError(w, "Server error", 500)
                                return
                        }
                        ph = sql.NullString{String: hash, Valid: true}
                }
                token := generateToken(32)
                res, err := db.Exec(`INSERT INTO share_links (owner_id, token, name, password_hash, public, active, created_at) VALUES (?,?,?,?,?,1,?)`,
                        userID, token, payload.Name, ph, boolToInt(payload.Public), time.Now().Format(time.RFC3339))
                if err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                shareID, _ := res.LastInsertId()
                for _, cid := range payload.ConnectionIDs {
                        if userOwnsConnection(cid, userID) {
                                db.Exec(`INSERT INTO share_items (share_id, connection_id) VALUES (?,?)`, shareID, cid)
                        }
                }
                for _, fid := range payload.FolderIDs {
                        if userOwnsFolder(fid, userID) {
                                db.Exec(`INSERT INTO share_items (share_id, folder_id) VALUES (?,?)`, shareID, fid)
                        }
                }
                for _, username := range payload.RecipientUsernames {
                        username = strings.TrimSpace(username)
                        if username == "" {
                                continue
                        }
                        var recipientID int
                        if err := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&recipientID); err == nil && recipientID != userID {
                                db.Exec(`INSERT OR IGNORE INTO share_members (share_id, user_id) VALUES (?,?)`, shareID, recipientID)
                        }
                }
                jsonOK(w, map[string]interface{}{"token": token, "share_url": "/share/" + token})
        default:
                jsonError(w, "Method not allowed", 405)
        }
}

func apiSharesByIDHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        idStr := strings.TrimPrefix(r.URL.Path, "/api/shares/")
        id, err := strconv.Atoi(strings.TrimSpace(idStr))
        if err != nil {
                jsonError(w, "Bad ID", 400)
                return
        }
        if r.Method == http.MethodDelete {
                res, err := db.Exec("DELETE FROM share_links WHERE id=? AND owner_id=?", id, userID)
                if err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                n, _ := res.RowsAffected()
                if n == 0 {
                        jsonError(w, "Not found", 404)
                        return
                }
                w.WriteHeader(http.StatusOK)
        } else {
                jsonError(w, "Method not allowed", 405)
        }
}

func apiShareByTokenHandler(w http.ResponseWriter, r *http.Request) {
        tokenPath := strings.TrimPrefix(r.URL.Path, "/api/share/")
        token := strings.TrimSpace(strings.SplitN(strings.TrimLeft(tokenPath, "/"), "/", 2)[0])
        if token == "" {
                jsonError(w, "Bad token", 400)
                return
        }
        share, ph, err := getShareByToken(token)
        if err != nil {
                jsonError(w, "Share not found", 404)
                return
        }
        if !share.Active {
                jsonError(w, "Share is not active", 403)
                return
        }
        // Password check
        if ph.Valid && ph.String != "" {
                if !hasShareAccess(r, share) {
                        if r.Method == http.MethodGet {
                                jsonOK(w, map[string]interface{}{
                                        "share":             map[string]interface{}{"token": share.Token, "name": share.Name},
                                        "password_required": true,
                                })
                                return
                        }
                        if r.Method == http.MethodPost {
                                var payload struct {
                                        Password string `json:"password"`
                                }
                                json.NewDecoder(r.Body).Decode(&payload)
                                if !checkPassword(ph.String, payload.Password) {
                                        jsonError(w, "Invalid password", 401)
                                        return
                                }
                                setShareAccessCookie(w, share.Token)
                        }
                }
        }
        connections := getSharedConnectionsForShare(share)
        var ownerUsername string
        db.QueryRow("SELECT username FROM users WHERE id=?", share.OwnerID).Scan(&ownerUsername)
        jsonOK(w, map[string]interface{}{
                "share": map[string]interface{}{
                        "token": share.Token, "name": share.Name,
                        "public": share.Public, "active": share.Active,
                        "owner_username": ownerUsername,
                },
                "connections": connections,
        })
}

func apiSharedHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
                jsonError(w, "Method not allowed", 405)
                return
        }
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        rows, err := db.Query(`
                SELECT s.id, s.owner_id, s.token, s.name, s.password_hash, s.public, s.active, s.created_at, u.username
                FROM share_links s
                JOIN users u ON u.id=s.owner_id
                WHERE s.active=1 AND s.owner_id != ? AND (
                        s.public=1 OR s.id IN (SELECT share_id FROM share_members WHERE user_id=?)
                ) ORDER BY s.created_at DESC`, userID, userID)
        if err != nil {
                jsonError(w, err.Error(), 500)
                return
        }
        defer rows.Close()
        results := []map[string]interface{}{}
        for rows.Next() {
                var s ShareLink
                var ph sql.NullString
                var ownerUsername string
                rows.Scan(&s.ID, &s.OwnerID, &s.Token, &s.Name, &ph, &s.Public, &s.Active, &s.CreatedAt, &ownerUsername)
                connections := getSharedConnectionsForShare(s)
                results = append(results, map[string]interface{}{
                        "share": map[string]interface{}{
                                "id": s.ID, "token": s.Token, "name": s.Name,
                                "public": s.Public, "active": s.Active,
                                "owner_username": ownerUsername,
                                "password_protected": ph.Valid && ph.String != "",
                        },
                        "connections": connections,
                })
        }
        jsonOK(w, results)
}

// ─── CONNECTIONS ─────────────────────────────────────

func apiConnectionsHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        switch r.Method {
        case http.MethodGet:
                rows, err := db.Query(`SELECT id,name,protocol,host,username,auth_method,password,private_key,key_path,folder_id FROM connections WHERE user_id=? ORDER BY name`, userID)
                if err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                defer rows.Close()
                conns := make([]Connection, 0)
                for rows.Next() {
                        var c Connection
                        rows.Scan(&c.ID, &c.Name, &c.Protocol, &c.Host, &c.Username, &c.AuthMethod, &c.Password, &c.PrivateKey, &c.KeyPath, &c.FolderID)
                        decryptConnectionSecrets(&c)
                        conns = append(conns, c)
                }
                jsonOK(w, conns)
        case http.MethodPost:
                var c Connection
                if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
                        jsonError(w, "Bad JSON", 400)
                        return
                }
                if c.Name == "" || c.Host == "" {
                        jsonError(w, "Name and host are required", 400)
                        return
                }
                if c.Protocol == "" {
                        c.Protocol = "SSH"
                }
                if c.AuthMethod == "" {
                        c.AuthMethod = "PASSWORD"
                }
                encryptConnectionSecrets(&c)
                res, err := db.Exec(`INSERT INTO connections (name,protocol,host,username,auth_method,password,private_key,key_path,folder_id,user_id) VALUES (?,?,?,?,?,?,?,?,?,?)`,
                        c.Name, c.Protocol, c.Host, c.Username, c.AuthMethod, c.Password, c.PrivateKey, c.KeyPath, c.FolderID, userID)
                if err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                id, _ := res.LastInsertId()
                c.ID = int(id)
                decryptConnectionSecrets(&c)
                w.WriteHeader(http.StatusCreated)
                jsonOK(w, c)
        default:
                jsonError(w, "Method not allowed", 405)
        }
}

func apiConnectionByIDHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        idStr := strings.TrimPrefix(r.URL.Path, "/api/connections/")
        id, err := strconv.Atoi(strings.TrimSpace(idStr))
        if err != nil {
                jsonError(w, "Bad ID", 400)
                return
        }
        if !userOwnsConnection(id, userID) {
                jsonError(w, "Not found", 404)
                return
        }
        switch r.Method {
        case http.MethodPut:
                var c Connection
                json.NewDecoder(r.Body).Decode(&c)
                if c.Protocol == "" {
                        c.Protocol = "SSH"
                }
                if c.AuthMethod == "" {
                        c.AuthMethod = "PASSWORD"
                }
                encryptConnectionSecrets(&c)
                if _, err := db.Exec(`UPDATE connections SET name=?,protocol=?,host=?,username=?,auth_method=?,password=?,private_key=?,key_path=?,folder_id=? WHERE id=?`,
                        c.Name, c.Protocol, c.Host, c.Username, c.AuthMethod, c.Password, c.PrivateKey, c.KeyPath, c.FolderID, id); err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                w.WriteHeader(http.StatusOK)
        case http.MethodDelete:
                db.Exec("DELETE FROM connections WHERE id=? AND user_id=?", id, userID)
                w.WriteHeader(http.StatusOK)
        default:
                jsonError(w, "Method not allowed", 405)
        }
}

func apiConnectionsBulkHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        if r.Method != http.MethodPost {
                jsonError(w, "Method not allowed", 405)
                return
        }
        var payload struct {
                Action   string `json:"action"`
                IDs      []int  `json:"ids"`
                FolderID *int   `json:"folder_id"`
        }
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                jsonError(w, "Bad JSON", 400)
                return
        }
        if len(payload.IDs) == 0 {
                jsonError(w, "No IDs provided", 400)
                return
        }
        ph := strings.TrimRight(strings.Repeat("?,", len(payload.IDs)), ",")
        args := make([]interface{}, len(payload.IDs))
        for i, id := range payload.IDs {
                args[i] = id
        }
        switch strings.ToLower(payload.Action) {
        case "delete":
                db.Exec("DELETE FROM connections WHERE id IN ("+ph+") AND user_id=?", append(args, userID)...)
        case "move":
                db.Exec("UPDATE connections SET folder_id=? WHERE id IN ("+ph+") AND user_id=?", append([]interface{}{payload.FolderID}, append(args, userID)...)...)
        default:
                jsonError(w, "Unsupported action", 400)
                return
        }
        w.WriteHeader(http.StatusOK)
}

// ─── FOLDERS ─────────────────────────────────────────

func apiFoldersHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        switch r.Method {
        case http.MethodGet:
                rows, _ := db.Query("SELECT id,name FROM folders WHERE user_id=? ORDER BY name", userID)
                defer rows.Close()
                folders := make([]Folder, 0)
                for rows.Next() {
                        var f Folder
                        rows.Scan(&f.ID, &f.Name)
                        folders = append(folders, f)
                }
                jsonOK(w, folders)
        case http.MethodPost:
                var f Folder
                json.NewDecoder(r.Body).Decode(&f)
                if f.Name == "" {
                        jsonError(w, "Name required", 400)
                        return
                }
                res, _ := db.Exec("INSERT INTO folders (name,user_id) VALUES (?,?)", f.Name, userID)
                id, _ := res.LastInsertId()
                f.ID = int(id)
                w.WriteHeader(http.StatusCreated)
                jsonOK(w, f)
        default:
                jsonError(w, "Method not allowed", 405)
        }
}

func apiFolderByIDHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        idStr := strings.TrimPrefix(r.URL.Path, "/api/folders/")
        id, _ := strconv.Atoi(strings.TrimSpace(idStr))
        if r.Method == http.MethodDelete {
                if !userOwnsFolder(id, userID) {
                        jsonError(w, "Not found", 404)
                        return
                }
                tx, _ := db.Begin()
                tx.Exec("UPDATE connections SET folder_id=NULL WHERE folder_id=? AND user_id=?", id, userID)
                tx.Exec("DELETE FROM folders WHERE id=? AND user_id=?", id, userID)
                tx.Commit()
                w.WriteHeader(http.StatusOK)
        } else {
                jsonError(w, "Method not allowed", 405)
        }
}

// ─── SESSIONS ─────────────────────────────────────────

func loadSessionRows(rows *sql.Rows) []Session {
        sessions := make([]Session, 0)
        for rows.Next() {
                var s Session
                var locked int
                var p1name, p2name *string
                if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt, &locked,
                        &s.Pane1ConnID, &s.Pane1Mode, &s.Pane2ConnID, &s.Pane2Mode,
                        &p1name, &p2name, &s.Layout); err != nil {
                        continue
                }
                s.Locked = locked == 1
                if p1name != nil {
                        s.Pane1ConnName = *p1name
                }
                if p2name != nil {
                        s.Pane2ConnName = *p2name
                }
                sessions = append(sessions, s)
        }
        return sessions
}

func apiSessionsHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        switch r.Method {
        case http.MethodGet:
                rows, err := db.Query(`
                        SELECT s.id, s.name, s.created_at, s.locked,
                               s.pane1_conn_id, s.pane1_mode, s.pane2_conn_id, s.pane2_mode,
                               c1.name, c2.name, s.layout
                        FROM sessions s
                        LEFT JOIN connections c1 ON s.pane1_conn_id = c1.id
                        LEFT JOIN connections c2 ON s.pane2_conn_id = c2.id
                        WHERE s.user_id=? ORDER BY s.created_at DESC`, userID)
                if err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                defer rows.Close()
                jsonOK(w, loadSessionRows(rows))
        case http.MethodPost:
                var s Session
                json.NewDecoder(r.Body).Decode(&s)
                if s.Name == "" {
                        s.Name = time.Now().Format("Jan 2, 2006 15:04")
                }
                if s.CreatedAt == "" {
                        s.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
                }
                locked := 0
                if s.Locked {
                        locked = 1
                }
                if len(s.Layout) > maxSessionLayout {
                        jsonError(w, "Session layout too large", 400)
                        return
                }
                res, err := db.Exec(`INSERT INTO sessions (name,created_at,locked,pane1_conn_id,pane1_mode,pane2_conn_id,pane2_mode,user_id,layout) VALUES (?,?,?,?,?,?,?,?,?)`,
                        s.Name, s.CreatedAt, locked, s.Pane1ConnID, s.Pane1Mode, s.Pane2ConnID, s.Pane2Mode, userID, s.Layout)
                if err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                id, _ := res.LastInsertId()
                s.ID = int(id)
                w.WriteHeader(http.StatusCreated)
                jsonOK(w, s)
                broadcastSessionUpdate()
        default:
                jsonError(w, "Method not allowed", 405)
        }
}

func apiSessionByIDHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        idStr := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
        id, err := strconv.Atoi(strings.TrimSpace(idStr))
        if err != nil {
                jsonError(w, "Bad ID", 400)
                return
        }
        if !userOwnsSession(id, userID) {
                jsonError(w, "Not found", 404)
                return
        }
        switch r.Method {
        case http.MethodPut:
                var s struct {
                        Session
                        Layout *string `json:"layout"` // nil = keep the stored layout
                }
                if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
                        jsonError(w, "Bad JSON", 400)
                        return
                }
                var curLocked int
                var curLayout string
                db.QueryRow("SELECT locked, layout FROM sessions WHERE id=?", id).Scan(&curLocked, &curLayout)
                layout := curLayout
                if s.Layout != nil {
                        if curLocked == 1 && s.Locked && *s.Layout != curLayout {
                                jsonError(w, "Session is locked — unlock it before overwriting its windows", 403)
                                return
                        }
                        if len(*s.Layout) > maxSessionLayout {
                                jsonError(w, "Session layout too large", 400)
                                return
                        }
                        layout = *s.Layout
                }
                if strings.TrimSpace(s.Name) == "" {
                        jsonError(w, "Name required", 400)
                        return
                }
                locked := 0
                if s.Locked {
                        locked = 1
                }
                if _, err := db.Exec(`UPDATE sessions SET name=?,locked=?,pane1_conn_id=?,pane1_mode=?,pane2_conn_id=?,pane2_mode=?,layout=? WHERE id=?`,
                        s.Name, locked, s.Pane1ConnID, s.Pane1Mode, s.Pane2ConnID, s.Pane2Mode, layout, id); err != nil {
                        jsonError(w, err.Error(), 500)
                        return
                }
                w.WriteHeader(http.StatusOK)
                broadcastSessionUpdate()
        case http.MethodDelete:
                var locked int
                db.QueryRow("SELECT locked FROM sessions WHERE id=?", id).Scan(&locked)
                if locked == 1 && r.URL.Query().Get("force") != "true" {
                        jsonError(w, "Session is locked. Add ?force=true to delete.", 403)
                        return
                }
                db.Exec("DELETE FROM sessions WHERE id=? AND user_id=?", id, userID)
                w.WriteHeader(http.StatusOK)
                broadcastSessionUpdate()
        default:
                jsonError(w, "Method not allowed", 405)
        }
}

// ─── EXPORT/IMPORT ───────────────────────────────────

func apiExportHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        if r.Method != http.MethodGet {
                jsonError(w, "Method not allowed", 405)
                return
        }
        fRows, _ := db.Query("SELECT id,name FROM folders WHERE user_id=? ORDER BY name", userID)
        defer fRows.Close()
        folders := []Folder{}
        for fRows.Next() {
                var f Folder
                fRows.Scan(&f.ID, &f.Name)
                folders = append(folders, f)
        }
        cRows, _ := db.Query(`SELECT id,name,protocol,host,username,auth_method,password,private_key,key_path,folder_id FROM connections WHERE user_id=? ORDER BY name`, userID)
        defer cRows.Close()
        conns := []Connection{}
        for cRows.Next() {
                var c Connection
                cRows.Scan(&c.ID, &c.Name, &c.Protocol, &c.Host, &c.Username, &c.AuthMethod, &c.Password, &c.PrivateKey, &c.KeyPath, &c.FolderID)
                decryptConnectionSecrets(&c)
                conns = append(conns, c)
        }
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Content-Disposition", `attachment; filename="wrm-config-`+AppVersion+`.json"`)
        json.NewEncoder(w).Encode(map[string]interface{}{"folders": folders, "connections": conns})
}

func apiImportHandler(w http.ResponseWriter, r *http.Request) {
        userID, ok := requireAuth(w, r)
        if !ok {
                return
        }
        if r.Method != http.MethodPost {
                jsonError(w, "Method not allowed", 405)
                return
        }
        var payload struct {
                Folders     []Folder     `json:"folders"`
                Connections []Connection `json:"connections"`
        }
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                jsonError(w, "Bad JSON", 400)
                return
        }
        tx, _ := db.Begin()
        defer tx.Rollback()
        folderIDs := map[string]*int{}
        for _, f := range payload.Folders {
                if f.Name == "" {
                        continue
                }
                var existingID int
                if err := tx.QueryRow("SELECT id FROM folders WHERE name=? AND user_id=?", f.Name, userID).Scan(&existingID); err == nil {
                        id := existingID
                        folderIDs[f.Name] = &id
                        continue
                }
                res, _ := tx.Exec("INSERT INTO folders (name, user_id) VALUES (?,?)", f.Name, userID)
                id, _ := res.LastInsertId()
                i := int(id)
                folderIDs[f.Name] = &i
        }
        for _, c := range payload.Connections {
                if c.Name == "" || c.Host == "" {
                        continue
                }
                var folderID *int
                if c.FolderID != nil {
                        for _, f := range payload.Folders {
                                if f.ID == *c.FolderID && folderIDs[f.Name] != nil {
                                        folderID = folderIDs[f.Name]
                                        break
                                }
                        }
                }
                encryptConnectionSecrets(&c)
                tx.Exec(`INSERT INTO connections (name,protocol,host,username,auth_method,password,private_key,key_path,folder_id,user_id) VALUES (?,?,?,?,?,?,?,?,?,?)`,
                        c.Name, c.Protocol, c.Host, c.Username, c.AuthMethod, c.Password, c.PrivateKey, c.KeyPath, folderID, userID)
        }
        tx.Commit()
        w.WriteHeader(http.StatusOK)
}

// ─── EVENTS WS ───────────────────────────────────────

func eventsHandler(w http.ResponseWriter, r *http.Request) {
        if _, err := currentUserID(r); err != nil {
                http.Error(w, "Unauthorized", 401)
                return
        }
        ws, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                return
        }
        defer ws.Close()
        ch := hub.subscribe()
        defer hub.unsubscribe(ch)
        // Reader: detects closed/dead clients (the old version leaked subscribers).
        gone := make(chan struct{})
        ws.SetReadDeadline(time.Now().Add(90 * time.Second))
        ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(90 * time.Second)); return nil })
        go func() {
                defer close(gone)
                for {
                        if _, _, err := ws.ReadMessage(); err != nil {
                                return
                        }
                }
        }()
        t := time.NewTicker(25 * time.Second)
        defer t.Stop()
        for {
                select {
                case <-gone:
                        return
                case <-t.C:
                        if ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)) != nil {
                                return
                        }
                case msg := <-ch:
                        ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
                        if ws.WriteMessage(websocket.TextMessage, msg) != nil {
                                return
                        }
                }
        }
}

// ─── SSH KEY HELPERS ──────────────────────────────────

func parsePrivateKeyBytes(data []byte) (ssh.Signer, error) {
        trimmed := strings.TrimSpace(string(data))
        if strings.HasPrefix(trimmed, "ssh-") || strings.HasPrefix(trimmed, "ecdsa-") {
                return nil, fmt.Errorf("this is a public key; use the private key file")
        }
        s, err := ssh.ParsePrivateKey([]byte(trimmed))
        if err != nil {
                return nil, fmt.Errorf("key parse error: %v", err)
        }
        return s, nil
}

func findPrivateForPub(pubPath string) ([]byte, string, error) {
        cands := []string{}
        if strings.HasSuffix(pubPath, ".pub") {
                cands = append(cands, strings.TrimSuffix(pubPath, ".pub"))
        }
        home, _ := os.UserHomeDir()
        for _, n := range []string{"id_rsa", "id_ed25519", "id_ecdsa", "id_dsa"} {
                cands = append(cands, filepath.Join(home, ".ssh", n))
        }
        for _, p := range cands {
                data, err := os.ReadFile(p)
                if err != nil {
                        continue
                }
                if strings.Contains(string(data), "PRIVATE KEY") {
                        return data, p, nil
                }
        }
        return nil, "", fmt.Errorf("private key not found")
}

func expandPath(p string) string {
        if strings.HasPrefix(p, "~/") || p == "~" {
                home, _ := os.UserHomeDir()
                return filepath.Join(home, p[2:])
        }
        return filepath.FromSlash(p)
}

func resolveKeyMaterial(c Connection) ([]byte, error) {
        switch c.AuthMethod {
        case "KEY":
                d := strings.TrimSpace(c.PrivateKey)
                if d == "" {
                        return nil, fmt.Errorf("private key not provided")
                }
                return []byte(d), nil
        case "KEY_FILE":
                p := expandPath(strings.TrimSpace(c.KeyPath))
                if p == "" {
                        return nil, fmt.Errorf("key path not provided")
                }
                data, err := os.ReadFile(p)
                if err != nil {
                        return nil, fmt.Errorf("cannot read '%s': %v", p, err)
                }
                t := strings.TrimSpace(string(data))
                if strings.HasSuffix(strings.ToLower(p), ".pub") || strings.HasPrefix(t, "ssh-") || strings.HasPrefix(t, "ecdsa-") {
                        priv, found, err := findPrivateForPub(p)
                        if err != nil {
                                return nil, fmt.Errorf("found .pub but no private key: %v", err)
                        }
                        log.Printf("Using key: %s", found)
                        return priv, nil
                }
                return data, nil
        case "KEY_AUTO":
                home, _ := os.UserHomeDir()
                for _, n := range []string{"id_rsa", "id_ed25519", "id_ecdsa", "id_dsa"} {
                        p := filepath.Join(home, ".ssh", n)
                        data, err := os.ReadFile(p)
                        if err != nil {
                                continue
                        }
                        if strings.Contains(string(data), "PRIVATE KEY") {
                                return data, nil
                        }
                }
                return nil, fmt.Errorf("no key found in ~/.ssh/")
        }
        return nil, fmt.Errorf("unknown auth method: %s", c.AuthMethod)
}

func buildAuthMethods(c Connection) ([]ssh.AuthMethod, error) {
        if c.AuthMethod == "PASSWORD" || c.AuthMethod == "" {
                return []ssh.AuthMethod{ssh.Password(c.Password)}, nil
        }
        kd, err := resolveKeyMaterial(c)
        if err != nil {
                return nil, err
        }
        sg, err := parsePrivateKeyBytes(kd)
        if err != nil {
                return nil, err
        }
        return []ssh.AuthMethod{ssh.PublicKeys(sg)}, nil
}

func ensurePort(host, protocol string) string {
        if strings.Contains(host, ":") {
                return host
        }
        if strings.ToUpper(protocol) == "FTP" || strings.ToUpper(protocol) == "FTPS" {
                return host + ":21"
        }
        return host + ":22"
}

func getSSHClient(c Connection) (*ssh.Client, error) {
        am, err := buildAuthMethods(c)
        if err != nil {
                return nil, err
        }
        return ssh.Dial("tcp", c.Host, &ssh.ClientConfig{
                User: c.Username, Auth: am,
                HostKeyCallback: ssh.InsecureIgnoreHostKey(),
                Timeout:         15 * time.Second,
        })
}

func isConnectionAccessible(r *http.Request, connID int) bool {
        if userID, err := currentUserID(r); err == nil && userOwnsConnection(connID, userID) {
                return true
        }
        token := strings.TrimSpace(r.URL.Query().Get("share_token"))
        if token == "" {
                return false
        }
        share, _, err := getShareByToken(token)
        if err != nil || !share.Active || !hasShareAccess(r, share) {
                return false
        }
        // Check if connection is in share
        var count int
        db.QueryRow(`SELECT COUNT(1) FROM share_items WHERE share_id=? AND connection_id=?`, share.ID, connID).Scan(&count)
        if count > 0 {
                return true
        }
        var folderID sql.NullInt64
        if err := db.QueryRow(`SELECT folder_id FROM connections WHERE id=?`, connID).Scan(&folderID); err == nil && folderID.Valid {
                db.QueryRow(`SELECT COUNT(1) FROM share_items WHERE share_id=? AND folder_id=?`, share.ID, folderID.Int64).Scan(&count)
                return count > 0
        }
        return false
}

// ─── SERVER-TO-SERVER TRANSFER ───────────────────────

func transferRemoteHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                jsonError(w, "Method not allowed", 405)
                return
        }
        var req struct {
                SourceConnID int      `json:"source_conn_id"`
                SourceFiles  []string `json:"source_files"`
                DestConnID   int      `json:"dest_conn_id"`
                DestPath     string   `json:"dest_path"`
                OnConflict   string   `json:"on_conflict"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                jsonError(w, "Invalid request body", 400)
                return
        }
        if req.SourceConnID == 0 || req.DestConnID == 0 || len(req.SourceFiles) == 0 {
                jsonError(w, "source_conn_id, dest_conn_id, and source_files are required", 400)
                return
        }
        if !isConnectionAccessible(r, req.SourceConnID) || !isConnectionAccessible(r, req.DestConnID) {
                jsonError(w, "Access denied", 403)
                return
        }

        // Load connections
        var srcConn Connection
        if err := db.QueryRow(`SELECT protocol,host,username,auth_method,password,private_key,key_path FROM connections WHERE id=?`, req.SourceConnID).
                Scan(&srcConn.Protocol, &srcConn.Host, &srcConn.Username, &srcConn.AuthMethod, &srcConn.Password, &srcConn.PrivateKey, &srcConn.KeyPath); err != nil {
                jsonError(w, "Source connection not found", 404)
                return
        }
        decryptConnectionSecrets(&srcConn)
        srcConn.Host = ensurePort(srcConn.Host, srcConn.Protocol)

        var dstConn Connection
        if err := db.QueryRow(`SELECT protocol,host,username,auth_method,password,private_key,key_path FROM connections WHERE id=?`, req.DestConnID).
                Scan(&dstConn.Protocol, &dstConn.Host, &dstConn.Username, &dstConn.AuthMethod, &dstConn.Password, &dstConn.PrivateKey, &dstConn.KeyPath); err != nil {
                jsonError(w, "Destination connection not found", 404)
                return
        }
        decryptConnectionSecrets(&dstConn)
        dstConn.Host = ensurePort(dstConn.Host, dstConn.Protocol)

        if strings.ToUpper(srcConn.Protocol) == "FTP" || strings.ToUpper(dstConn.Protocol) == "FTP" {
                jsonError(w, "FTP not supported for server-to-server transfer", 400)
                return
        }

        // Setup streaming response IMMEDIATELY so frontend gets feedback
        w.Header().Set("Content-Type", "application/x-ndjson")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("X-Accel-Buffering", "no")
        flusher, canFlush := w.(http.Flusher)

        sendInitial := func(msg interface{}) {
                json.NewEncoder(w).Encode(msg)
                if canFlush { flusher.Flush() }
        }

        ctx := r.Context()
        clientGone := ctx.Done()

        // Connect to source SFTP
        sendInitial(map[string]interface{}{"type": "status", "message": "Connecting to source server..."})
        srcSSHClient, err := getSSHClient(srcConn)
        if err != nil {
                sendInitial(map[string]interface{}{"type": "error", "error": "Source SSH failed: " + err.Error()})
                return
        }
        defer srcSSHClient.Close()
        sendInitial(map[string]interface{}{"type": "status", "message": "Source SSH connected. Opening SFTP..."})
        srcSFTP, err := sftp.NewClient(srcSSHClient)
        if err != nil {
                sendInitial(map[string]interface{}{"type": "error", "error": "Source SFTP failed: " + err.Error()})
                return
        }
        defer srcSFTP.Close()
        sendInitial(map[string]interface{}{"type": "status", "message": "Source ready. Connecting to destination..."})

        // Connect to destination SFTP
        dstSSHClient, err := getSSHClient(dstConn)
        if err != nil {
                sendInitial(map[string]interface{}{"type": "error", "error": "Destination SSH failed: " + err.Error()})
                return
        }
        defer dstSSHClient.Close()
        sendInitial(map[string]interface{}{"type": "status", "message": "Destination SSH connected. Opening SFTP..."})
        dstSFTP, err := sftp.NewClient(dstSSHClient)
        if err != nil {
                sendInitial(map[string]interface{}{"type": "error", "error": "Destination SFTP failed: " + err.Error()})
                return
        }
        defer dstSFTP.Close()
        sendInitial(map[string]interface{}{"type": "status", "message": "Both servers connected. Starting transfer..."})

        if req.DestPath != "" && req.DestPath != "." {
                dstSFTP.MkdirAll(req.DestPath)
        }

        var sendMu sync.Mutex
        sendMsg := func(msg interface{}) error {
                sendMu.Lock()
                defer sendMu.Unlock()
                select {
                case <-clientGone:
                        return fmt.Errorf("client disconnected")
                default:
                }
                if err := json.NewEncoder(w).Encode(msg); err != nil {
                        return err
                }
                if canFlush {
                        flusher.Flush()
                }
                return nil
        }

        type transferResult struct {
                File   string `json:"file"`
                Status string `json:"status"`
                Size   int64  `json:"size"`
                Error  string `json:"error,omitempty"`
        }

        // Ensure directory exists on destination
        ensureDir := func(dirPath string) error {
                if dirPath == "" || dirPath == "/" || dirPath == "." {
                        return nil
                }
                info, err := dstSFTP.Stat(dirPath)
                if err == nil {
                        if info.IsDir() {
                                return nil
                        }
                        dstSFTP.Remove(dirPath)
                }
                if err := dstSFTP.MkdirAll(dirPath); err != nil {
                        return fmt.Errorf("mkdir %s: %v", dirPath, err)
                }
                info2, err := dstSFTP.Stat(dirPath)
                if err != nil || !info2.IsDir() {
                        return fmt.Errorf("cannot create directory %s", dirPath)
                }
                return nil
        }

        // Transfer a single file + stream results
        transferFile := func(srcPath, dstPath string) transferResult {
                result := transferResult{File: srcPath, Status: "ok"}

                srcFileObj, err := srcSFTP.Open(srcPath)
                if err != nil {
                        result.Status = "failed"
                        result.Error = "Cannot open source: " + err.Error()
                        return result
                }

                srcInfo, err := srcFileObj.Stat()
                if err != nil {
                        srcFileObj.Close()
                        result.Status = "failed"
                        result.Error = "Cannot stat source: " + err.Error()
                        return result
                }
                result.Size = srcInfo.Size()

                // Send file_start before beginning transfer
                sendMsg(map[string]interface{}{
                        "type": "file_start",
                        "file": srcPath,
                        "size": srcInfo.Size(),
                })

                dstDir := path.Dir(dstPath) // remote paths always use "/", even when WRM runs on Windows
                if err := ensureDir(dstDir); err != nil {
                        srcFileObj.Close()
                        result.Status = "failed"
                        result.Error = "Cannot create directory: " + err.Error()
                        return result
                }

                // Check conflict
                if dstInfo, err := dstSFTP.Stat(dstPath); err == nil && !dstInfo.IsDir() {
                        switch req.OnConflict {
                        case "skip":
                                srcFileObj.Close()
                                result.Status = "skipped"
                                result.Size = 0
                                return result
                        case "newer":
                                if !srcInfo.ModTime().After(dstInfo.ModTime()) {
                                        srcFileObj.Close()
                                        result.Status = "skipped"
                                        result.Size = 0
                                        return result
                                }
                        case "abort":
                                srcFileObj.Close()
                                result.Status = "conflict"
                                result.Error = "File exists: " + dstPath
                                return result
                        }
                }

                dstFileObj, err := dstSFTP.Create(dstPath)
                if err != nil {
                        srcFileObj.Close()
                        result.Status = "failed"
                        result.Error = "Cannot create destination: " + err.Error()
                        return result
                }

                buf := make([]byte, 1024*1024)
                var written int64
                lastProgress := time.Now()
                for {
                        select {
                        case <-clientGone:
                                srcFileObj.Close()
                                dstFileObj.Close()
                                result.Status = "failed"
                                result.Error = "Client disconnected"
                                return result
                        default:
                        }
                        n, readErr := srcFileObj.Read(buf)
                        if n > 0 {
                                nw, writeErr := dstFileObj.Write(buf[:n])
                                if writeErr != nil {
                                        srcFileObj.Close()
                                        dstFileObj.Close()
                                        result.Status = "failed"
                                        result.Error = "Write error: " + writeErr.Error()
                                        return result
                                }
                                written += int64(nw)
                                // Send progress every 500ms
                                if time.Since(lastProgress) >= 500*time.Millisecond {
                                        sendMsg(map[string]interface{}{
                                                "type":    "file_progress",
                                                "file":    srcPath,
                                                "written": written,
                                                "total":   srcInfo.Size(),
                                        })
                                        lastProgress = time.Now()
                                }
                        }
                        if readErr != nil {
                                break
                        }
                }
                srcFileObj.Close()
                dstFileObj.Close()

                if !srcInfo.ModTime().IsZero() {
                        dstSFTP.Chtimes(dstPath, srcInfo.ModTime(), srcInfo.ModTime())
                }

                log.Printf("Transfer OK: %s -> %s (%d bytes)", srcPath, dstPath, written)
                return result
        }

        // Collect all files to transfer (flatten dirs)
        type fileJob struct {
                srcPath string
                dstPath string
        }
        var jobs []fileJob

        for _, srcFile := range req.SourceFiles {
                fileName := path.Base(strings.TrimSuffix(srcFile, "/"))
                dstFullPath := path.Join(req.DestPath, fileName)

                srcInfo, err := srcSFTP.Stat(srcFile)
                if err != nil {
                        r := transferResult{File: srcFile, Status: "failed", Error: "Source not found: " + err.Error()}
                        sendMsg(map[string]interface{}{"type": "file", "result": r})
                        continue
                }

                if srcInfo.IsDir() {
                        ensureDir(dstFullPath)
                        var collectFiles func(string)
                        collectFiles = func(dir string) {
                                entries, err := srcSFTP.ReadDir(dir)
                                if err != nil {
                                        r := transferResult{File: dir, Status: "failed", Error: "Cannot list: " + err.Error()}
                                        sendMsg(map[string]interface{}{"type": "file", "result": r})
                                        return
                                }
                                for _, entry := range entries {
                                        sPath := dir + "/" + entry.Name()
                                        dPath := path.Join(dstFullPath, strings.TrimPrefix(sPath, srcFile))
                                        if entry.IsDir() {
                                                collectFiles(sPath)
                                        } else {
                                                jobs = append(jobs, fileJob{srcPath: sPath, dstPath: dPath})
                                        }
                                }
                        }
                        collectFiles(srcFile)
                } else {
                        jobs = append(jobs, fileJob{srcPath: srcFile, dstPath: dstFullPath})
                }
        }

        // Recount total size from actual jobs
        var totalSize int64
        for _, j := range jobs {
                if info, err := srcSFTP.Stat(j.srcPath); err == nil {
                        totalSize += info.Size()
                }
        }

        // Send progress with correct file count
        if err := sendMsg(map[string]interface{}{
                "type":      "progress",
                "total":     len(jobs),
                "totalSize": totalSize,
        }); err != nil {
                log.Printf("Transfer cancelled: client disconnected")
                return
        }

        // Transfer files with worker pool (4 concurrent workers)
        const numWorkers = 4
        var mu sync.Mutex
        okCount, failCount, skipCount := 0, 0, 0
        var errorList []string
        var clientDisconnected atomic.Bool
        jobCh := make(chan fileJob, len(jobs))
        for _, j := range jobs {
                jobCh <- j
        }
        close(jobCh)

        var wg sync.WaitGroup
        for w := 0; w < numWorkers; w++ {
                wg.Add(1)
                go func() {
                        defer wg.Done()
                        for j := range jobCh {
                                if clientDisconnected.Load() {
                                        return
                                }
                                select {
                                case <-clientGone:
                                        clientDisconnected.Store(true)
                                        return
                                default:
                                }
                                r := transferFile(j.srcPath, j.dstPath)
                                mu.Lock()
                                sendMsg(map[string]interface{}{"type": "file", "result": r})
                                switch r.Status {
                                case "ok": okCount++
                                case "skipped": skipCount++
                                default:
                                        failCount++
                                        if r.Error != "" {
                                                errorList = append(errorList, j.srcPath+": "+r.Error)
                                        }
                                }
                                mu.Unlock()
                        }
                }()
        }
        wg.Wait()

        // Send done
        if !clientDisconnected.Load() {
                sendMsg(map[string]interface{}{
                        "type":      "done",
                        "ok":        okCount,
                        "failed":    failCount,
                        "skipped":   skipCount,
                        "errorList": errorList,
                })
        }
}


