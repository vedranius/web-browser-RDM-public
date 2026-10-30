package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

// Databases used by earlier WRM builds can contain tables with the names v10 uses but a
// different layout. They must be kept aside and recreated, not break the upgrade.
func TestLegacyTablesUpgraded(t *testing.T) {
	legacy, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "legacy.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer legacy.Close()
	for _, q := range []string{
		`CREATE TABLE collab_messages (id INTEGER PRIMARY KEY AUTOINCREMENT, share_token TEXT NOT NULL,
			user TEXT NOT NULL, type TEXT NOT NULL DEFAULT 'chat', text TEXT NOT NULL DEFAULT '',
			file_name TEXT NOT NULL DEFAULT '', timestamp TEXT NOT NULL)`,
		`INSERT INTO collab_messages (share_token, user, text, timestamp) VALUES ('tok', 'ana', 'hello', '2025-01-01')`,
		`CREATE TABLE audit_log (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at TEXT NOT NULL, action TEXT NOT NULL)`,
		`CREATE INDEX ix_audit_ts ON audit_log(created_at)`,
		`CREATE TABLE app_settings (key TEXT PRIMARY KEY, value TEXT NOT NULL DEFAULT '')`,
		`INSERT INTO app_settings (key, value) VALUES ('max_file_size_bytes', '100')`,
	} {
		if _, err := legacy.Exec(q); err != nil {
			t.Fatalf("%v\n%s", err, q)
		}
	}

	saved := db
	db = legacy
	defer func() { db = saved }()

	ensureTable("collab_messages", `CREATE TABLE IF NOT EXISTS collab_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT, share_id INTEGER NOT NULL, ts TEXT NOT NULL,
		pkey TEXT NOT NULL DEFAULT '', name TEXT NOT NULL DEFAULT '', kind TEXT NOT NULL DEFAULT 'chat',
		text TEXT NOT NULL DEFAULT '')`, "id", "share_id", "ts", "pkey", "name", "kind", "text")
	ensureTable("audit_log", `CREATE TABLE IF NOT EXISTS audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT, ts TEXT NOT NULL, user_id INTEGER DEFAULT NULL,
		username TEXT NOT NULL DEFAULT '', ip TEXT NOT NULL DEFAULT '', action TEXT NOT NULL,
		target TEXT NOT NULL DEFAULT '', details TEXT NOT NULL DEFAULT '')`,
		"id", "ts", "user_id", "username", "ip", "action", "target", "details")
	ensureTable("app_settings", `CREATE TABLE IF NOT EXISTS app_settings (key TEXT PRIMARY KEY, value TEXT NOT NULL)`, "key", "value")

	// New layout works.
	if _, err := db.Exec(`INSERT INTO collab_messages (share_id, ts, pkey, name, kind, text) VALUES (1, 'now', 'u:1', 'ana', 'chat', 'hi')`); err != nil {
		t.Fatalf("insert into upgraded collab_messages: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO audit_log (ts, action) VALUES ('now', 'test')`); err != nil {
		t.Fatalf("insert into upgraded audit_log: %v", err)
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS ix_audit_ts ON audit_log(ts)`); err != nil {
		t.Fatalf("index on upgraded audit_log: %v", err)
	}
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='ix_audit_ts' AND tbl_name='audit_log'`).Scan(&n)
	if n != 1 {
		t.Fatalf("ix_audit_ts not on the new audit_log")
	}

	// Old data is kept under a new name.
	var old string
	if err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name LIKE 'collab_messages_old_%'`).Scan(&old); err != nil {
		t.Fatalf("legacy collab_messages not kept: %v", err)
	}
	var text string
	if err := db.QueryRow(`SELECT text FROM "` + old + `"`).Scan(&text); err != nil || text != "hello" {
		t.Fatalf("legacy data = %q, %v", text, err)
	}

	// A compatible table is left alone.
	var v string
	if err := db.QueryRow(`SELECT value FROM app_settings WHERE key='max_file_size_bytes'`).Scan(&v); err != nil || v != "100" {
		t.Fatalf("compatible app_settings was touched: %q, %v", v, err)
	}
	db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name LIKE 'app_settings_old_%'`).Scan(&n)
	if n != 0 {
		t.Fatalf("compatible app_settings was renamed")
	}
}

func TestStaticOldURLsRedirect(t *testing.T) {
	h := staticHandler()
	for path, want := range map[string]int{
		"/static/":                http.StatusFound,
		"/static/index.html":      http.StatusFound,
		"/static/vendor/":         http.StatusNotFound,
		"/static/nope.js":         http.StatusNotFound,
		"/static/vendor/xterm.js": http.StatusOK,
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		if rec.Code != want {
			t.Errorf("%s: status %d, want %d", path, rec.Code, want)
		}
		if want == http.StatusFound && rec.Header().Get("Location") != "/" {
			t.Errorf("%s: redirects to %q, want /", path, rec.Header().Get("Location"))
		}
	}
}

func TestBrowserURL(t *testing.T) {
	for _, c := range []struct {
		addr string
		tls  bool
		want string
	}{
		{":8080", false, "http://localhost:8080/"},
		{"0.0.0.0:9000", true, "https://localhost:9000/"},
		{"127.0.0.1:8080", false, "http://127.0.0.1:8080/"},
		{"[::1]:8443", true, "https://[::1]:8443/"},
	} {
		if got := browserURL(c.addr, c.tls); got != c.want {
			t.Errorf("browserURL(%q, %v) = %q, want %q", c.addr, c.tls, got, c.want)
		}
	}
}

// A table whose CREATE statement has a NOT NULL column without a default that the caller
// did not list (addon_instances.addon_id) must not be set aside again on every start.
func TestInitDBIsIdempotent(t *testing.T) {
	saved := db
	defer func() { db = saved }()
	t.Setenv("DB_PATH", filepath.Join(t.TempDir(), "restart.db"))
	initDB()
	if _, err := db.Exec(`INSERT INTO addon_instances (addon_owner, addon_id, owner_id, name) VALUES (0, 'demo.panel', 0, 'kept')`); err != nil {
		t.Fatal(err)
	}
	db.Close()
	initDB()
	defer db.Close()
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name LIKE '%\_old\_%' ESCAPE '\'`)
	if err != nil {
		t.Fatal(err)
	}
	var moved []string
	for rows.Next() {
		var n string
		rows.Scan(&n)
		moved = append(moved, n)
	}
	rows.Close()
	if len(moved) > 0 {
		t.Fatalf("a second start set tables aside: %v", moved)
	}
	var name string
	if err := db.QueryRow(`SELECT name FROM addon_instances WHERE addon_id='demo.panel'`).Scan(&name); err != nil || name != "kept" {
		t.Fatalf("addon instance lost after restart: %q, %v", name, err)
	}
}

func TestCreateColumns(t *testing.T) {
	got := createColumns(`CREATE TABLE IF NOT EXISTS x (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		"Name" TEXT NOT NULL DEFAULT '',
		price NUMERIC(10, 2) NOT NULL,
		UNIQUE (id, name),
		PRIMARY KEY (id)
	)`)
	want := []string{"id", "name", "price"}
	if len(got) != len(want) {
		t.Fatalf("createColumns = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("createColumns = %v, want %v", got, want)
		}
	}
}

// Instances set aside by the earlier start-up bug come back with their ids, so the
// addon_store and addon_windows rows that point at them work again.
func TestAddonInstancesRestored(t *testing.T) {
	saved := db
	defer func() { db = saved }()
	t.Setenv("DB_PATH", filepath.Join(t.TempDir(), "restore.db"))
	initDB()
	defer db.Close()
	if _, err := db.Exec(`DELETE FROM addon_instance_restore_migration WHERE id=1`); err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		`DELETE FROM addon_instances`,
		`CREATE TABLE addon_instances_old_20261027080042 AS SELECT * FROM addon_instances`,
		`INSERT INTO addon_instances_old_20261027080042 (id, addon_owner, addon_id, owner_id, name, created_at, created_by) VALUES (2, 0, 'demo.panel', 1, 'Old', '', '')`,
		`CREATE TABLE addon_instances_old_20261028073353 AS SELECT * FROM addon_instances`,
		`INSERT INTO addon_instances_old_20261028073353 (id, addon_owner, addon_id, owner_id, name, created_at, created_by) VALUES (1, 0, 'demo.panel', 1, 'Default', '', ''), (3, 0, 'demo.panel', 0, 'shared', '', '')`,
		`CREATE TABLE addon_instances_old_20261028074752 AS SELECT * FROM addon_instances`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatalf("%v\n%s", err, q)
		}
	}
	restoreAddonInstances()
	var ids []int
	rows, _ := db.Query(`SELECT id FROM addon_instances ORDER BY id`)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	rows.Close()
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 3 {
		t.Fatalf("restored ids = %v, want [1 3] from the newest non-empty copy", ids)
	}
	// A table that already has instances is left alone.
	db.Exec(`DELETE FROM addon_instances WHERE id=3`)
	restoreAddonInstances()
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM addon_instances`).Scan(&n)
	if n != 1 {
		t.Fatalf("restore touched a non-empty table: %d rows", n)
	}
	if _, err := db.Exec(`DELETE FROM addon_instances`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	initDB()
	defer db.Close()
	if err := db.QueryRow(`SELECT COUNT(*) FROM addon_instances`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("deleted addon instances were restored again after restart: %d rows", n)
	}
}
