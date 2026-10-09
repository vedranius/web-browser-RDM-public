// Git workspace of WRM PRO: compare services on servers with GitLab / GitHub or offline bundles.
// Loaded on demand by openGit() in index.html; uses its helpers (t, esc, showToast, apiError,
// uiConfirm, fmtTime, conns, folders, currentUser).
'use strict';

Object.assign(LANGS.en, {
  git_title: 'Git workspace', git_close: 'Close', git_tab_overview: 'Overview', git_tab_services: 'Services', git_tab_installs: 'Installations', git_tab_settings: 'Settings',
  git_check_now: 'Check now', git_checking: 'Checking…', git_last_check: 'Last check', git_never: 'never', git_no_checks: 'Checks are not allowed for your account.',
  git_st_ok: 'Up to date', git_st_update: 'Needs update', git_st_review: 'Review', git_st_error: 'Error', git_st_gone: 'Gone', git_st_unknown: 'Not checked',
  git_fs_ok: 'ok', git_fs_old: 'old', git_fs_modified: 'modified', git_fs_missing: 'missing', git_fs_extra: 'extra', git_fs_protected: 'protected',
  git_behind: '{n} behind', git_vmd_missing: 'missing', git_version_md: 'VERSION.md', git_target: 'Target', git_latest_tag: 'Latest tag',
  git_all_states: 'All states', git_all_envs: 'All environments', git_all_folders: 'All folders', git_all_tags: 'All tags', git_search: 'Search…',
  git_ref_choice: 'Ref', git_ref_catalog: 'As in the catalog', git_ref_default: 'Default branch of each project', git_ref_branch: 'One branch for all',
  git_ref_override: 'Ref for this run', git_ref_override_ph: 'tag:latest, tag:v2 or a branch (empty = as chosen above)',
  git_server: 'Server', git_service: 'Service', git_path: 'Path', git_env: 'Environment', git_state: 'State', git_files: 'Files', git_checked: 'Checked', git_actions: 'Actions',
  git_empty_overview: 'Nothing to show yet. Add a source or import a bundle, fill in the catalog, choose the servers in Settings and run a check.',
  git_targets: 'Targets', git_refresh_targets: 'Refresh targets', git_source: 'Source', git_warnings: 'Warnings', git_no_target: 'no target yet',
  git_add_service: 'Add service', git_suggest: 'Suggest from repository', git_import_catalog: 'Import catalog', git_export_catalog: 'Export catalog',
  git_import_merge: 'Merge with the current catalog (OK) or replace it (Cancel)?', git_export_bundle: 'Export bundle', git_export_bundle_none: 'Select services with a target first.',
  git_from_bundle: 'From imported bundles (not in the catalog)', git_edit: 'Edit', git_delete: 'Delete', git_delete_q: 'Delete service {name} from the catalog?', git_save: 'Save', git_cancel: 'Cancel',
  git_f_name: 'Name', git_f_project: 'Project (group/name)', git_f_ref: 'Ref', git_f_ref_h: 'tag:latest, tag:v1 or a branch', git_f_branch: 'Branch (for tags)', git_f_tag_filter: 'Tag filter (regex)',
  git_f_subdir: 'Subdirectory = installation root', git_f_kind: 'Kind', git_kind_app: 'App', git_kind_library: 'Library', git_kind_tool: 'Tool',
  git_f_include: 'Include globs', git_f_exclude: 'Exclude globs', git_f_protected: 'Protected (never overwritten, per host)', git_f_fingerprint: 'Fingerprint files', git_f_hint: 'Install hints (directory names)',
  git_lines_h: 'One per line or comma-separated', git_suggest_project: 'Project', git_suggest_ref: 'Branch (empty = default)', git_suggest_go: 'Suggest', git_repo_files: 'Repository files',
  git_discover: 'Discover', git_add_install: 'Add by hand', git_details: 'Details', git_forget: 'Forget', git_forget_q: 'Forget this installation? Discovery may find it again.',
  git_units: 'Units', git_no_units: 'No systemd or supervisor unit mentions this directory.', git_diff: 'Diff', git_diff_title: 'Server → target', git_binary: 'Binary file',
  git_no_diff: 'No differences in the text (line endings or whitespace only).', git_capped: 'Only the first 5000 files were checked.', git_files_only: 'Show only changes',
  git_sources: 'Git sources', git_add_source: 'Add source', git_kind: 'Kind', git_url: 'Address', git_token: 'Read-only token', git_token_keep: 'leave empty to keep the stored token',
  git_test: 'Test', git_projects_found: '{n} projects', git_bundles: 'Offline bundles', git_import_bundle: 'Import bundle (.tar.gz)', git_bundle_age: 'created {d}',
  git_bundle_imported: 'Bundle imported: {n} files', git_source_mode: 'Targets from', git_mode_auto: 'Automatic (Git API, otherwise the newest bundle)', git_mode_live: 'Git API',
  git_mode_bundle: 'Offline bundle', git_newest_bundle: 'Newest bundle with the service', git_history_depth: 'History depth (earlier versions per file)',
  git_servers: 'Servers to check', git_servers_h: 'SSH connections. Checks only read files, nothing is changed on servers.', git_filter: 'Filter…', git_select_all: 'All', git_select_none: 'None',
  git_catalog_settings: 'Catalog settings', git_groups: 'Groups / organisations', git_fallback: 'Fallback branches', git_roots: 'Server roots (searched to depth 3)', git_ignore: 'Ignored directories (globs)',
  git_monitor: 'Periodic checks', git_monitor_on: 'Check my servers every {n} minutes while WRM runs and notify me', git_monitor_off: 'Periodic checks are turned off by the administrator.',
  git_monitor_h: 'Choose the channels for "Git" events in Settings → Notifications. Checks never update anything.', git_saved: 'Saved',
  git_n_installs: '{n} installations', git_none: '—', git_api_calls: '{n} API calls', git_from: 'from', git_bundle_only: 'bundle',
});
Object.assign(LANGS.hr, {
  git_title: 'Git radni prostor', git_close: 'Zatvori', git_tab_overview: 'Pregled', git_tab_services: 'Servisi', git_tab_installs: 'Instalacije', git_tab_settings: 'Postavke',
  git_check_now: 'Provjeri sada', git_checking: 'Provjera…', git_last_check: 'Zadnja provjera', git_never: 'nikad', git_no_checks: 'Provjere nisu dopuštene vašem računu.',
  git_st_ok: 'Ažurno', git_st_update: 'Treba ažurirati', git_st_review: 'Pregledati', git_st_error: 'Greška', git_st_gone: 'Nema je', git_st_unknown: 'Nije provjereno',
  git_fs_ok: 'ok', git_fs_old: 'staro', git_fs_modified: 'izmijenjeno', git_fs_missing: 'nedostaje', git_fs_extra: 'višak', git_fs_protected: 'zaštićeno',
  git_behind: '{n} iza', git_vmd_missing: 'nema', git_version_md: 'VERSION.md', git_target: 'Cilj', git_latest_tag: 'Zadnji tag',
  git_all_states: 'Sva stanja', git_all_envs: 'Sva okruženja', git_all_folders: 'Sve mape', git_all_tags: 'Sve oznake', git_search: 'Traži…',
  git_ref_choice: 'Ref', git_ref_catalog: 'Kao u katalogu', git_ref_default: 'Zadana grana svakog projekta', git_ref_branch: 'Jedna grana za sve',
  git_ref_override: 'Ref za ovu provjeru', git_ref_override_ph: 'tag:latest, tag:v2 ili grana (prazno = kao gore)',
  git_server: 'Server', git_service: 'Servis', git_path: 'Putanja', git_env: 'Okruženje', git_state: 'Stanje', git_files: 'Datoteke', git_checked: 'Provjereno', git_actions: 'Akcije',
  git_empty_overview: 'Još nema ničega. Dodajte izvor ili uvezite bundle, popunite katalog, odaberite servere u Postavkama i pokrenite provjeru.',
  git_targets: 'Ciljevi', git_refresh_targets: 'Osvježi ciljeve', git_source: 'Izvor', git_warnings: 'Upozorenja', git_no_target: 'još nema cilja',
  git_add_service: 'Dodaj servis', git_suggest: 'Prijedlog iz repozitorija', git_import_catalog: 'Uvezi katalog', git_export_catalog: 'Izvezi katalog',
  git_import_merge: 'Spojiti s trenutnim katalogom (U redu) ili ga zamijeniti (Odustani)?', git_export_bundle: 'Izvezi bundle', git_export_bundle_none: 'Najprije odaberite servise koji imaju cilj.',
  git_from_bundle: 'Iz uvezenih bundleova (nisu u katalogu)', git_edit: 'Uredi', git_delete: 'Obriši', git_delete_q: 'Obrisati servis {name} iz kataloga?', git_save: 'Spremi', git_cancel: 'Odustani',
  git_f_name: 'Naziv', git_f_project: 'Projekt (grupa/naziv)', git_f_ref: 'Ref', git_f_ref_h: 'tag:latest, tag:v1 ili grana', git_f_branch: 'Grana (za tagove)', git_f_tag_filter: 'Filtar tagova (regex)',
  git_f_subdir: 'Poddirektorij = korijen instalacije', git_f_kind: 'Vrsta', git_kind_app: 'Aplikacija', git_kind_library: 'Biblioteka', git_kind_tool: 'Alat',
  git_f_include: 'Uključi (globovi)', git_f_exclude: 'Isključi (globovi)', git_f_protected: 'Zaštićeno (nikad se ne prepisuje, po hostu)', git_f_fingerprint: 'Datoteke za prepoznavanje', git_f_hint: 'Nazivi direktorija instalacije',
  git_lines_h: 'Jedno po retku ili odvojeno zarezom', git_suggest_project: 'Projekt', git_suggest_ref: 'Grana (prazno = zadana)', git_suggest_go: 'Predloži', git_repo_files: 'Datoteke repozitorija',
  git_discover: 'Pronađi', git_add_install: 'Dodaj ručno', git_details: 'Detalji', git_forget: 'Zaboravi', git_forget_q: 'Zaboraviti ovu instalaciju? Pretraga je može ponovno pronaći.',
  git_units: 'Jedinice', git_no_units: 'Nijedna systemd ili supervisor jedinica ne spominje ovaj direktorij.', git_diff: 'Razlike', git_diff_title: 'Server → cilj', git_binary: 'Binarna datoteka',
  git_no_diff: 'Nema razlika u tekstu (samo završeci redaka ili razmaci).', git_capped: 'Provjereno je samo prvih 5000 datoteka.', git_files_only: 'Prikaži samo promjene',
  git_sources: 'Git izvori', git_add_source: 'Dodaj izvor', git_kind: 'Vrsta', git_url: 'Adresa', git_token: 'Token samo za čitanje', git_token_keep: 'ostavite prazno za zadržavanje spremljenog tokena',
  git_test: 'Testiraj', git_projects_found: '{n} projekata', git_bundles: 'Offline bundleovi', git_import_bundle: 'Uvezi bundle (.tar.gz)', git_bundle_age: 'izrađen {d}',
  git_bundle_imported: 'Bundle uvezen: {n} datoteka', git_source_mode: 'Ciljevi iz', git_mode_auto: 'Automatski (Git API, inače najnoviji bundle)', git_mode_live: 'Git API',
  git_mode_bundle: 'Offline bundle', git_newest_bundle: 'Najnoviji bundle sa servisom', git_history_depth: 'Dubina povijesti (ranije verzije po datoteci)',
  git_servers: 'Serveri za provjeru', git_servers_h: 'SSH veze. Provjere samo čitaju datoteke, na serverima se ništa ne mijenja.', git_filter: 'Filtar…', git_select_all: 'Sve', git_select_none: 'Ništa',
  git_catalog_settings: 'Postavke kataloga', git_groups: 'Grupe / organizacije', git_fallback: 'Zamjenske grane', git_roots: 'Korijeni na serverima (pretraga do dubine 3)', git_ignore: 'Zanemareni direktoriji (globovi)',
  git_monitor: 'Periodične provjere', git_monitor_on: 'Provjeravaj moje servere svakih {n} minuta dok WRM radi i obavijesti me', git_monitor_off: 'Administrator je isključio periodične provjere.',
  git_monitor_h: 'Kanale za "Git" događaje odaberite u Postavke → Obavijesti. Provjere nikad ništa ne ažuriraju.', git_saved: 'Spremljeno',
  git_n_installs: '{n} instalacija', git_none: '—', git_api_calls: '{n} API poziva', git_from: 'iz', git_bundle_only: 'bundle',
});

Object.assign(LANGS.en, {
  git_tab_runs: 'Runs', git_update: 'Update', git_upgrade: 'Upgrade…', git_rollback: 'Rollback…', git_restart: 'Restart…', git_stamp: 'Stamp',
  git_select_installs: 'Select installations first.', git_restart_pending: 'restart pending', git_fs_changed: 'differs',
  git_k_update: 'Update', git_k_upgrade: 'Upgrade', git_k_rollback: 'Rollback', git_k_stamp: 'Stamp', git_k_restart: 'Restart',
  git_rs_scheduled: 'Scheduled', git_rs_starting: 'Starting', git_rs_running: 'Running', git_rs_done: 'Done', git_rs_failed: 'Failed', git_rs_cancelled: 'Cancelled',
  git_rs_skipped: 'Skipped', git_rs_interrupted: 'Interrupted', git_rs_pending: 'Waiting', git_rs_ok: 'OK', git_rs_rolled_back: 'Rolled back',
  git_plan_loading: 'Preparing…', git_upgrade_ref: 'Branch or tag to upgrade to', git_upgrade_ref_h: 'A branch (main) or a tag (tag:v2, tag:latest)', git_show_changes: 'Show changes',
  git_upgrade_warn: 'An upgrade switches the installations to another branch or version. Read the summary of the changes before you continue.',
  git_changes_summary: '{changed} to change, {missing} new, {local} changed locally, {removed} no longer in the target (they stay on the server)',
  git_prod: 'production', git_files_none: 'Nothing to change: the files are up to date.', git_show_all: 'Show all files', git_local: 'changed locally',
  git_local_warn: 'Ticked files with local changes are overwritten: the changes made on the server are lost (they stay in the backup).',
  git_check_cmd: 'Custom check for {app} (optional, runs in the installation directory; exit code 0 = OK)', git_ignore_imports: 'Skip the check for dangling Python imports',
  git_restart_choice: 'Service restart', git_rm_none: 'Do not restart (only show what needs a restart)', git_rm_now: 'Restart right after a successful update',
  git_rm_at: 'Restart at a chosen time', git_rm_delay: 'Restart after a delay (minutes)', git_restart_warn: 'These units will be restarted: {list}',
  git_no_units_restart: 'No systemd or supervisor unit is known for these installations (check them first).', git_when: 'When', git_run_now: 'Now', git_run_at: 'At a chosen time (maintenance window)',
  git_grace_h: 'A scheduled run is skipped when WRM was not running at its time for more than {n} minutes. A notification is sent before and after.',
  git_set_override: 'Follow this ref from now on (ref override for these services)', git_after_upgrade: 'Comparisons follow the service\'s ref (catalog or override) until you change it.',
  git_continue: 'Continue', git_back: 'Back', git_confirm_title: 'Confirm: {kind}', git_confirm_q: '{kind} of {n} installation(s) on {m} server(s):',
  git_type_host: 'Production: type the server name {name} to confirm', git_start: 'Start', git_schedule: 'Schedule', git_scheduled_ok: 'Scheduled for {t}',
  git_nothing_selected: 'No files are selected.', git_time_past: 'Choose a time in the future.', git_progress: 'Run {label}', git_stop: 'Stop after the current server',
  git_cancel_run: 'Cancel run', git_cancel_q: 'Cancel this run?', git_dangling: 'Dangling imports', git_health: 'Health check', git_backup: 'Backup',
  git_backups: 'Backups on the server', git_no_backups: 'No backups in .deploy-bak.', git_backup_files: '{n} files', git_backup_unknown: 'No record of this backup: only the files in it are restored.',
  git_force_stamp: 'Overwrite an existing VERSION.md', git_stamp_h: 'Stamp writes VERSION.md where an installation is up to date (and the file is missing).',
  git_not_current: 'not up to date', git_has_vmd: 'VERSION.md exists', git_runs_empty: 'No runs yet.', git_created: 'Created', git_scheduled_for: 'Scheduled for',
  git_run: 'Run', git_kind_col: 'Action', git_items: 'Installations', git_written: '{n} written', git_deleted: '{n} deleted', git_waits_update: 'after the update',
  git_restore_hint: 'Restores the files of this backup; files that did not exist then are deleted. The current files are backed up first.', git_versions: 'Versions',
  git_refresh: 'Refresh', git_unit_list: 'Units', git_no_permission: 'Not allowed for your account.',
  git_step_connect: 'Log in', git_step_scan: 'Compare', git_step_prepare: 'Prepare', git_step_imports: 'Python imports', git_step_backups: 'Backups', git_step_backup: 'Backup',
  git_step_write: 'Write', git_step_checks: 'Checks', git_step_rollback: 'Rollback', git_step_version: 'VERSION.md', git_step_log: 'updates.jsonl', git_step_restart: 'Restart',
});
Object.assign(LANGS.hr, {
  git_tab_runs: 'Izvođenja', git_update: 'Ažuriraj', git_upgrade: 'Nadogradi…', git_rollback: 'Vrati…', git_restart: 'Ponovno pokreni…', git_stamp: 'Označi',
  git_select_installs: 'Najprije odaberite instalacije.', git_restart_pending: 'čeka ponovno pokretanje', git_fs_changed: 'razlikuje se',
  git_k_update: 'Ažuriranje', git_k_upgrade: 'Nadogradnja', git_k_rollback: 'Vraćanje', git_k_stamp: 'Označavanje', git_k_restart: 'Ponovno pokretanje',
  git_rs_scheduled: 'Zakazano', git_rs_starting: 'Pokreće se', git_rs_running: 'U tijeku', git_rs_done: 'Gotovo', git_rs_failed: 'Nije uspjelo', git_rs_cancelled: 'Otkazano',
  git_rs_skipped: 'Preskočeno', git_rs_interrupted: 'Prekinuto', git_rs_pending: 'Čeka', git_rs_ok: 'U redu', git_rs_rolled_back: 'Vraćeno',
  git_plan_loading: 'Priprema…', git_upgrade_ref: 'Grana ili tag za nadogradnju', git_upgrade_ref_h: 'Grana (main) ili tag (tag:v2, tag:latest)', git_show_changes: 'Prikaži promjene',
  git_upgrade_warn: 'Nadogradnja prebacuje instalacije na drugu granu ili verziju. Prije nastavka pročitajte sažetak promjena.',
  git_changes_summary: '{changed} za promjenu, {missing} novih, {local} lokalno izmijenjenih, {removed} više nije u cilju (ostaju na serveru)',
  git_prod: 'produkcija', git_files_none: 'Nema promjena: datoteke su ažurne.', git_show_all: 'Prikaži sve datoteke', git_local: 'lokalno izmijenjeno',
  git_local_warn: 'Označene datoteke s lokalnim izmjenama bit će prepisane: izmjene na serveru gube se (ostaju u sigurnosnoj kopiji).',
  git_check_cmd: 'Vlastita provjera za {app} (neobavezno, pokreće se u direktoriju instalacije; izlazni kod 0 = u redu)', git_ignore_imports: 'Preskoči provjeru visećih Python importa',
  git_restart_choice: 'Ponovno pokretanje servisa', git_rm_none: 'Ne pokreći ponovno (samo prikaži što treba ponovno pokrenuti)', git_rm_now: 'Ponovno pokreni odmah nakon uspješnog ažuriranja',
  git_rm_at: 'Ponovno pokreni u odabrano vrijeme', git_rm_delay: 'Ponovno pokreni nakon odgode (minute)', git_restart_warn: 'Ponovno će se pokrenuti ove jedinice: {list}',
  git_no_units_restart: 'Za ove instalacije nije poznata nijedna systemd ili supervisor jedinica (najprije ih provjerite).', git_when: 'Kada', git_run_now: 'Sada', git_run_at: 'U odabrano vrijeme (prozor održavanja)',
  git_grace_h: 'Zakazano izvođenje preskače se ako WRM u to vrijeme nije radio dulje od {n} minuta. Obavijest se šalje prije i poslije.',
  git_set_override: 'Ubuduće prati ovaj ref (zamjenski ref za ove servise)', git_after_upgrade: 'Usporedbe prate ref servisa (katalog ili zamjenski ref) dok ga ne promijenite.',
  git_continue: 'Nastavi', git_back: 'Natrag', git_confirm_title: 'Potvrda: {kind}', git_confirm_q: '{kind}: {n} instalacija na {m} servera:',
  git_type_host: 'Produkcija: upišite naziv servera {name} za potvrdu', git_start: 'Pokreni', git_schedule: 'Zakaži', git_scheduled_ok: 'Zakazano za {t}',
  git_nothing_selected: 'Nije odabrana nijedna datoteka.', git_time_past: 'Odaberite vrijeme u budućnosti.', git_progress: 'Izvođenje {label}', git_stop: 'Zaustavi nakon trenutnog servera',
  git_cancel_run: 'Otkaži izvođenje', git_cancel_q: 'Otkazati ovo izvođenje?', git_dangling: 'Viseći importi', git_health: 'Provjera stanja', git_backup: 'Sigurnosna kopija',
  git_backups: 'Sigurnosne kopije na serveru', git_no_backups: 'U .deploy-bak nema sigurnosnih kopija.', git_backup_files: '{n} datoteka', git_backup_unknown: 'Nema zapisa o ovoj kopiji: vraćaju se samo datoteke u njoj.',
  git_force_stamp: 'Prepiši postojeći VERSION.md', git_stamp_h: 'Označavanje upisuje VERSION.md gdje je instalacija ažurna (a datoteka nedostaje).',
  git_not_current: 'nije ažurno', git_has_vmd: 'VERSION.md postoji', git_runs_empty: 'Još nema izvođenja.', git_created: 'Izrađeno', git_scheduled_for: 'Zakazano za',
  git_run: 'Izvođenje', git_kind_col: 'Akcija', git_items: 'Instalacije', git_written: 'zapisano: {n}', git_deleted: 'obrisano: {n}', git_waits_update: 'nakon ažuriranja',
  git_restore_hint: 'Vraća datoteke iz ove kopije; datoteke koje tada nisu postojale brišu se. Trenutne datoteke najprije se spremaju u novu kopiju.', git_versions: 'Verzije',
  git_refresh: 'Osvježi', git_unit_list: 'Jedinice', git_no_permission: 'Nije dopušteno vašem računu.',
  git_step_connect: 'Prijava', git_step_scan: 'Usporedba', git_step_prepare: 'Priprema', git_step_imports: 'Python importi', git_step_backups: 'Sigurnosne kopije', git_step_backup: 'Sigurnosna kopija',
  git_step_write: 'Zapisivanje', git_step_checks: 'Provjere', git_step_rollback: 'Vraćanje', git_step_version: 'VERSION.md', git_step_log: 'updates.jsonl', git_step_restart: 'Ponovno pokretanje',
});

Object.assign(LANGS.en, {
  git_install: 'Install…', git_transfer: 'Transfer…', git_k_install: 'Install', git_k_transfer: 'Transfer', git_new_server: 'New server',
  git_w_server: 'Target server', git_w_services: 'Services', git_w_ref: 'Version', git_w_ref_target: 'Target ({v})', git_w_ref_other: 'Branch or tag…', git_w_ref_bundle: 'Bundle {b}',
  git_w_pick_service: 'Choose at least one service.', git_w_dirs: 'Target directories', git_w_roots: 'Server roots', git_w_root_missing: 'does not exist', git_w_root_ro: 'not writable',
  git_w_files: '{n} files, {kb} KB', git_w_check: 'Check the target', git_w_checks: 'Pre-checks', git_w_free: '{free} MB free, {need} MB needed', git_w_tools: 'Tools',
  git_w_missing_tool: 'missing', git_w_existing: 'An installation exists here ({n} files{v}).', git_w_overwrite: 'Replace it (the replaced files are backed up to .deploy-bak first)',
  git_w_protected: 'Per-host files', git_w_protected_h: 'Protected files are never taken from the repository. Fill them from a template, copy them from another installation or type them.',
  git_w_no_slots: 'The catalog lists no protected files for this service.', git_w_m_template: 'From a template', git_w_m_copy: 'Copy from', git_w_m_manual: 'Type the content', git_w_m_skip: 'Leave it out',
  git_w_tmpl_na: 'not available', git_w_in_repo: 'in the repository', git_w_typed: 'Typed value', git_w_vault_pw: '🔑 {n}: password', git_w_vault_user: '🔑 {n}: user name',
  git_w_vault_h: 'Values from the vault are read on the server side only and are not stored with the run.',
  git_w_unit: 'Create a service unit (optional)', git_w_unit_kind: 'Kind', git_w_unit_name: 'Name', git_w_unit_user: 'Run as user (empty = default)', git_w_unit_cmd: 'Command (absolute path)',
  git_w_unit_dir: 'Working directory', git_w_unit_enable: 'Start at boot (systemctl enable)',
  git_w_unit_warn: 'The unit file is written to the server as root (or with sudo -n) and is never written over an existing one. WRM does not test the command: check it.',
  git_w_start: 'Start the service after the install (restart of its units)', git_w_start_h: 'Units that mention the directory (also an existing or the new one) are started; their health is checked.',
  git_w_blocked: 'Fix the problems above first.', git_w_confirm: '{kind}: {n} service(s) on {server}', git_w_type: 'Type the server name {name} to confirm (production or overwrite)',
  git_w_source: 'Source', git_w_excluded: '{n} paths are not copied (logs, caches, backups, ignored directories)', git_w_tr_h: 'Copies the installation with its per-host configuration. Virtual environments and node_modules are not copied: recreate them on the target.',
  git_w_units_not_copied: 'Units of the source are not copied: {list}', git_w_same: 'Choose another server or directory than the source.', git_w_env: 'Environment',
  git_step_precheck: 'Pre-checks', git_step_fetch: 'Contents', git_step_receive: 'Copy to the server', git_step_register: 'Register', git_step_unit: 'Unit', git_step_source: 'Source',
});
Object.assign(LANGS.hr, {
  git_install: 'Instaliraj…', git_transfer: 'Prenesi…', git_k_install: 'Instalacija', git_k_transfer: 'Prijenos', git_new_server: 'Novi server',
  git_w_server: 'Ciljni server', git_w_services: 'Servisi', git_w_ref: 'Verzija', git_w_ref_target: 'Cilj ({v})', git_w_ref_other: 'Grana ili tag…', git_w_ref_bundle: 'Bundle {b}',
  git_w_pick_service: 'Odaberite barem jedan servis.', git_w_dirs: 'Ciljni direktoriji', git_w_roots: 'Korijeni na serveru', git_w_root_missing: 'ne postoji', git_w_root_ro: 'nije zapisiv',
  git_w_files: '{n} datoteka, {kb} KB', git_w_check: 'Provjeri cilj', git_w_checks: 'Provjere prije instalacije', git_w_free: 'slobodno {free} MB, potrebno {need} MB', git_w_tools: 'Alati',
  git_w_missing_tool: 'nema', git_w_existing: 'Ovdje postoji instalacija ({n} datoteka{v}).', git_w_overwrite: 'Zamijeni je (zamijenjene datoteke najprije se spremaju u .deploy-bak)',
  git_w_protected: 'Datoteke po hostu', git_w_protected_h: 'Zaštićene datoteke nikad se ne uzimaju iz repozitorija. Popunite ih iz predloška, kopirajte iz druge instalacije ili upišite.',
  git_w_no_slots: 'Katalog ne navodi zaštićene datoteke za ovaj servis.', git_w_m_template: 'Iz predloška', git_w_m_copy: 'Kopiraj iz', git_w_m_manual: 'Upiši sadržaj', git_w_m_skip: 'Izostavi',
  git_w_tmpl_na: 'nije dostupno', git_w_in_repo: 'u repozitoriju', git_w_typed: 'Upisana vrijednost', git_w_vault_pw: '🔑 {n}: lozinka', git_w_vault_user: '🔑 {n}: korisničko ime',
  git_w_vault_h: 'Vrijednosti iz trezora čitaju se samo na strani poslužitelja i ne spremaju se uz izvođenje.',
  git_w_unit: 'Izradi jedinicu servisa (neobavezno)', git_w_unit_kind: 'Vrsta', git_w_unit_name: 'Naziv', git_w_unit_user: 'Pokreni kao korisnik (prazno = zadano)', git_w_unit_cmd: 'Naredba (apsolutna putanja)',
  git_w_unit_dir: 'Radni direktorij', git_w_unit_enable: 'Pokreni pri podizanju sustava (systemctl enable)',
  git_w_unit_warn: 'Datoteka jedinice zapisuje se na server kao root (ili sa sudo -n) i nikad ne prepisuje postojeću. WRM ne testira naredbu: provjerite je.',
  git_w_start: 'Pokreni servis nakon instalacije (ponovno pokretanje njegovih jedinica)', git_w_start_h: 'Pokreću se jedinice koje spominju direktorij (i postojeće i nova); provjerava se njihovo stanje.',
  git_w_blocked: 'Najprije riješite gornje probleme.', git_w_confirm: '{kind}: {n} servis(a) na {server}', git_w_type: 'Upišite naziv servera {name} za potvrdu (produkcija ili zamjena)',
  git_w_source: 'Izvor', git_w_excluded: '{n} putanja se ne kopira (zapisi, priručne datoteke, sigurnosne kopije, zanemareni direktoriji)', git_w_tr_h: 'Kopira instalaciju s njezinom konfiguracijom po hostu. Virtualna okruženja i node_modules se ne kopiraju: izradite ih ponovno na cilju.',
  git_w_units_not_copied: 'Jedinice izvora se ne kopiraju: {list}', git_w_same: 'Odaberite drugi server ili direktorij od izvora.', git_w_env: 'Okruženje',
  git_step_precheck: 'Provjere', git_step_fetch: 'Sadržaj', git_step_receive: 'Kopiranje na server', git_step_register: 'Upis', git_step_unit: 'Jedinica', git_step_source: 'Izvor',
});

Object.assign(LANGS.en, {
  git_tab_gitignore: '.gitignore', gi_service: 'Service', gi_load: 'Load', gi_loading: 'Reading the repository and the latest checks…',
  gi_repo: 'Repository', gi_file: 'File', gi_branch: 'Branch', gi_current: 'Current .gitignore', gi_no_current: 'The repository has no .gitignore here yet.',
  gi_warnings: 'Committed files that look like secrets or per-host configuration', gi_reason_secret: 'looks like a secret', gi_reason_protected: 'matches a protected glob',
  gi_warn_advice: 'Commit a template {tmpl} instead (WRM fills it on install), add {path} to .gitignore and remove it from the repository with git rm --cached {path}.',
  gi_extra: 'Files only on servers (state extra, from the latest checks)', gi_extra_none: 'No extra files in the latest checks of this service.', gi_size: 'Size', gi_pattern: 'Pattern',
  gi_stacks: 'Standard patterns', gi_detected: 'found: {f}', gi_protected: 'Protected globs (per-host files)', gi_custom: 'More patterns', gi_custom_h: 'One per line, e.g. /data/ or *.local.ini',
  gi_lists: 'Catalog lists of this service', gi_lists_h: 'Exclude: never compared or deployed. Protected: per-host files, never overwritten.', gi_not_catalog: 'This service comes from a bundle: add it to the catalog to edit its lists.',
  gi_preview: 'Preview', gi_download: 'Download .gitignore', gi_mr: 'Create merge request…', gi_mr_gh: 'Create pull request…', gi_nothing: 'Nothing to add: the .gitignore already has these patterns.',
  gi_added: '{n} patterns added', gi_diff: 'Current → new', gi_mr_title: 'Merge request to {repo}', gi_mr_q: 'WRM creates the branch {branch}, commits {file} with these patterns and opens a merge request against {base} in:',
  gi_mr_type: 'Type the repository name {name} to confirm', gi_mr_go: 'Create', gi_mr_done: 'Created: {url}', gi_mr_default: 'The default output is a download; WRM writes to the Git server only when you choose this.',
  gi_note: 'Note',
  git_write_token: 'Write token (optional)', git_write_token_h: 'Only for .gitignore merge / pull requests. Leave empty to keep the stored one.',
  git_webhook: 'Webhook…', git_wh_title: 'Incoming webhook: {name}', git_wh_h: 'A push, tag or release in GitLab / GitHub, or a POST from a CI job, starts an immediate check of the affected services instead of waiting for the interval.',
  git_wh_off: 'The webhook is off.', git_wh_enable: 'Turn on', git_wh_renew: 'New secret', git_wh_disable: 'Turn off', git_wh_url: 'URL', git_wh_secret: 'Secret (shown only now: copy it)',
  git_wh_last: 'Last delivery', git_wh_gitlab: 'GitLab: Settings → Webhooks, the URL, the secret as "Secret token", events "Push" and "Tag push".',
  git_wh_github: 'GitHub: Settings → Webhooks, the URL, content type application/json, the secret, events "push" and "release".',
  git_wh_generic: 'Jenkins or any CI: POST with the header X-WRM-Token: <secret>, optionally a JSON body {"project": "group/name"} or {"services": ["name"]}.',
  git_feeds: 'Bundles from CI artifacts', git_feeds_h: 'Fetch the newest bundle from a URL (GitLab job artifacts API, Jenkins artifact URL or any HTTPS URL) and use it as the offline source.',
  git_add_feed: 'Add artifact URL', git_feed_url: 'URL of the bundle (.tar.gz)', git_feed_header: 'Auth header name', git_feed_value: 'Auth header value', git_feed_value_h: 'PRIVATE-TOKEN: <token> for GitLab, Authorization: Basic user:token for Jenkins. Leave empty to keep the stored value.',
  git_feed_interval: 'Fetch every (minutes, 0 = only by hand)', git_feed_fetch: 'Fetch now', git_feed_imported: 'Bundle imported', git_feed_unchanged: 'Unchanged since the last fetch', git_feed_last: 'Last fetched',
  git_deploy_method: 'Deploy method', git_dm_files: 'Files over SSH (WRM writes the files)', git_dm_jenkins: 'CI pipeline: Jenkins job', git_dm_gitlab: 'CI pipeline: GitLab pipeline',
  git_dm_h: 'For services built elsewhere: an update or upgrade triggers the job with the parameters server, install_path and version, and WRM follows its status.',
  git_dm_job: 'Job URL', git_dm_user: 'Jenkins user', git_dm_token: 'API token (or the job\'s trigger token without a user)', git_dm_gl_url: 'GitLab address (empty = the Git source)', git_dm_project: 'Project (empty = the service\'s project)',
  git_dm_ref: 'Branch to run (empty = the target\'s branch)', git_dm_trigger: 'Pipeline trigger token', git_dm_keep: 'leave empty to keep the stored token',
  git_ci_item: 'Deployed by CI ({kind}): WRM triggers the job with server, install_path and version and follows its status; no files are written.', git_ci_link: 'Open in CI', git_step_ci: 'CI pipeline',
});
Object.assign(LANGS.hr, {
  git_tab_gitignore: '.gitignore', gi_service: 'Servis', gi_load: 'Učitaj', gi_loading: 'Čitanje repozitorija i zadnjih provjera…',
  gi_repo: 'Repozitorij', gi_file: 'Datoteka', gi_branch: 'Grana', gi_current: 'Trenutni .gitignore', gi_no_current: 'Repozitorij ovdje još nema .gitignore.',
  gi_warnings: 'Datoteke u repozitoriju koje izgledaju kao tajne ili konfiguracija po hostu', gi_reason_secret: 'izgleda kao tajna', gi_reason_protected: 'odgovara zaštićenom globu',
  gi_warn_advice: 'Umjesto nje spremite predložak {tmpl} (WRM ga popunjava pri instalaciji), dodajte {path} u .gitignore i uklonite je iz repozitorija s git rm --cached {path}.',
  gi_extra: 'Datoteke samo na serverima (stanje višak, iz zadnjih provjera)', gi_extra_none: 'U zadnjim provjerama ovog servisa nema viška datoteka.', gi_size: 'Veličina', gi_pattern: 'Uzorak',
  gi_stacks: 'Standardni uzorci', gi_detected: 'pronađeno: {f}', gi_protected: 'Zaštićeni globovi (datoteke po hostu)', gi_custom: 'Dodatni uzorci', gi_custom_h: 'Jedan po retku, npr. /data/ ili *.local.ini',
  gi_lists: 'Popisi ovog servisa u katalogu', gi_lists_h: 'Isključi: nikad se ne uspoređuje ni ne isporučuje. Zaštićeno: datoteke po hostu, nikad se ne prepisuju.', gi_not_catalog: 'Ovaj servis dolazi iz bundlea: dodajte ga u katalog da biste uređivali njegove popise.',
  gi_preview: 'Pregled', gi_download: 'Preuzmi .gitignore', gi_mr: 'Izradi merge request…', gi_mr_gh: 'Izradi pull request…', gi_nothing: 'Nema ničega za dodati: .gitignore već ima ove uzorke.',
  gi_added: 'dodano uzoraka: {n}', gi_diff: 'Trenutno → novo', gi_mr_title: 'Merge request u {repo}', gi_mr_q: 'WRM izrađuje granu {branch}, sprema {file} s ovim uzorcima i otvara merge request prema {base} u:',
  gi_mr_type: 'Upišite naziv repozitorija {name} za potvrdu', gi_mr_go: 'Izradi', gi_mr_done: 'Izrađeno: {url}', gi_mr_default: 'Zadani izlaz je preuzimanje; WRM piše na Git server samo kad odaberete ovo.',
  gi_note: 'Napomena',
  git_write_token: 'Token za pisanje (neobavezno)', git_write_token_h: 'Samo za merge / pull requestove za .gitignore. Ostavite prazno za zadržavanje spremljenog.',
  git_webhook: 'Webhook…', git_wh_title: 'Dolazni webhook: {name}', git_wh_h: 'Push, tag ili release u GitLabu / GitHubu, ili POST iz CI posla, odmah pokreće provjeru zahvaćenih servisa umjesto čekanja intervala.',
  git_wh_off: 'Webhook je isključen.', git_wh_enable: 'Uključi', git_wh_renew: 'Nova tajna', git_wh_disable: 'Isključi', git_wh_url: 'URL', git_wh_secret: 'Tajna (prikazuje se samo sada: kopirajte je)',
  git_wh_last: 'Zadnja isporuka', git_wh_gitlab: 'GitLab: Settings → Webhooks, URL, tajna kao "Secret token", događaji "Push" i "Tag push".',
  git_wh_github: 'GitHub: Settings → Webhooks, URL, content type application/json, tajna, događaji "push" i "release".',
  git_wh_generic: 'Jenkins ili bilo koji CI: POST sa zaglavljem X-WRM-Token: <tajna>, neobavezno JSON tijelo {"project": "grupa/naziv"} ili {"services": ["naziv"]}.',
  git_feeds: 'Bundleovi iz CI artefakata', git_feeds_h: 'Dohvati najnoviji bundle s URL-a (GitLab API za artefakte posla, URL artefakta u Jenkinsu ili bilo koji HTTPS URL) i koristi ga kao offline izvor.',
  git_add_feed: 'Dodaj URL artefakta', git_feed_url: 'URL bundlea (.tar.gz)', git_feed_header: 'Naziv zaglavlja za prijavu', git_feed_value: 'Vrijednost zaglavlja za prijavu', git_feed_value_h: 'PRIVATE-TOKEN: <token> za GitLab, Authorization: Basic korisnik:token za Jenkins. Ostavite prazno za zadržavanje spremljene vrijednosti.',
  git_feed_interval: 'Dohvaćaj svakih (minuta, 0 = samo ručno)', git_feed_fetch: 'Dohvati sada', git_feed_imported: 'Bundle uvezen', git_feed_unchanged: 'Nepromijenjeno od zadnjeg dohvata', git_feed_last: 'Zadnji dohvat',
  git_deploy_method: 'Način isporuke', git_dm_files: 'Datoteke preko SSH-a (WRM zapisuje datoteke)', git_dm_jenkins: 'CI pipeline: Jenkins posao', git_dm_gitlab: 'CI pipeline: GitLab pipeline',
  git_dm_h: 'Za servise koji se grade drugdje: ažuriranje ili nadogradnja pokreće posao s parametrima server, install_path i version, a WRM prati njegovo stanje.',
  git_dm_job: 'URL posla', git_dm_user: 'Jenkins korisnik', git_dm_token: 'API token (ili token za okidanje posla bez korisnika)', git_dm_gl_url: 'Adresa GitLaba (prazno = Git izvor)', git_dm_project: 'Projekt (prazno = projekt servisa)',
  git_dm_ref: 'Grana za pokretanje (prazno = grana cilja)', git_dm_trigger: 'Token za okidanje pipelinea', git_dm_keep: 'ostavite prazno za zadržavanje spremljenog tokena',
  git_ci_item: 'Isporučuje CI ({kind}): WRM pokreće posao s parametrima server, install_path i version i prati njegovo stanje; datoteke se ne zapisuju.', git_ci_link: 'Otvori u CI-ju', git_step_ci: 'CI pipeline',
});

Object.assign(LANGS.en, {
  git_all_services: 'All services', git_all_servers: 'All servers', git_search_srv: 'Search server, host, path…',
  git_check_selected: 'Check selected ({n})', git_check_visible: 'Check visible ({n})', git_check_stale: 'Check only stale', git_check_cell: 'Check this cell',
  git_stale_never: 'never checked', git_stale_60: 'older than 1 hour', git_stale_1440: 'older than 24 hours', git_stale_10080: 'older than 7 days',
  git_sel_row: 'Select every installation of this server', git_sel_col: 'Select every installation of this service', git_sel_all: 'Select every visible installation',
  git_clear_sel: 'Clear selection', git_cancel_check: 'Cancel', git_job_progress: 'Checking {done}/{total}', git_job_targets: 'Refreshing targets…',
  git_job_done: 'Check finished: {n} installations.', git_job_cancelled: 'Check cancelled.', git_job_errors: '{n} could not be checked.', git_nothing_stale: 'No visible installation is that old.',
  git_q_queued: 'queued', git_q_running: 'checking…', git_checked_ago: 'checked {t}', git_ago_now: 'just now', git_ago_min: '{n} min ago', git_ago_h: '{n} h ago', git_ago_d: '{n} d ago',
  git_last_good: 'last good: {s}, {t}',
  git_files_all: 'Files in the folder and in Git', git_ff_all: 'All', git_ff_diff: 'Differing', git_ff_tracked: 'Tracked only', git_ff_untracked: 'Not tracked', git_files_search: 'Search files…',
  git_files_loading: 'Reading the folder on the server…', git_files_count: '{n} of {total}', git_files_truncated: 'The folder has more than {n} entries: only the first {n} are listed.',
  git_files_partial: 'Only the catalog\'s files of the repository are known: refresh the targets to see the other repository files too.',
  git_files_notcmp: '{n} repository files outside the catalog were not compared yet (Diff compares them).', git_files_more: 'Showing the first {n} rows: narrow the list with the filters.',
  git_files_failed: 'The folder could not be read: {e}', git_files_stored: 'Files of the last check', git_files_info: 'Files the catalog does not track are informational: they never change the state of the installation.',
  git_size: 'Size', git_mtime: 'Modified', git_reason: 'Why', git_open: 'Open',
  git_fs_same: 'same', git_fs_differs: 'differs', git_fs_unknown: 'not compared', git_fs_symlink: 'symlink', git_fs_unreadable: 'unreadable',
  git_r_tracked: 'tracked', git_r_excluded: 'excluded by the catalog', git_r_excluded_p: 'excluded by the catalog ({p})', git_r_not_included: 'not in the include list',
  git_r_protected: 'protected', git_r_protected_p: 'protected ({p})', git_r_extra: 'not in Git', git_r_tool: 'written by the deploy tool', git_r_compiled: 'compiled file',
  git_r_not_tracked: 'in Git, not tracked', git_r_no_target: 'no target yet', git_link_outside: 'points out of the installation (not followed)', git_dir_unreadable: 'directory: permission denied',
  git_informational: 'Excluded by the catalog — informational, not updated.', git_info_untracked: 'Not tracked by the catalog — informational, not updated.',
  git_view_unified: 'Unified', git_view_split: 'Side by side', git_from_server: 'Server', git_from_git: 'Git',
  git_binary_same: 'Binary file: the same on the server and in Git (by hash).', git_binary_differs: 'Binary file: differs between the server and Git (by hash).',
  git_big_same: 'The file is too large to show: the same on the server and in Git (by hash).', git_big_differs: 'The file is too large to show: it differs between the server and Git (by hash).',
  git_missing_srv: 'The file is not on the server: everything below would be added.', git_view_binary: 'Binary file, {s}, SHA-256 {h}', git_view_big: 'The file is too large to show ({s}), SHA-256 {h}',
  git_crlf: 'Line endings are ignored (CRLF = LF).', git_view_side: 'Open from', git_target_err: 'Target: {e}',
  git_commits: 'Commits', git_commits_behind: 'Commits behind: {n}', git_commits_none: 'The server is at the target commit.', git_commits_loading: 'Loading commits…',
  git_commits_range: '{f} → {t}', git_commits_more: 'The newest {n} are shown.',
});
Object.assign(LANGS.hr, {
  git_all_services: 'Svi servisi', git_all_servers: 'Svi serveri', git_search_srv: 'Traži server, host, putanju…',
  git_check_selected: 'Provjeri odabrane ({n})', git_check_visible: 'Provjeri prikazane ({n})', git_check_stale: 'Provjeri samo zastarjele', git_check_cell: 'Provjeri ovu ćeliju',
  git_stale_never: 'nikad provjerene', git_stale_60: 'starije od 1 sata', git_stale_1440: 'starije od 24 sata', git_stale_10080: 'starije od 7 dana',
  git_sel_row: 'Odaberi sve instalacije ovog servera', git_sel_col: 'Odaberi sve instalacije ovog servisa', git_sel_all: 'Odaberi sve prikazane instalacije',
  git_clear_sel: 'Poništi odabir', git_cancel_check: 'Prekini', git_job_progress: 'Provjera {done}/{total}', git_job_targets: 'Osvježavanje ciljeva…',
  git_job_done: 'Provjera završena: {n} instalacija.', git_job_cancelled: 'Provjera prekinuta.', git_job_errors: 'Nije moguće provjeriti: {n}.', git_nothing_stale: 'Nijedna prikazana instalacija nije toliko stara.',
  git_q_queued: 'na čekanju', git_q_running: 'provjera…', git_checked_ago: 'provjereno {t}', git_ago_now: 'upravo', git_ago_min: 'prije {n} min', git_ago_h: 'prije {n} h', git_ago_d: 'prije {n} d',
  git_last_good: 'zadnje uspješno: {s}, {t}',
  git_files_all: 'Datoteke u direktoriju i u Gitu', git_ff_all: 'Sve', git_ff_diff: 'Različite', git_ff_tracked: 'Samo praćene', git_ff_untracked: 'Nepraćene', git_files_search: 'Traži datoteke…',
  git_files_loading: 'Čitanje direktorija na serveru…', git_files_count: '{n} od {total}', git_files_truncated: 'Direktorij ima više od {n} stavki: prikazano je samo prvih {n}.',
  git_files_partial: 'Poznate su samo datoteke repozitorija iz kataloga: osvježite ciljeve da biste vidjeli i ostale datoteke repozitorija.',
  git_files_notcmp: 'Datoteke repozitorija izvan kataloga koje još nisu uspoređene: {n} (usporedite ih gumbom Razlike).', git_files_more: 'Prikazano je prvih {n} redaka: suzite popis filtrima.',
  git_files_failed: 'Direktorij nije moguće pročitati: {e}', git_files_stored: 'Datoteke zadnje provjere', git_files_info: 'Datoteke koje katalog ne prati samo su informativne: nikad ne mijenjaju stanje instalacije.',
  git_size: 'Veličina', git_mtime: 'Izmijenjeno', git_reason: 'Zašto', git_open: 'Otvori',
  git_fs_same: 'isto', git_fs_differs: 'razlikuje se', git_fs_unknown: 'nije uspoređeno', git_fs_symlink: 'simbolička veza', git_fs_unreadable: 'nečitljivo',
  git_r_tracked: 'praćeno', git_r_excluded: 'isključeno katalogom', git_r_excluded_p: 'isključeno katalogom ({p})', git_r_not_included: 'nije na popisu uključenih',
  git_r_protected: 'zaštićeno', git_r_protected_p: 'zaštićeno ({p})', git_r_extra: 'nije u Gitu', git_r_tool: 'zapisuje alat za isporuku', git_r_compiled: 'prevedena datoteka',
  git_r_not_tracked: 'u Gitu, nije praćeno', git_r_no_target: 'još nema cilja', git_link_outside: 'vodi izvan instalacije (ne slijedi se)', git_dir_unreadable: 'direktorij: pristup odbijen',
  git_informational: 'Isključeno katalogom — informativno, ne ažurira se.', git_info_untracked: 'Katalog ovo ne prati — informativno, ne ažurira se.',
  git_view_unified: 'Objedinjeno', git_view_split: 'Usporedno', git_from_server: 'Server', git_from_git: 'Git',
  git_binary_same: 'Binarna datoteka: ista na serveru i u Gitu (po hashu).', git_binary_differs: 'Binarna datoteka: razlikuje se između servera i Gita (po hashu).',
  git_big_same: 'Datoteka je prevelika za prikaz: ista na serveru i u Gitu (po hashu).', git_big_differs: 'Datoteka je prevelika za prikaz: razlikuje se između servera i Gita (po hashu).',
  git_missing_srv: 'Datoteke nema na serveru: sve ispod bilo bi dodano.', git_view_binary: 'Binarna datoteka, {s}, SHA-256 {h}', git_view_big: 'Datoteka je prevelika za prikaz ({s}), SHA-256 {h}',
  git_crlf: 'Završeci redaka se zanemaruju (CRLF = LF).', git_view_side: 'Otvori iz', git_target_err: 'Cilj: {e}',
  git_commits: 'Commitovi', git_commits_behind: 'Commitova iza cilja: {n}', git_commits_none: 'Server je na commitu cilja.', git_commits_loading: 'Učitavanje commitova…',
  git_commits_range: '{f} → {t}', git_commits_more: 'Prikazano je najnovijih {n}.',
});

(function () {
  const G = {state: null, tab: 'overview', f: {state: '', env: '', folder: '', tag: '', q: '', app: '', conn: ''}, onlyChanges: true, busy: false, sel: new Set()};
  const $ = id => document.getElementById(id);
  const tf = (k, vars) => { let s = t(k); for (const [a, b] of Object.entries(vars || {})) s = s.split('{' + a + '}').join(b); return s; };
  const lines = s => String(s || '').split(/[\n,]/).map(x => x.trim()).filter(Boolean);
  const json = (method, body) => ({method, headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body)});

  const css = `
#git-ws { position: fixed; inset: 0; z-index: 900; background: var(--bg); color: var(--text); display: flex; flex-direction: column; font-size: 13px; }
#git-ws .g-head { display: flex; align-items: center; gap: 10px; padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--bg2); flex-wrap: wrap; }
#git-ws .g-head h2 { font-size: 15px; margin: 0; white-space: nowrap; }
#git-ws .g-tabs { display: flex; gap: 4px; flex-wrap: wrap; }
#git-ws .g-tabs button { background: transparent; border: 1px solid transparent; color: var(--text2); padding: 6px 12px; border-radius: 8px; cursor: pointer; }
#git-ws .g-tabs button.on { background: var(--accent-d); border-color: var(--accent); color: var(--text); }
#git-ws .g-sp { flex: 1; }
#git-ws .g-body { flex: 1; overflow: auto; padding: 14px 16px 40px; }
#git-ws .g-bar { display: flex; gap: 6px; flex-wrap: wrap; align-items: center; margin: 0 0 12px; }
#git-ws .g-bar select, #git-ws .g-bar input { width: auto; min-width: 120px; flex: 0 1 auto; }
#git-ws .g-card { background: var(--bg2); border: 1px solid var(--border); border-radius: 10px; padding: 12px 14px; margin-bottom: 14px; }
#git-ws .g-card h3 { margin: 0 0 10px; font-size: 13px; }
#git-ws .g-scroll { overflow: auto; max-width: 100%; }
#git-ws table { border-collapse: collapse; width: 100%; font-size: 12px; }
#git-ws th { text-align: left; font-size: 10px; text-transform: uppercase; letter-spacing: .05em; color: var(--text3); padding: 6px 8px; border-bottom: 1px solid var(--border); white-space: nowrap; }
#git-ws td { padding: 6px 8px; border-bottom: 1px solid var(--border); vertical-align: top; }
#git-ws .g-matrix td.g-cell { min-width: 150px; }
#git-ws .g-inst { display: block; border-radius: 7px; padding: 4px 7px; margin: 2px 0; cursor: pointer; border: 1px solid var(--border2); background: var(--bg3); }
#git-ws .g-inst:hover { border-color: var(--accent); }
#git-ws .g-st { display: inline-block; font-size: 10.5px; padding: 1px 7px; border-radius: 10px; font-weight: 600; white-space: nowrap; }
#git-ws .s-ok { background: rgba(34,197,94,.14); color: var(--green); }
#git-ws .s-update, #git-ws .s-old, #git-ws .s-missing { background: rgba(245,158,11,.14); color: var(--yellow); }
#git-ws .s-review, #git-ws .s-modified, #git-ws .s-error { background: rgba(239,68,68,.14); color: var(--red); }
#git-ws .s-gone, #git-ws .s-unknown, #git-ws .s-extra, #git-ws .s-protected { background: var(--bg4); color: var(--text2); }
#git-ws .g-env { font-size: 10px; color: var(--purple); margin-left: 4px; }
#git-ws .g-mut { color: var(--text3); font-size: 11px; }
#git-ws .g-mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; word-break: break-all; }
#git-ws .g-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 10px; }
#git-ws .g-grid label { display: block; font-size: 11px; color: var(--text2); margin-bottom: 3px; }
#git-ws textarea { width: 100%; min-height: 54px; font-family: inherit; }
#git-ws .g-conns { max-height: 260px; overflow: auto; border: 1px solid var(--border); border-radius: 8px; padding: 6px; }
#git-ws .g-conns label { display: flex; gap: 6px; align-items: center; padding: 2px 4px; }
#git-ws .g-conns input { width: auto; }
#git-modal { position: fixed; inset: 0; z-index: 950; background: rgba(0,0,0,.55); display: flex; align-items: flex-start; justify-content: center; padding: 30px 12px; overflow: auto; }
#git-modal .g-dlg { background: var(--bg2); border: 1px solid var(--border2); border-radius: 12px; width: min(980px, 100%); padding: 16px 18px; }
#git-modal .g-dlg h3 { margin: 0 0 10px; font-size: 14px; word-break: break-all; }
#git-modal .g-diff { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; background: var(--bg); border: 1px solid var(--border); border-radius: 8px; overflow: auto; max-height: 60vh; }
#git-modal .g-diff div { white-space: pre; padding: 0 8px; }
#git-modal .g-diff .d-add { background: rgba(34,197,94,.13); color: #bbf7d0; }
#git-modal .g-diff .d-del { background: rgba(239,68,68,.13); color: #fecaca; }
#git-modal .g-diff .d-hunk { background: var(--accent-d); color: var(--text2); }
#git-modal table { border-collapse: collapse; width: 100%; font-size: 12px; }
#git-modal td, #git-modal th { padding: 4px 6px; border-bottom: 1px solid var(--border); text-align: left; }
#git-modal .g-st { display: inline-block; font-size: 10.5px; padding: 1px 7px; border-radius: 10px; font-weight: 600; }
#git-modal .s-ok { background: rgba(34,197,94,.14); color: var(--green); }
#git-modal .s-old, #git-modal .s-missing, #git-modal .s-update { background: rgba(245,158,11,.14); color: var(--yellow); }
#git-modal .s-modified, #git-modal .s-review, #git-modal .s-error { background: rgba(239,68,68,.14); color: var(--red); }
#git-modal .s-extra, #git-modal .s-protected, #git-modal .s-gone, #git-modal .s-unknown { background: var(--bg4); color: var(--text2); }
#git-modal .g-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 10px; }
#git-modal .g-grid label { display: block; font-size: 11px; color: var(--text2); margin-bottom: 3px; }
#git-modal textarea { width: 100%; min-height: 48px; font-family: inherit; }
#git-modal .g-btns { display: flex; gap: 8px; justify-content: flex-end; margin-top: 12px; flex-wrap: wrap; }
#git-ws input:not([type=checkbox]), #git-ws select, #git-ws textarea, #git-modal input:not([type=checkbox]), #git-modal select, #git-modal textarea {
  background: var(--bg3); color: var(--text); border: 1px solid var(--border2); border-radius: 7px; padding: 6px 9px; font-size: 12.5px; box-sizing: border-box; outline: none; }
#git-ws .g-grid input, #git-ws .g-grid select, #git-modal .g-grid input, #git-modal .g-grid select { width: 100%; }
#git-ws input:focus, #git-ws select:focus, #git-ws textarea:focus, #git-modal input:focus, #git-modal select:focus, #git-modal textarea:focus { border-color: var(--accent); }
#git-modal .g-mut { color: var(--text3); font-size: 11px; }
#git-modal .g-mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; word-break: break-all; }
#git-modal .g-env { font-size: 10px; color: var(--purple); margin-left: 6px; }
#git-modal .g-card { border: 1px solid var(--border); border-radius: 9px; padding: 10px 12px; margin: 10px 0; }
#git-modal .g-card h3 { margin: 0 0 6px; font-size: 12.5px; }
#git-modal .g-card > label:not(.g-opt), #git-modal .g-conf label { display: block; font-size: 11px; color: var(--text2); margin-bottom: 3px; }
#git-modal .g-opt > span { flex: 1 1 0; min-width: 160px; }
#git-modal .g-conf { margin-top: 10px; }
#git-modal .g-conf input { width: min(320px, 100%); }
#git-ws .g-pend, #git-modal .g-pend { font-size: 10px; color: var(--yellow); margin-left: 4px; white-space: nowrap; }
#git-ws .g-actions { display: flex; gap: 6px; flex-wrap: wrap; align-items: center; }
#git-modal .g-item { border: 1px solid var(--border); border-radius: 9px; padding: 10px 12px; margin: 10px 0; background: var(--bg3); }
#git-modal .g-item h4 { margin: 0 0 6px; font-size: 12.5px; word-break: break-all; }
#git-modal .g-prod { background: rgba(239,68,68,.16); color: var(--red); font-size: 10px; padding: 1px 6px; border-radius: 9px; font-weight: 700; text-transform: uppercase; margin-left: 4px; }
#git-modal .g-hideok .g-okrow { display: none; }
#git-modal .g-files { max-height: 32vh; overflow: auto; }
#git-modal .g-opt { display: flex; gap: 6px; align-items: center; margin: 4px 0; flex-wrap: wrap; }
#git-modal .g-opt input[type=radio], #git-modal .g-opt input[type=checkbox], #git-modal td input[type=checkbox], #git-modal td input[type=radio] { width: auto; }
#git-modal .g-log { font-size: 11.5px; margin: 6px 0 0; padding: 0; list-style: none; }
#git-modal .g-log li { padding: 2px 0; border-top: 1px dashed var(--border); word-break: break-word; white-space: pre-wrap; }
#git-modal .s-done, #git-modal .s-rs-ok { background: rgba(34,197,94,.14); color: var(--green); }
#git-modal .s-failed, #git-modal .s-rolled_back { background: rgba(239,68,68,.14); color: var(--red); }
#git-modal .s-scheduled, #git-modal .s-running, #git-modal .s-starting, #git-modal .s-pending, #git-modal .s-changed { background: rgba(245,158,11,.14); color: var(--yellow); }
#git-modal .s-cancelled, #git-modal .s-skipped, #git-modal .s-interrupted { background: var(--bg4); color: var(--text2); }
#git-ws .s-done, #git-ws .s-rs-ok { background: rgba(34,197,94,.14); color: var(--green); }
#git-ws .s-failed, #git-ws .s-rolled_back { background: rgba(239,68,68,.14); color: var(--red); }
#git-ws .s-scheduled, #git-ws .s-running, #git-ws .s-starting, #git-ws .s-pending { background: rgba(245,158,11,.14); color: var(--yellow); }
#git-ws .s-cancelled, #git-ws .s-skipped, #git-ws .s-interrupted { background: var(--bg4); color: var(--text2); }
#git-ws .g-diff { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; background: var(--bg); border: 1px solid var(--border); border-radius: 8px; overflow: auto; max-height: 50vh; }
#git-ws .g-diff div { white-space: pre; padding: 0 8px; }
#git-ws .g-diff .d-add { background: rgba(34,197,94,.13); color: var(--green); }
#git-ws .g-diff .d-del { background: rgba(239,68,68,.13); color: var(--red); }
#git-ws .g-diff .d-hunk { color: var(--text3); }
#git-ws .g-pats { display: flex; flex-wrap: wrap; gap: 4px 12px; }
#git-ws .g-pats label { display: inline-flex; gap: 5px; align-items: center; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; }
#git-ws .g-pats input, #git-ws td input[type=checkbox] { width: auto; }
#git-ws .g-warn { border-color: var(--red); }
#git-ws pre.g-pre { background: var(--bg); border: 1px solid var(--border); border-radius: 8px; padding: 8px; max-height: 30vh; overflow: auto; font-size: 11.5px; margin: 0; white-space: pre-wrap; word-break: break-all; }
@media (max-width: 640px) { #git-ws .g-head { padding: 8px 10px; } #git-ws .g-body { padding: 10px; } #git-ws .g-bar select, #git-ws .g-bar input { min-width: 0; flex: 1 1 140px; } }
#git-ws .g-cb, #git-modal .g-cb { width: auto; margin: 0 5px 0 0; vertical-align: middle; cursor: pointer; }
#git-ws .g-inst.sel { border-color: var(--accent); box-shadow: inset 0 0 0 1px var(--accent); }
#git-ws .g-inst .g-ago { font-size: 10.5px; color: var(--text3); }
#git-ws .g-cellbtn { background: transparent; border: 1px dashed var(--border2); color: var(--text2); border-radius: 6px; padding: 1px 7px; font-size: 11px; margin-top: 2px; cursor: pointer; }
#git-ws .g-cellbtn:hover { border-color: var(--accent); color: var(--text); }
#git-ws .g-cellbtn:disabled { opacity: .4; cursor: default; }
#git-ws .g-prog { display: inline-flex; gap: 6px; align-items: center; color: var(--text2); font-size: 12px; white-space: nowrap; }
.g-spin { display: inline-block; width: 10px; height: 10px; border: 2px solid var(--border2); border-top-color: var(--accent); border-radius: 50%; animation: g-spin .8s linear infinite; vertical-align: -1px; }
@keyframes g-spin { to { transform: rotate(360deg); } }
#git-ws .g-q { display: inline-flex; gap: 5px; align-items: center; font-size: 10.5px; color: var(--yellow); font-weight: 600; }
#git-ws .g-hostname { font-size: 10.5px; color: var(--text3); font-weight: 400; }
#git-ws .g-checkbar { padding: 6px 8px; border: 1px solid var(--border); border-radius: 9px; background: var(--bg2); }
#git-ws .g-checkbar select { min-width: 0; }
#git-modal .g-seg { display: inline-flex; border: 1px solid var(--border2); border-radius: 8px; overflow: hidden; flex-wrap: wrap; }
#git-modal .g-seg button { background: transparent; color: var(--text2); border: 0; border-right: 1px solid var(--border2); border-radius: 0; padding: 4px 10px; font-size: 12px; }
#git-modal .g-seg button:last-child { border-right: 0; }
#git-modal .g-seg button.on { background: var(--accent-d); color: var(--text); }
#git-modal .g-flist { max-height: 45vh; overflow: auto; border: 1px solid var(--border); border-radius: 8px; }
#git-modal .g-flist table td { vertical-align: top; }
#git-modal .g-flist .g-why { color: var(--text2); font-size: 11px; }
#git-modal .g-flist .g-pat { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
#git-modal .g-untracked td:first-child { opacity: .85; }
#git-modal .s-same { background: rgba(34,197,94,.10); color: var(--green); }
#git-modal .s-differs { background: rgba(139,92,246,.16); color: var(--purple); }
#git-modal .s-symlink, #git-modal .s-unknown { background: var(--bg4); color: var(--text2); }
#git-modal .s-unreadable { background: rgba(239,68,68,.10); color: var(--red); }
#git-modal .g-diff .d-add { color: var(--green); }
#git-modal .g-diff .d-del { color: var(--red); }
#git-modal .g-split { table-layout: fixed; width: 100%; border-collapse: collapse; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; }
#git-modal .g-split td { white-space: pre-wrap; word-break: break-all; padding: 0 6px; border: 0; vertical-align: top; }
#git-modal .g-split td.ln { width: 3.4em; color: var(--text3); text-align: right; user-select: none; padding: 0 4px; }
#git-modal .g-split td.d-del { background: rgba(239,68,68,.13); color: var(--red); }
#git-modal .g-split td.d-add { background: rgba(34,197,94,.13); color: var(--green); }
#git-modal .g-split td.d-none { background: var(--bg3); }
#git-modal .g-split tr.d-hunk td { background: var(--accent-d); color: var(--text2); }
#git-modal pre.g-code { background: var(--bg); border: 1px solid var(--border); border-radius: 8px; padding: 8px; max-height: 60vh; overflow: auto; font-size: 11.5px; margin: 0; white-space: pre; tab-size: 4; }
#git-modal .g-dhead { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; margin: 12px 0 6px; }
#git-modal .g-dhead h3 { margin: 0; flex: 1 1 200px; }
#git-modal .g-commits table td { white-space: normal; }
@media (max-width: 640px) { #git-modal { padding: 8px 4px; } #git-modal .g-dlg { padding: 12px 10px; } #git-modal .g-hm { display: none; } #git-modal .g-flist { max-height: 55vh; } }
`;

  function stBadge(s) { return `<span class="g-st s-${esc(s)}">${esc(t('git_st_' + s))}</span>`; }
  function fsBadge(s) { return `<span class="g-st s-${esc(s)}">${esc(t('git_fs_' + s))}</span>`; }

  async function load() {
    const r = await fetch('/api/git');
    if (!r.ok) { showToast(await apiError(r), 'error'); return false; }
    G.state = await r.json();
    return true;
  }

  function open() {
    if (!$('git-style')) { const s = document.createElement('style'); s.id = 'git-style'; s.textContent = css; document.head.appendChild(s); }
    let ws = $('git-ws');
    if (!ws) {
      ws = document.createElement('div');
      ws.id = 'git-ws';
      ws.innerHTML = `<div class="g-head"><h2>⎇ ${esc(t('git_title'))}</h2><div class="g-tabs" id="git-tabs"></div><span class="g-sp"></span>
        <span class="g-mut" id="git-last"></span><button id="git-check" onclick="gitWorkspace.check()">${esc(t('git_check_now'))}</button>
        <button class="btn-sec" onclick="gitWorkspace.close()" title="${esc(t('git_close'))}">✕</button></div><div class="g-body" id="git-body"></div>`;
      document.body.appendChild(ws);
      document.addEventListener('keydown', e => { if (e.key === 'Escape' && $('git-ws') && !$('git-modal')) close(); });
    }
    ws.style.display = 'flex';
    load().then(ok => { if (ok) { render(); resumeJob(); } });
  }

  function close() { const ws = $('git-ws'); if (ws) ws.style.display = 'none'; closeModal(); }

  function render() {
    const s = G.state;
    $('git-tabs').innerHTML = ['overview', 'services', 'installs', 'runs', 'gitignore', 'settings'].map(k =>
      `<button class="${G.tab === k ? 'on' : ''}" onclick="gitWorkspace.tab('${k}')">${esc(t('git_tab_' + k))}</button>`).join('');
    $('git-last').textContent = `${t('git_last_check')}: ${s.last_check_at ? fmtTime(s.last_check_at) : t('git_never')}`;
    $('git-check').disabled = !s.can_check || G.busy || jobRunning();
    $('git-check').title = s.can_check ? '' : t('git_no_checks');
    $('git-check').textContent = G.busy || jobRunning() ? t('git_checking') : t('git_check_now');
    const body = $('git-body');
    if (G.tab === 'overview') body.innerHTML = overviewHTML();
    else if (G.tab === 'services') body.innerHTML = servicesHTML();
    else if (G.tab === 'installs') body.innerHTML = installsHTML();
    else if (G.tab === 'runs') body.innerHTML = runsHTML();
    else if (G.tab === 'gitignore') body.innerHTML = ignoreHTML();
    else body.innerHTML = settingsHTML();
  }

  function setTab(k) { G.tab = k; render(); if (k === 'runs') loadRuns(); if (k === 'gitignore' && !I.info && !I.loading) ignoreLoad(); }

  // ── overview ──
  function filtered() {
    const f = G.f, q = f.q.toLowerCase();
    return G.state.installs.filter(i => (!f.state || i.state === f.state) && (!f.env || i.env === f.env) &&
      (!f.folder || String(i.folder_id || '') === f.folder) && (!f.tag || (i.tags || []).includes(f.tag)) &&
      (!f.app || i.app === f.app) && (!f.conn || String(i.conn_id) === f.conn) &&
      (!q || (i.conn_name + ' ' + (i.host || '') + ' ' + i.app + ' ' + i.path + ' ' + (i.env || '')).toLowerCase().includes(q)));
  }

  function filterBar() {
    const s = G.state, f = G.f;
    const envs = [...new Set(s.installs.map(i => i.env).filter(Boolean))].sort();
    const fids = [...new Set(s.installs.map(i => i.folder_id).filter(Boolean))];
    const tags = [...new Set(s.installs.flatMap(i => i.tags || []))].sort();
    const appNames = [...new Set(s.installs.map(i => i.app))].sort();
    const servers = [...new Map(s.installs.map(i => [i.conn_id, i.conn_name])).entries()].sort((a, b) => a[1].localeCompare(b[1]));
    const opt = (v, label, cur) => `<option value="${esc(v)}" ${String(cur) === String(v) ? 'selected' : ''}>${esc(label)}</option>`;
    const fname = id => ((folders || []).find(x => x.id === id) || {}).name || ('#' + id);
    return `<select onchange="gitWorkspace.filter('app', this.value)">${opt('', t('git_all_services'), f.app)}${appNames.map(x => opt(x, x, f.app)).join('')}</select>
      <select onchange="gitWorkspace.filter('conn', this.value)">${opt('', t('git_all_servers'), f.conn)}${servers.map(([id, n]) => opt(id, n, f.conn)).join('')}</select>
      <select onchange="gitWorkspace.filter('state', this.value)">${opt('', t('git_all_states'), f.state)}${['update', 'review', 'ok', 'error', 'gone', 'unknown'].map(x => opt(x, t('git_st_' + x), f.state)).join('')}</select>
      <select onchange="gitWorkspace.filter('env', this.value)">${opt('', t('git_all_envs'), f.env)}${envs.map(x => opt(x, x, f.env)).join('')}</select>
      <select onchange="gitWorkspace.filter('folder', this.value)">${opt('', t('git_all_folders'), f.folder)}${fids.map(x => opt(x, fname(x), f.folder)).join('')}</select>
      <select onchange="gitWorkspace.filter('tag', this.value)">${opt('', t('git_all_tags'), f.tag)}${tags.map(x => opt(x, x, f.tag)).join('')}</select>
      <input type="search" id="git-q" placeholder="${esc(t('git_search_srv'))}" value="${esc(f.q)}" oninput="gitWorkspace.filter('q', this.value, true)">`;
  }

  function refBar() {
    const st = G.state.settings;
    return `<label class="g-mut">${esc(t('git_ref_choice'))}</label>
      <select onchange="gitWorkspace.setRef(this.value)">${['catalog', 'default', 'branch'].map(m => `<option value="${m}" ${st.ref_mode === m ? 'selected' : ''}>${esc(t('git_ref_' + m))}</option>`).join('')}</select>
      ${st.ref_mode === 'branch' ? `<input id="git-ref-branch" value="${esc(st.ref_branch || '')}" placeholder="main" onchange="gitWorkspace.setRef('branch')" style="min-width:110px">` : ''}`;
  }

  function overviewHTML() {
    const s = G.state;
    const apps = [...s.catalog.apps.map(a => a.name), ...(s.bundle_apps || [])];
    const list = filtered();
    const targets = Object.fromEntries(s.targets.map(x => [x.app, x]));
    let html = `<div class="g-bar">${refBar()}<button class="btn-sec btn-sm" onclick="gitWorkspace.refresh()" ${s.can_check ? '' : 'disabled'}>↻ ${esc(t('git_refresh_targets'))}</button></div>`;
    html += `<div class="g-card"><h3>${esc(t('git_targets'))}</h3><div class="g-scroll"><table><tr><th>${esc(t('git_service'))}</th><th>${esc(t('git_target'))}</th><th>${esc(t('git_latest_tag'))}</th><th>${esc(t('git_source'))}</th><th>${esc(t('git_warnings'))}</th></tr>` +
      (apps.length ? apps.map(a => { const x = targets[a]; return `<tr><td><b>${esc(a)}</b></td><td>${x ? `<span class="g-mono">${esc(x.version || '—')}</span> <span class="g-mut">${esc(x.branch || '')} · ${esc(x.commit_date || '')}</span>` : `<span class="g-mut">${esc(t('git_no_target'))}</span>`}</td>
        <td class="g-mono">${esc((x && x.latest_tag) || '—')}</td><td class="g-mut">${esc(x ? x.source : '')}</td><td>${x && x.error ? `<span class="g-st s-error">${esc(x.error)}</span>` : esc(((x && x.warnings) || []).join('; '))}</td></tr>`; }).join('')
        : `<tr><td colspan="5" class="g-mut">${esc(t('git_empty_overview'))}</td></tr>`) + '</table></div></div>';
    html += `<div class="g-bar">${filterBar()}</div>${checkBar(list)}`;
    const servers = [...new Map(list.map(i => [i.conn_id, i])).entries()].sort((a, b) => a[1].conn_name.localeCompare(b[1].conn_name));
    if (!servers.length) return html + `<div class="hint-box">${esc(t('git_empty_overview'))}</div>`;
    const cols = [...apps, ...new Set(list.map(i => i.app))].filter((a, k, arr) => arr.indexOf(a) === k && list.some(i => i.app === a));
    const ids = l => l.map(i => i.id);
    const running = jobRunning();
    html += `<div class="g-card g-scroll"><table class="g-matrix"><tr><th>${selBox(ids(list), t('git_sel_all'))}${esc(t('git_server'))}</th>${cols.map(a => `<th>${selBox(ids(list.filter(i => i.app === a)), t('git_sel_col'))}${esc(a)}</th>`).join('')}</tr>` +
      servers.map(([cid, x]) => { const row = list.filter(i => i.conn_id === cid);
        return `<tr><td>${selBox(ids(row), t('git_sel_row'))}<b>${esc(x.conn_name)}</b>${x.host ? `<div class="g-hostname g-mono">${esc(x.host)}</div>` : ''}</td>${cols.map(a => { const cell = row.filter(i => i.app === a);
          return `<td class="g-cell">${cell.map(instChip).join('')}${cell.length && G.state.can_check ? `<button class="g-cellbtn" title="${esc(t('git_check_cell'))}" ${running ? 'disabled' : ''} onclick="gitWorkspace.checkIds([${ids(cell).join(',')}])">↻ ${esc(t('git_check_cell'))}</button>` : ''}</td>`; }).join('')}</tr>`; }).join('') + '</table></div>';
    return html;
  }

  function instChip(i) {
    const v = (i.version_md && i.version_md.version) || '';
    return `<span class="g-inst${G.sel.has(i.id) ? ' sel' : ''}" onclick="gitWorkspace.details(${i.id})" title="${esc(i.path)}"><input type="checkbox" class="g-cb" ${G.sel.has(i.id) ? 'checked' : ''} onclick="event.stopPropagation()" onchange="gitWorkspace.sel([${i.id}], this.checked)">${stateCell(i)}${i.env ? `<span class="g-env">${esc(i.env)}</span>` : ''}${i.restart_pending ? `<span class="g-pend" title="${esc(fmtTime(i.restart_pending))}">↻ ${esc(t('git_restart_pending'))}</span>` : ''}
      <div class="g-mut">${esc(t('git_version_md'))}: <span class="g-mono">${v ? esc(v) : esc(t('git_vmd_missing'))}</span>${i.behind ? ' · ' + esc(tf('git_behind', {n: i.behind})) : ''}</div>
      <div class="g-ago" title="${esc(i.checked_at ? fmtTime(i.checked_at) : '')}">${lastGood(i)}${esc(tf('git_checked_ago', {t: ago(i.checked_at)}))}</div></span>`;
  }

  // ── selection and check jobs ──
  const J = {job: null, items: {}, seq: 0, timer: null, stale: 'never'};
  const jobRunning = () => !!(J.job && J.job.state === 'running');

  function ago(ts) {
    if (!ts) return t('git_never');
    const s = Math.max(0, (Date.now() - new Date(ts).getTime()) / 1000);
    if (s < 60) return t('git_ago_now');
    if (s < 3600) return tf('git_ago_min', {n: Math.floor(s / 60)});
    if (s < 86400) return tf('git_ago_h', {n: Math.floor(s / 3600)});
    return tf('git_ago_d', {n: Math.floor(s / 86400)});
  }

  function lastGood(i) { return i.state === 'error' && i.last_ok_state ? esc(tf('git_last_good', {s: t('git_st_' + i.last_ok_state), t: ago(i.last_ok_at)})) + ' · ' : ''; }

  // the state badge, or the job's spinner while the installation is queued or being checked
  function stateCell(i) {
    const it = jobRunning() && J.items[i.id];
    if (it && it.state === 'running') return `<span class="g-q"><span class="g-spin"></span>${esc(t('git_q_running'))}</span>`;
    if (it && it.state === 'queued') return `<span class="g-q">⏳ ${esc(t('git_q_queued'))}</span>`;
    return stBadge(i.state);
  }

  function selBox(ids, title) {
    const on = ids.length > 0 && ids.every(id => G.sel.has(id));
    return `<input type="checkbox" class="g-cb" title="${esc(title)}" ${on ? 'checked' : ''} ${ids.length ? '' : 'disabled'} onclick="event.stopPropagation()" onchange="gitWorkspace.sel([${ids.join(',')}], this.checked)">`;
  }

  function checkBar(list) {
    const running = jobRunning(), can = G.state.can_check, n = G.sel.size;
    const j = J.job || {};
    const stale = ['never', '60', '1440', '10080'];
    return `<div class="g-bar g-checkbar">
      <button class="btn-sm" onclick="gitWorkspace.checkSel()" ${can && !running && n ? '' : 'disabled'}>✓ ${esc(tf('git_check_selected', {n}))}</button>
      <button class="btn-sec btn-sm" onclick="gitWorkspace.checkVisible()" ${can && !running && list.length ? '' : 'disabled'}>${esc(tf('git_check_visible', {n: list.length}))}</button>
      <span style="display:inline-flex;gap:4px;align-items:center"><button class="btn-sec btn-sm" onclick="gitWorkspace.checkStale()" ${can && !running && list.length ? '' : 'disabled'}>${esc(t('git_check_stale'))}</button>
      <select id="git-stale" onchange="gitWorkspace.setStale(this.value)">${stale.map(x => `<option value="${x}" ${J.stale === x ? 'selected' : ''}>${esc(t('git_stale_' + x))}</option>`).join('')}</select></span>
      ${n ? `<button class="btn-sec btn-sm" onclick="gitWorkspace.clearSel()">${esc(t('git_clear_sel'))}</button>` : ''}
      ${running ? `<span class="g-sp"></span><span class="g-prog"><span class="g-spin"></span>${esc(j.phase === 'targets' ? t('git_job_targets') : tf('git_job_progress', {done: j.done || 0, total: j.total || 0}))}</span>
        <button class="btn-danger btn-sm" onclick="gitWorkspace.cancelCheck()">✕ ${esc(t('git_cancel_check'))}</button>` : ''}</div>`;
  }

  // keepView re-renders the tab and keeps the focus, the caret and the scroll positions
  function keepView() {
    if (!$('git-ws') || $('git-ws').style.display === 'none') return;
    const ae = document.activeElement, id = ae && ae.id, pos = ae && typeof ae.selectionStart === 'number' ? ae.selectionStart : null;
    const scr = [...document.querySelectorAll('#git-body .g-scroll')].map(e => [e.scrollLeft, e.scrollTop]);
    const top = $('git-body') ? $('git-body').scrollTop : 0;
    render();
    document.querySelectorAll('#git-body .g-scroll').forEach((e, k) => { if (scr[k]) { e.scrollLeft = scr[k][0]; e.scrollTop = scr[k][1]; } });
    if ($('git-body')) $('git-body').scrollTop = top;
    if (id && $(id) && !$('git-modal')) { $(id).focus(); if (pos !== null) try { $(id).setSelectionRange(pos, pos); } catch (e) { /* not a text input */ } }
  }

  function sel(ids, on) { ids.forEach(id => on ? G.sel.add(id) : G.sel.delete(id)); keepView(); }
  function clearSel() { G.sel.clear(); keepView(); }
  function setStale(v) { J.stale = v; }

  function mergeInstalls(list) {
    const by = new Map(list.map(i => [i.id, i]));
    G.state.installs = G.state.installs.map(i => by.has(i.id) ? by.get(i.id) : i);
    const have = new Set(G.state.installs.map(i => i.id));
    list.forEach(i => { if (!have.has(i.id)) G.state.installs.push(i); });
  }

  function setJob(j) {
    J.job = j;
    J.items = {};
    (j.items || []).forEach(it => { J.items[it.install_id] = it; });
    if (j.installs) mergeInstalls(j.installs);
    J.seq = j.seq || 0;
  }

  async function startJob(body) {
    if (jobRunning()) return;
    const r = await fetch('/api/git/check/jobs', json('POST', body));
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    J.seq = 0;
    setJob(await r.json());
    keepView();
    pollJob();
  }

  function pollJob() {
    clearTimeout(J.timer);
    J.timer = setTimeout(async () => {
      let j = null;
      try {
        const r = await fetch('/api/git/check/jobs/current?since=' + J.seq);
        if (r.ok) j = (await r.json()).job;
      } catch (e) { /* network: try again */ }
      if (!j) { if (jobRunning()) J.timer = setTimeout(pollJob, 3000); return; }
      setJob(j);
      if (j.state === 'running') { keepView(); pollJob(); return; }
      const errs = (j.items || []).filter(x => x.state === 'error').length;
      showToast((j.state === 'cancelled' ? t('git_job_cancelled') : tf('git_job_done', {n: (j.items || []).filter(x => x.state === 'done').length})) +
        (errs ? ' ' + tf('git_job_errors', {n: errs}) : ''), errs ? 'warning' : 'success', 5000);
      if (await load()) keepView();
    }, 900);
  }

  async function resumeJob() {
    try {
      const r = await fetch('/api/git/check/jobs/current');
      if (!r.ok) return;
      const j = (await r.json()).job;
      if (j && j.state === 'running') { J.seq = 0; setJob(j); keepView(); pollJob(); }
    } catch (e) { /* no job */ }
  }

  function checkIds(ids) { if (ids.length) startJob({install_ids: ids}); }
  function checkSel() { if (!G.sel.size) { showToast(t('git_select_installs'), 'warning'); return; } checkIds([...G.sel]); }
  function checkVisible() { checkIds(filtered().map(i => i.id)); }
  function checkStale() {
    const ids = filtered().map(i => i.id);
    if (!ids.length) return;
    const body = J.stale === 'never' ? {install_ids: ids, only_never: true} : {install_ids: ids, stale_minutes: Number(J.stale)};
    const cutoff = J.stale === 'never' ? 0 : Date.now() - Number(J.stale) * 60e3;
    if (!filtered().some(i => !i.checked_at || (cutoff && new Date(i.checked_at).getTime() <= cutoff))) { showToast(t('git_nothing_stale'), 'info'); return; }
    startJob(body);
  }

  async function cancelCheck() {
    const r = await fetch('/api/git/check/jobs/current/cancel', {method: 'POST'});
    if (!r.ok) showToast(await apiError(r), 'error');
  }

  // ── services ──
  function servicesHTML() {
    const s = G.state;
    const targets = Object.fromEntries(s.targets.map(x => [x.app, x]));
    const row = (a, fromBundle) => { const x = targets[a.name || a]; const name = a.name || a;
      return `<tr><td><input type="checkbox" class="git-sel" value="${esc(name)}" ${x && !x.error ? '' : 'disabled'} style="width:auto"></td><td><b>${esc(name)}</b>${fromBundle ? ` <span class="g-mut">(${esc(t('git_bundle_only'))})</span>` : ''}</td>
        <td class="g-mono">${esc(a.project || (x && x.project) || '')}</td><td class="g-mono">${esc(a.ref || (x && x.ref) || '')}${a.subdir ? ' · ' + esc(a.subdir) + '/' : ''}</td>
        <td>${x ? `<span class="g-mono">${esc(x.version)}</span> <span class="g-mut">${x.files} ${esc(t('git_files').toLowerCase())}</span>` : '<span class="g-mut">—</span>'}</td>
        <td style="white-space:nowrap">${fromBundle ? '' : `<button class="btn-sec btn-sm" onclick="gitWorkspace.editApp('${esc(name)}')">${esc(t('git_edit'))}</button> <button class="btn-danger btn-sm" onclick="gitWorkspace.delApp('${esc(name)}')">✕</button>`}</td></tr>`; };
    return `<div class="g-bar"><button onclick="gitWorkspace.editApp('')">＋ ${esc(t('git_add_service'))}</button>
      <button class="btn-sec" onclick="gitWorkspace.suggest()">✨ ${esc(t('git_suggest'))}</button>
      <button class="btn-sec" onclick="document.getElementById('git-cat-file').click()">⬆ ${esc(t('git_import_catalog'))}</button><input type="file" id="git-cat-file" accept=".json,application/json" style="display:none" onchange="gitWorkspace.importCatalog(this)">
      <button class="btn-sec" onclick="gitWorkspace.exportCatalog()">⬇ ${esc(t('git_export_catalog'))}</button>
      <button class="btn-sec" onclick="gitWorkspace.exportBundle()">📦 ${esc(t('git_export_bundle'))}</button>
      <button class="btn-sec" onclick="gitWorkspace.refresh()" ${s.can_check ? '' : 'disabled'}>↻ ${esc(t('git_refresh_targets'))}</button></div>
      <div class="g-card g-scroll"><table><tr><th></th><th>${esc(t('git_f_name'))}</th><th>${esc(t('git_f_project'))}</th><th>${esc(t('git_f_ref'))}</th><th>${esc(t('git_target'))}</th><th></th></tr>
      ${s.catalog.apps.map(a => row(a, false)).join('')}${(s.bundle_apps || []).map(a => row(a, true)).join('')}</table></div>`;
  }

  function appForm(a) {
    const f = (k, label, v, ph) => `<div><label>${esc(t(label))}</label><input id="ga-${k}" value="${esc(v || '')}" placeholder="${esc(ph || '')}"></div>`;
    const ta = (k, label, v) => `<div><label>${esc(t(label))}</label><textarea id="ga-${k}" placeholder="${esc(t('git_lines_h'))}">${esc((v || []).join('\n'))}</textarea></div>`;
    return `<div class="g-grid">${f('name', 'git_f_name', a.name)}${f('project', 'git_f_project', a.project, 'group/project')}${f('ref', 'git_f_ref', a.ref || 'tag:latest', t('git_f_ref_h'))}
      ${f('branch', 'git_f_branch', a.branch, 'main')}${f('tag_filter', 'git_f_tag_filter', a.tag_filter, '^v\\d')}${f('subdir', 'git_f_subdir', a.subdir, 'linux')}
      <div><label>${esc(t('git_f_kind'))}</label><select id="ga-kind">${['app', 'library', 'tool'].map(k => `<option value="${k}" ${a.kind === k ? 'selected' : ''}>${esc(t('git_kind_' + k))}</option>`).join('')}</select></div></div>
      <div class="g-grid" style="margin-top:10px">${ta('fingerprint', 'git_f_fingerprint', a.fingerprint)}${ta('protected', 'git_f_protected', a.protected)}${ta('install_hint', 'git_f_hint', a.install_hint)}
      ${ta('include', 'git_f_include', a.include)}${ta('exclude', 'git_f_exclude', a.exclude)}</div>`;
  }

  function readAppForm() {
    const v = k => $('ga-' + k).value.trim();
    return {name: v('name'), project: v('project'), ref: v('ref'), branch: v('branch'), tag_filter: v('tag_filter'), subdir: v('subdir'), kind: v('kind'),
      fingerprint: lines(v('fingerprint')), protected: lines(v('protected')), install_hint: lines(v('install_hint')), include: lines(v('include')), exclude: lines(v('exclude'))};
  }

  function editApp(name, preset, files) {
    const a = preset || G.state.catalog.apps.find(x => x.name === name) || {kind: 'app', ref: 'tag:latest', protected: ['config*', '*.ini', '.env']};
    modal(`<h3>${esc(name || t('git_add_service'))}</h3>${appForm(a)}${name ? pipelineForm(name) : ''}${files ? `<details style="margin-top:10px"><summary class="g-mut">${esc(t('git_repo_files'))} (${files.length})</summary><div class="g-mono" style="max-height:200px;overflow:auto">${files.map(esc).join('<br>')}</div></details>` : ''}
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="ga-save">${esc(t('git_save'))}</button></div>`);
    $('ga-save').onclick = async () => {
      const app = readAppForm();
      const r = await fetch('/api/git/catalog/apps/' + encodeURIComponent(name || app.name), json('PUT', app));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = await r.json();
      if (name && !(await savePipeline(app.name))) return;
      closeModal(); render(); showToast(t('git_saved'), 'success');
    };
    if ($('gp-kind')) pipelineKind();
  }

  async function delApp(name) {
    if (!(await uiConfirm(tf('git_delete_q', {name}), {danger: true, okText: t('git_delete')}))) return;
    const r = await fetch('/api/git/catalog/apps/' + encodeURIComponent(name), {method: 'DELETE'});
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    G.state = await r.json(); render();
  }

  async function suggest() {
    const src = G.state.sources.find(x => x.kind !== 'bundle');
    let projects = [];
    if (src) { const r = await fetch(`/api/git/sources/${src.id}/test`, {method: 'POST'}); if (r.ok) projects = (await r.json()).projects || []; }
    modal(`<h3>✨ ${esc(t('git_suggest'))}</h3><div class="g-grid"><div><label>${esc(t('git_suggest_project'))}</label><input id="gs-project" list="gs-projects" placeholder="group/project">
      <datalist id="gs-projects">${projects.map(p => `<option value="${esc(p)}">`).join('')}</datalist></div>
      <div><label>${esc(t('git_suggest_ref'))}</label><input id="gs-ref"></div></div>
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gs-go">${esc(t('git_suggest_go'))}</button></div>`);
    $('gs-go').onclick = async () => {
      $('gs-go').disabled = true;
      const r = await fetch('/api/git/suggest', json('POST', {project: $('gs-project').value.trim(), ref: $('gs-ref').value.trim()}));
      $('gs-go').disabled = false;
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      const j = await r.json();
      editApp('', j.app, j.files);
    };
  }

  function importCatalog(input) {
    const file = input.files[0]; input.value = '';
    if (!file) return;
    file.text().then(async text => {
      let body;
      try { body = JSON.parse(text); } catch (e) { showToast(e.message, 'error'); return; }
      const merge = await uiConfirm(t('git_import_merge'), {okText: 'OK'});
      const r = await fetch('/api/git/catalog/import?mode=' + (merge ? 'merge' : 'replace'), json('POST', body));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = await r.json(); render(); showToast(t('git_saved'), 'success');
    });
  }

  async function download(r, fallback) {
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const m = /filename="([^"]+)"/.exec(r.headers.get('Content-Disposition') || '');
    const url = URL.createObjectURL(await r.blob());
    const a = document.createElement('a'); a.href = url; a.download = m ? m[1] : fallback; document.body.appendChild(a); a.click(); a.remove();
    setTimeout(() => URL.revokeObjectURL(url), 2000);
  }

  function exportCatalog() { fetch('/api/git/catalog/export').then(r => download(r, 'catalog.json')); }

  function exportBundle() {
    const apps = [...document.querySelectorAll('.git-sel:checked')].map(x => x.value);
    if (!apps.length) { showToast(t('git_export_bundle_none'), 'warning'); return; }
    fetch('/api/git/bundles/export', json('POST', {apps})).then(r => download(r, 'bundle.tar.gz'));
  }

  // ── installations ──
  function installsHTML() {
    const list = filtered();
    const sshConns = (conns || []).filter(c => c.protocol === 'SSH').sort((a, b) => a.name.localeCompare(b.name));
    const apps = [...G.state.catalog.apps.map(a => a.name), ...(G.state.bundle_apps || [])];
    return `<div class="g-bar">${filterBar()}<button class="btn-sec" onclick="gitWorkspace.check(true)" ${G.state.can_check && !jobRunning() ? '' : 'disabled'}>🔎 ${esc(t('git_discover'))}</button>
      <button class="btn-sm" onclick="gitWorkspace.install()" ${can('install') ? '' : `disabled title="${esc(t('git_no_permission'))}"`}>＋ ${esc(t('git_install'))}</button></div>
      ${checkBar(list)}<div class="g-bar g-actions">${actionButtons('')}</div>
      <div class="g-card g-scroll"><table><tr><th>${selBox(list.map(i => i.id), t('git_sel_all'))}</th><th>${esc(t('git_server'))}</th><th>${esc(t('git_service'))}</th><th>${esc(t('git_path'))}</th><th>${esc(t('git_env'))}</th><th>${esc(t('git_state'))}</th>
      <th>${esc(t('git_version_md'))}</th><th>${esc(t('git_target'))}</th><th>${esc(t('git_checked'))}</th><th></th></tr>
      ${list.map(i => `<tr><td><input type="checkbox" class="git-isel g-cb" value="${i.id}" ${G.sel.has(i.id) ? 'checked' : ''} onchange="gitWorkspace.sel([${i.id}], this.checked)"></td><td>${esc(i.conn_name)}${i.host ? `<div class="g-hostname g-mono">${esc(i.host)}</div>` : ''}</td><td><b>${esc(i.app)}</b>${i.restart_pending ? `<span class="g-pend" title="${esc(fmtTime(i.restart_pending))}">↻ ${esc(t('git_restart_pending'))}</span>` : ''}</td><td class="g-mono">${esc(i.path)}</td><td>${esc(i.env || '')}</td>
        <td>${stateCell(i)}${i.behind ? ` <span class="g-mut">${esc(tf('git_behind', {n: i.behind}))}</span>` : ''}${i.error ? `<div class="g-mut">${esc(i.error)}</div>` : ''}${i.state === 'error' && i.last_ok_state ? `<div class="g-mut">${lastGood(i).replace(/ · $/, '')}</div>` : ''}</td>
        <td class="g-mono">${esc((i.version_md && i.version_md.version) || t('git_vmd_missing'))}</td><td class="g-mono">${esc(i.target || '—')}</td><td class="g-mut" title="${esc(i.checked_at ? fmtTime(i.checked_at) : '')}">${esc(ago(i.checked_at))}</td>
        <td style="white-space:nowrap">${G.state.can_check ? `<button class="btn-sec btn-sm" title="${esc(t('git_check_cell'))}" ${jobRunning() ? 'disabled' : ''} onclick="gitWorkspace.checkIds([${i.id}])">↻</button> ` : ''}<button class="btn-sec btn-sm" onclick="gitWorkspace.details(${i.id})">${esc(t('git_details'))}</button> <button class="btn-danger btn-sm" onclick="gitWorkspace.forget(${i.id})">✕</button></td></tr>`).join('')}</table></div>
      <div class="g-card"><h3>${esc(t('git_add_install'))}</h3><div class="g-bar"><select id="gi-conn">${sshConns.map(c => `<option value="${c.id}">${esc(c.name)}</option>`).join('')}</select>
      <select id="gi-app">${apps.map(a => `<option>${esc(a)}</option>`).join('')}</select><input id="gi-path" placeholder="/opt/service" style="min-width:200px"><input id="gi-env" placeholder="${esc(t('git_env'))}">
      <button class="btn-sec" onclick="gitWorkspace.addInstall()">＋</button></div></div>`;
  }

  async function addInstall() {
    const r = await fetch('/api/git/installs', json('POST', {conn_id: Number($('gi-conn').value), app: $('gi-app').value, path: $('gi-path').value.trim(), env: $('gi-env').value.trim()}));
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    G.state = await r.json(); render();
  }

  async function forget(id) {
    if (!(await uiConfirm(t('git_forget_q'), {danger: true, okText: t('git_forget')}))) return;
    const r = await fetch('/api/git/installs/' + id, {method: 'DELETE'});
    if (r.ok) { G.state = await r.json(); render(); }
  }

  // ── installation details: every file, diff, view, commits ──
  const F = {id: 0, list: null, err: '', loading: false, filter: 'all', q: '', diff: null, mode: window.innerWidth < 700 ? 'unified' : 'split'};

  async function details(id) {
    const r = await fetch('/api/git/installs/' + id);
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const i = await r.json();
    G.detail = i;
    Object.assign(F, {id, list: null, err: '', loading: G.state.can_check, q: '', diff: null});
    const v = i.version_md || {};
    modal(`<h3>${esc(i.conn_name)} · ${esc(i.app)} · <span class="g-mono">${esc(i.path)}</span></h3>
      <div class="g-bar">${stBadge(i.state)}${i.env ? `<span class="g-env">${esc(i.env)}</span>` : ''}${Object.entries(i.counts || {}).map(([k, n]) => `${fsBadge(k)} ${n}`).join(' ')}
        <span class="g-mut" title="${esc(i.checked_at ? fmtTime(i.checked_at) : '')}">${lastGood(i)}${esc(tf('git_checked_ago', {t: ago(i.checked_at)}))}</span></div>
      <div class="g-bar g-actions">${actionButtons(i.id)}</div>
      ${i.error ? `<div class="hint-box bad">${esc(i.error)}</div>` : ''}${i.capped ? `<div class="hint-box warn">${esc(t('git_capped'))}</div>` : ''}
      <div class="g-grid"><div><label>${esc(t('git_version_md'))}</label><span class="g-mono">${v.version ? esc(v.version) : esc(t('git_vmd_missing'))}</span>
        ${v.ref || v.commit ? `<div class="g-mut">${[v.ref ? esc(v.ref) : '', v.commit ? `<span class="g-mono">${esc(v.commit)}</span>` : ''].filter(Boolean).join(' · ')}</div>` : ''}${v.updated ? `<div class="g-mut">${esc(v.updated)} ${esc(v.user || '')}${v.bundle ? ' · ' + esc(v.bundle) : ''}</div>` : ''}</div>
        <div><label>${esc(t('git_target'))}</label><span class="g-mono">${esc(i.target || '—')}</span><div id="git-commits" class="g-commits"></div></div>
        <div><label>${esc(t('git_units'))}</label>${(i.units || []).length ? i.units.map(u => `<div class="g-mono">${esc(u.kind)}: ${esc(u.name)}</div>`).join('') : `<span class="g-mut">${esc(t('git_no_units'))}</span>`}</div></div>
      <h3 style="margin:14px 0 6px">${esc(t('git_files_all'))}</h3><div id="git-files">${filesHTML()}</div>
      <div id="git-diff"></div><div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`);
    if (G.state.can_check) { loadFiles(id); loadCommits(i); }
  }

  async function loadFiles(id) {
    let list = null, err = '';
    try {
      const r = await fetch(`/api/git/installs/${id}/files`);
      if (r.ok) list = await r.json(); else err = await apiError(r);
    } catch (e) { err = e.message; }
    if (F.id !== id || !$('git-files')) return;
    Object.assign(F, {list, err, loading: false});
    paintFiles();
  }

  function paintFiles() {
    const el = $('git-files');
    if (!el) return;
    const ae = document.activeElement, focused = ae && ae.id === 'git-fq', pos = focused ? ae.selectionStart : 0;
    const box = el.querySelector('.g-flist'), top = box ? box.scrollTop : 0;
    el.innerHTML = filesHTML();
    const nb = el.querySelector('.g-flist');
    if (nb) nb.scrollTop = top;
    if (focused && $('git-fq')) { $('git-fq').focus(); $('git-fq').setSelectionRange(pos, pos); }
  }

  // the stored comparison of the last check (tracked and extra files)
  function storedHTML(i) {
    const files = (i.files || []).filter(f => !G.onlyChanges || (f.state !== 'ok' && f.state !== 'protected'));
    return `<label style="display:flex;gap:6px;align-items:center;margin:8px 0 6px"><input type="checkbox" style="width:auto" ${G.onlyChanges ? 'checked' : ''} onchange="gitWorkspace.toggleChanges(this.checked)">${esc(t('git_files_only'))}</label>
      <div class="g-flist"><table><tr><th>${esc(t('git_files_stored'))}</th><th>${esc(t('git_state'))}</th><th></th></tr>
      ${files.map(f => `<tr><td class="g-mono">${esc(f.path)}</td><td>${fsBadge(f.state)}${f.behind ? ` <span class="g-mut">${esc(tf('git_behind', {n: f.behind}))} (${esc(f.tag || '')} ${esc(f.date || '')})</span>` : ''}${f.note ? ` <span class="g-mut">${esc(f.note)}</span>` : ''}</td>
        <td>${f.diffable && f.state !== 'ok' && G.state.can_check ? `<button class="btn-sec btn-sm" onclick="gitWorkspace.diff(${i.id}, '${esc(encodeURIComponent(f.path))}')">${esc(t('git_diff'))}</button>` : ''}</td></tr>`).join('')}</table></div>`;
  }

  const DIFFERING = ['old', 'modified', 'differs', 'missing', 'extra'];

  function reasonText(f) {
    if (f.pattern && (f.reason === 'excluded' || f.reason === 'protected')) return tf('git_r_' + f.reason + '_p', {p: f.pattern});
    return t('git_r_' + f.reason);
  }

  function filesHTML() {
    const i = G.detail;
    if (!G.state.can_check) return storedHTML(i);
    if (F.loading) return `<div class="hint-box">⏳ ${esc(t('git_files_loading'))}</div>`;
    if (F.err || !F.list) return `<div class="hint-box bad">${esc(tf('git_files_failed', {e: F.err}))}</div>${storedHTML(i)}`;
    const L = F.list, q = F.q.toLowerCase();
    const pass = f => (F.filter === 'all' || (F.filter === 'diff' && DIFFERING.includes(f.state)) || (F.filter === 'tracked' && f.tracked) || (F.filter === 'untracked' && !f.tracked)) &&
      (!q || f.path.toLowerCase().includes(q));
    const rows = L.files.filter(pass), max = 1500;
    const seg = k => `<button class="${F.filter === k ? 'on' : ''}" onclick="gitWorkspace.fileFilter('${k}')">${esc(t('git_ff_' + k))}</button>`;
    const enc = p => esc(encodeURIComponent(p));
    const row = f => {
      const name = f.kind === 'symlink' ? `${esc(f.path)} <span class="g-mut">→ ${esc(f.link || '?')}</span>${f.outside ? `<div style="color:var(--red);font-family:inherit;font-size:11px">${esc(t('git_link_outside'))}</div>` : ''}`
        : f.kind === 'dir' ? `${esc(f.path)}/ <span class="g-mut">${esc(t('git_dir_unreadable'))}</span>` : esc(f.path);
      return `<tr class="${f.tracked ? '' : 'g-untracked'}"><td class="g-mono">${name}</td>
        <td class="g-mut g-hm" style="white-space:nowrap">${f.on_server && f.kind === 'file' ? esc(fmtBytes(f.size)) : f.size ? esc(fmtBytes(f.size)) : ''}</td>
        <td class="g-mut g-hm" style="white-space:nowrap">${f.mtime ? esc(fmtTime(f.mtime)) : ''}</td>
        <td>${fsBadge(f.state)}${f.behind ? ` <span class="g-mut">${esc(tf('git_behind', {n: f.behind}))}${f.tag ? ' (' + esc(f.tag) + ')' : ''}</span>` : ''}</td>
        <td class="g-why">${esc(reasonText(f))}</td>
        <td style="white-space:nowrap">${f.viewable ? `<button class="btn-sec btn-sm" onclick="gitWorkspace.viewFile('${enc(f.path)}', '')">${esc(t('git_open'))}</button> ` : ''}${f.diffable && f.state !== 'ok' && f.state !== 'same' ? `<button class="btn-sec btn-sm" onclick="gitWorkspace.diff(${i.id}, '${enc(f.path)}')">${esc(t('git_diff'))}</button>` : ''}</td></tr>`;
    };
    return `${L.target_error ? `<div class="hint-box warn">${esc(tf('git_target_err', {e: L.target_error}))}</div>` : ''}
      ${L.truncated ? `<div class="hint-box warn">${esc(tf('git_files_truncated', {n: 10000}))}</div>` : ''}
      ${L.git_partial && L.target ? `<div class="hint-box">${esc(t('git_files_partial'))}</div>` : ''}
      ${L.not_compared ? `<div class="g-mut">${esc(tf('git_files_notcmp', {n: L.not_compared}))}</div>` : ''}
      <div class="g-bar" style="margin:6px 0"><span class="g-seg">${['all', 'diff', 'tracked', 'untracked'].map(seg).join('')}</span>
        <input type="search" id="git-fq" placeholder="${esc(t('git_files_search'))}" value="${esc(F.q)}" oninput="gitWorkspace.fileSearch(this.value)" style="flex:1 1 160px;min-width:0">
        <span class="g-mut">${esc(tf('git_files_count', {n: rows.length, total: L.files.length}))}</span></div>
      <div class="g-mut" style="margin-bottom:6px">${esc(t('git_files_info'))}</div>
      <div class="g-flist"><table><tr><th>${esc(t('git_path'))}</th><th class="g-hm">${esc(t('git_size'))}</th><th class="g-hm">${esc(t('git_mtime'))}</th><th>${esc(t('git_state'))}</th><th>${esc(t('git_reason'))}</th><th></th></tr>
      ${rows.slice(0, max).map(row).join('')}</table></div>${rows.length > max ? `<div class="g-mut">${esc(tf('git_files_more', {n: max}))}</div>` : ''}`;
  }

  let fqTimer = null;
  function fileFilter(k) { F.filter = k; paintFiles(); }
  function fileSearch(v) { F.q = v; clearTimeout(fqTimer); fqTimer = setTimeout(paintFiles, 200); }

  async function loadCommits(i) {
    const el = $('git-commits');
    const v = i.version_md || {};
    if (!el || !i.target || !(v.commit || v.version)) return;
    el.innerHTML = `<div class="g-mut">⏳ ${esc(t('git_commits_loading'))}</div>`;
    let c = null, err = '';
    try {
      const r = await fetch(`/api/git/installs/${i.id}/commits`);
      if (r.ok) c = await r.json(); else err = await apiError(r);
    } catch (e) { err = e.message; }
    if (F.id !== i.id || !$('git-commits')) return;
    if (!c) { $('git-commits').innerHTML = `<div class="g-mut">${esc(t('git_commits'))}: ${esc(err)}</div>`; return; }
    $('git-commits').innerHTML = !c.count ? `<div class="g-mut">✓ ${esc(t('git_commits_none'))}</div>`
      : `<details><summary><b>${esc(tf('git_commits_behind', {n: c.count}))}</b> <span class="g-mut g-mono">${esc(tf('git_commits_range', {f: c.from, t: c.to}))}</span></summary>
        <div class="g-flist" style="max-height:30vh;margin-top:6px"><table>${c.commits.map(x => `<tr><td class="g-mono">${esc(x.sha.slice(0, 10))}</td><td>${esc(x.title)}</td><td class="g-mut g-hm">${esc(x.author)}</td><td class="g-mut g-hm" style="white-space:nowrap">${esc(x.date ? fmtTime(x.date) : '')}</td></tr>`).join('')}</table></div>
        ${c.commits.length < c.count ? `<div class="g-mut">${esc(tf('git_commits_more', {n: c.commits.length}))}</div>` : ''}</details>`;
  }

  async function diff(id, path) {
    const out = $('git-diff');
    out.innerHTML = `<div class="hint-box">⏳</div>`;
    out.scrollIntoView({block: 'nearest'});
    const r = await fetch(`/api/git/installs/${id}/diff?path=${path}`);
    if (!r.ok) { out.innerHTML = `<div class="hint-box bad">${esc(await apiError(r))}</div>`; return; }
    F.diff = await r.json();
    paintDiff();
  }

  function setMode(m) { F.mode = m; paintDiff(); }

  // splitRows pairs removed and added lines of a unified diff for the side-by-side view
  function splitRows(lines) {
    const rows = [];
    let del = [], add = [];
    const flush = () => { for (let k = 0; k < Math.max(del.length, add.length); k++) rows.push({l: del[k] || null, r: add[k] || null}); del = []; add = []; };
    for (const l of lines) {
      if (l.op === '-') del.push(l);
      else if (l.op === '+') add.push(l);
      else { flush(); rows.push(l.op === '@' ? {hunk: l} : {l, r: l, ctx: true}); }
    }
    flush();
    return rows;
  }

  function paintDiff() {
    const out = $('git-diff'), d = F.diff;
    if (!out || !d) return;
    const info = d.informational ? `<div class="hint-box warn">${esc(d.reason === 'excluded' || d.reason === 'not_included' ? t('git_informational') : t('git_info_untracked'))}${d.pattern ? ` <span class="g-mono">(${esc(d.pattern)})</span>` : ''}</div>` : '';
    const head = `<div class="g-dhead"><h3>${esc(t('git_diff_title'))}: <span class="g-mono">${esc(d.path)}</span> → ${esc(d.target || '')} ${fsBadge(d.state)}</h3>
      ${d.lines && d.lines.length ? `<span class="g-seg"><button class="${F.mode === 'unified' ? 'on' : ''}" onclick="gitWorkspace.diffMode('unified')">${esc(t('git_view_unified'))}</button><button class="${F.mode === 'split' ? 'on' : ''}" onclick="gitWorkspace.diffMode('split')">${esc(t('git_view_split'))}</button></span>` : ''}</div>`;
    let body;
    if (d.binary) body = `<div class="hint-box">${esc(t(d.same ? 'git_binary_same' : 'git_binary_differs'))}</div>`;
    else if (d.too_big) body = `<div class="hint-box">${esc(t(d.same ? 'git_big_same' : 'git_big_differs'))}</div>`;
    else if (!d.lines.length) body = `<div class="hint-box">${esc(t('git_no_diff'))}</div>`;
    else if (F.mode === 'split') {
      const cell = (l, side) => !l ? `<td class="ln"></td><td class="d-none"></td>` : `<td class="ln">${side === 'l' ? l.a || '' : l.b || ''}</td><td class="${l.op === '-' ? 'd-del' : l.op === '+' ? 'd-add' : ''}">${esc(l.s)}</td>`;
      body = `<div class="g-diff"><table class="g-split">${splitRows(d.lines).map(x => x.hunk ? `<tr class="d-hunk"><td class="ln"></td><td>@@ ${x.hunk.a}</td><td class="ln"></td><td>@@ ${x.hunk.b}</td></tr>` : `<tr>${cell(x.l, 'l')}${cell(x.r, 'r')}</tr>`).join('')}</table></div>`;
    } else {
      body = '<div class="g-diff">' + d.lines.map(l => l.op === '@' ? `<div class="d-hunk">@@ -${l.a} +${l.b} @@</div>` : `<div class="${l.op === '+' ? 'd-add' : l.op === '-' ? 'd-del' : ''}">${esc(l.op + ' ' + l.s)}</div>`).join('') + '</div>';
    }
    const notes = [d.missing_on_server ? t('git_missing_srv') : '', d.crlf_server !== d.crlf_target && d.lines && d.lines.length ? t('git_crlf') : ''].filter(Boolean);
    out.innerHTML = head + info + (notes.length ? `<div class="g-mut" style="margin-bottom:6px">${esc(notes.join(' '))}</div>` : '') + body;
    out.scrollIntoView({block: 'nearest'});
  }

  async function viewFile(path, side) {
    const out = $('git-diff'), id = F.id;
    if (!out) return;
    out.innerHTML = `<div class="hint-box">⏳</div>`;
    out.scrollIntoView({block: 'nearest'});
    const r = await fetch(`/api/git/installs/${id}/file?path=${path}${side ? '&side=' + side : ''}`);
    if (!r.ok) { out.innerHTML = `<div class="hint-box bad">${esc(await apiError(r))}</div>`; return; }
    const v = await r.json();
    const f = ((F.list || {}).files || []).find(x => x.path === v.path) || {};
    const other = v.side === 'server' ? (f.in_git && f.reason !== 'protected' ? 'git' : '') : (f.on_server ? 'server' : '');
    const head = `<div class="g-dhead"><h3>${esc(t('git_open'))}: <span class="g-mono">${esc(v.path)}</span> · ${esc(v.side === 'git' ? t('git_from_git') + (v.target ? ' ' + v.target : '') : t('git_from_server'))}</h3>
      ${other ? `<button class="btn-sec btn-sm" onclick="gitWorkspace.viewFile('${esc(encodeURIComponent(v.path))}', '${other}')">${esc(t('git_view_side'))} ${esc(t(other === 'git' ? 'git_from_git' : 'git_from_server'))}</button>` : ''}</div>`;
    let body;
    if (v.too_big) body = `<div class="hint-box">${esc(tf('git_view_big', {s: fmtBytes(v.size), h: v.hash || ''}))}</div>`;
    else if (v.binary) body = `<div class="hint-box">${esc(tf('git_view_binary', {s: fmtBytes(v.size), h: v.hash || ''}))}</div>`;
    else body = `<pre class="g-code">${esc(v.text)}</pre>`;
    out.innerHTML = head + body;
    out.scrollIntoView({block: 'nearest'});
  }

  // ── settings ──
  function settingsHTML() {
    const s = G.state, st = s.settings, cat = s.catalog;
    const apiSources = s.sources.filter(x => x.kind !== 'bundle'), bundles = s.sources.filter(x => x.kind === 'bundle');
    const sshConns = (conns || []).filter(c => c.protocol === 'SSH').sort((a, b) => a.name.localeCompare(b.name));
    const sel = new Set(st.conn_ids || []);
    return `<div class="g-card"><h3>${esc(t('git_sources'))}</h3>
      ${apiSources.length ? `<table>${apiSources.map(x => `<tr><td><b>${esc(x.name)}</b> <span class="g-mut">${esc(x.kind)}</span><div class="g-mono g-mut">${esc(x.url)}</div>${x.last_error ? `<div class="g-st s-error">${esc(x.last_error)}</div>` : ''}</td>
        <td style="white-space:nowrap"><button class="btn-sec btn-sm" onclick="gitWorkspace.testSource(${x.id})">${esc(t('git_test'))}</button> <button class="btn-sec btn-sm" ${s.can_check ? '' : 'disabled'} onclick="gitWorkspace.webhook(${x.id})">${x.hook_id ? '🔔 ' : ''}${esc(t('git_webhook'))}</button> <button class="btn-sec btn-sm" onclick="gitWorkspace.editSource(${x.id})">${esc(t('git_edit'))}</button> <button class="btn-danger btn-sm" onclick="gitWorkspace.delSource(${x.id})">✕</button></td></tr>`).join('')}</table>` : ''}
      <div class="g-bar" style="margin-top:8px"><button class="btn-sec" onclick="gitWorkspace.editSource(0)">＋ ${esc(t('git_add_source'))}</button></div></div>
      <div class="g-card"><h3>${esc(t('git_bundles'))}</h3>
      ${bundles.length ? `<table>${bundles.map(x => `<tr><td><b>${esc(x.bundle_id || x.name)}</b> <span class="g-mut">${esc(tf('git_bundle_age', {d: x.created || '?'}))} · ${esc(x.created_by || '')}</span><div class="g-mono g-mut">${esc(x.origin || '')}</div><div class="g-mut">${esc((x.apps || []).join(', '))}</div></td>
        <td><button class="btn-danger btn-sm" onclick="gitWorkspace.delSource(${x.id})">✕</button></td></tr>`).join('')}</table>` : ''}
      <div class="g-bar" style="margin-top:8px"><button class="btn-sec" onclick="document.getElementById('git-bundle-file').click()">⬆ ${esc(t('git_import_bundle'))}</button><input type="file" id="git-bundle-file" accept=".gz,.tgz,application/gzip" style="display:none" onchange="gitWorkspace.importBundle(this)"></div></div>
      ${feedsHTML()}
      <div class="g-card"><h3>${esc(t('git_source_mode'))}</h3><div class="g-grid">
        <div><label>${esc(t('git_source_mode'))}</label><select id="gx-mode">${['auto', 'live', 'bundle'].map(m => `<option value="${m}" ${st.mode === m ? 'selected' : ''}>${esc(t('git_mode_' + m))}</option>`).join('')}</select></div>
        <div><label>${esc(t('git_mode_bundle'))}</label><select id="gx-bundle"><option value="0">${esc(t('git_newest_bundle'))}</option>${bundles.map(b => `<option value="${b.id}" ${st.bundle_id === b.id ? 'selected' : ''}>${esc(b.bundle_id || b.name)}</option>`).join('')}</select></div>
        ${apiSources.length ? `<div><label>${esc(t('git_source'))}</label><select id="gx-source">${apiSources.map(x => `<option value="${x.id}" ${st.source_id === x.id ? 'selected' : ''}>${esc(x.name)}</option>`).join('')}</select></div>` : ''}
        <div><label>${esc(t('git_history_depth'))}</label><input id="gx-depth" type="number" min="1" max="50" value="${st.history_depth}"></div></div></div>
      <div class="g-card"><h3>${esc(t('git_servers'))}</h3><div class="g-mut" style="margin-bottom:6px">${esc(t('git_servers_h'))}</div>
        <div class="g-bar"><input type="search" placeholder="${esc(t('git_filter'))}" oninput="gitWorkspace.filterConns(this.value)">
        <button class="btn-sec btn-sm" onclick="gitWorkspace.selConns(true)">${esc(t('git_select_all'))}</button><button class="btn-sec btn-sm" onclick="gitWorkspace.selConns(false)">${esc(t('git_select_none'))}</button></div>
        <div class="g-conns" id="gx-conns">${sshConns.map(c => `<label data-q="${esc((c.name + ' ' + c.host + ' ' + (c.tags || []).join(' ')).toLowerCase())}"><input type="checkbox" value="${c.id}" ${sel.has(c.id) ? 'checked' : ''}>${esc(c.name)} <span class="g-mut">${esc(c.host)}</span></label>`).join('')}</div></div>
      <div class="g-card"><h3>${esc(t('git_monitor'))}</h3>${s.interval_minutes > 0 ? `<label style="display:flex;gap:6px;align-items:center"><input type="checkbox" id="gx-monitor" style="width:auto" ${st.monitor ? 'checked' : ''}>${esc(tf('git_monitor_on', {n: s.interval_minutes}))}</label>` : `<div class="g-mut">${esc(t('git_monitor_off'))}</div>`}
        <div class="g-mut" style="margin-top:6px">${esc(t('git_monitor_h'))}</div></div>
      <div class="g-bar"><button onclick="gitWorkspace.saveSettings()">${esc(t('git_save'))}</button></div>
      <div class="g-card"><h3>${esc(t('git_catalog_settings'))}</h3><div class="g-grid">
        <div><label>${esc(t('git_url'))}</label><input id="gc-url" value="${esc(cat.gitlab.url || '')}" placeholder="https://git.example.com"></div>
        <div><label>${esc(t('git_groups'))}</label><textarea id="gc-groups">${esc((cat.gitlab.groups || []).join('\n'))}</textarea></div>
        <div><label>${esc(t('git_fallback'))}</label><textarea id="gc-fallback">${esc((cat.fallback_branches || []).join('\n'))}</textarea></div>
        <div><label>${esc(t('git_roots'))}</label><textarea id="gc-roots">${esc((cat.server_roots || []).join('\n'))}</textarea></div>
        <div><label>${esc(t('git_ignore'))}</label><textarea id="gc-ignore">${esc((cat.ignore_dirs || []).join('\n'))}</textarea></div></div>
        <div class="g-bar" style="margin-top:8px"><button onclick="gitWorkspace.saveCatalogSettings()">${esc(t('git_save'))}</button></div></div>`;
  }

  function filterConns(q) { q = q.toLowerCase(); document.querySelectorAll('#gx-conns label').forEach(l => { l.style.display = l.dataset.q.includes(q) ? '' : 'none'; }); }
  function selConns(on) { document.querySelectorAll('#gx-conns label').forEach(l => { if (l.style.display !== 'none') l.querySelector('input').checked = on; }); }

  async function putSettings(patch) {
    const body = Object.assign({}, G.state.settings, patch);
    const r = await fetch('/api/git/settings', json('PUT', body));
    if (!r.ok) { showToast(await apiError(r), 'error'); return false; }
    G.state = await r.json(); render(); return true;
  }

  function saveSettings() {
    const mon = $('gx-monitor');
    putSettings({mode: $('gx-mode').value, bundle_id: Number($('gx-bundle').value) || 0, source_id: Number(($('gx-source') || {}).value) || 0,
      history_depth: Number($('gx-depth').value) || 0, conn_ids: [...document.querySelectorAll('#gx-conns input:checked')].map(x => Number(x.value)),
      monitor: mon ? mon.checked : G.state.settings.monitor}).then(ok => ok && showToast(t('git_saved'), 'success'));
  }

  async function saveCatalogSettings() {
    const cat = Object.assign({}, G.state.catalog, {gitlab: {url: $('gc-url').value.trim(), groups: lines($('gc-groups').value)},
      fallback_branches: lines($('gc-fallback').value), server_roots: lines($('gc-roots').value), ignore_dirs: lines($('gc-ignore').value)});
    const r = await fetch('/api/git/catalog', json('PUT', cat));
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    G.state = await r.json(); render(); showToast(t('git_saved'), 'success');
  }

  function setRef(mode) {
    const b = $('git-ref-branch');
    if (mode === 'branch' && !(b && b.value.trim())) { G.state.settings.ref_mode = 'branch'; render(); setTimeout(() => $('git-ref-branch') && $('git-ref-branch').focus(), 30); return; }
    putSettings({ref_mode: mode, ref_branch: b ? b.value.trim() : G.state.settings.ref_branch});
  }

  function editSource(id) {
    const x = G.state.sources.find(s => s.id === id) || {kind: 'gitlab', name: '', url: G.state.catalog.gitlab.url || ''};
    modal(`<h3>${esc(id ? x.name : t('git_add_source'))}</h3><div class="g-grid">
      <div><label>${esc(t('git_kind'))}</label><select id="gsrc-kind"><option value="gitlab" ${x.kind === 'gitlab' ? 'selected' : ''}>GitLab</option><option value="github" ${x.kind === 'github' ? 'selected' : ''}>GitHub</option></select></div>
      <div><label>${esc(t('git_f_name'))}</label><input id="gsrc-name" value="${esc(x.name)}"></div>
      <div><label>${esc(t('git_url'))}</label><input id="gsrc-url" value="${esc(x.url)}" placeholder="https://git.example.com"></div>
      <div><label>${esc(t('git_token'))}</label><input id="gsrc-token" type="password" autocomplete="new-password" placeholder="${id && x.has_token ? esc(t('git_token_keep')) : ''}"></div>
      <div><label>${esc(t('git_write_token'))}</label><input id="gsrc-wtoken" type="password" autocomplete="new-password" placeholder="${id && x.has_write_token ? esc(t('git_token_keep')) : ''}"></div></div>
      <div class="g-mut" style="margin-top:6px">${esc(t('git_write_token_h'))}</div>
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gsrc-save">${esc(t('git_save'))}</button></div>`);
    $('gsrc-save').onclick = async () => {
      const body = {kind: $('gsrc-kind').value, name: $('gsrc-name').value.trim(), url: $('gsrc-url').value.trim()};
      const tok = $('gsrc-token').value;
      if (tok || !id) body.token = tok;
      if ($('gsrc-wtoken').value) body.write_token = $('gsrc-wtoken').value;
      const r = await fetch(id ? '/api/git/sources/' + id : '/api/git/sources', json(id ? 'PUT' : 'POST', body));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = await r.json(); closeModal(); render();
    };
  }

  async function testSource(id) {
    const r = await fetch(`/api/git/sources/${id}/test`, {method: 'POST'});
    if (!r.ok) showToast(await apiError(r), 'error');
    else showToast(tf('git_projects_found', {n: ((await r.json()).projects || []).length}), 'success');
    load().then(ok => ok && render());
  }

  async function delSource(id) {
    if (!(await uiConfirm(t('git_delete') + '?', {danger: true, okText: t('git_delete')}))) return;
    const r = await fetch('/api/git/sources/' + id, {method: 'DELETE'});
    if (r.ok) { G.state = await r.json(); render(); }
  }

  async function importBundle(input) {
    const file = input.files[0]; input.value = '';
    if (!file) return;
    const r = await fetch('/api/git/bundles', {method: 'POST', headers: {'Content-Type': 'application/gzip'}, body: file});
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const j = await r.json();
    G.state = j.state; render();
    showToast(tf('git_bundle_imported', {n: j.result.files}) + (j.result.warnings.length ? ' · ⚠ ' + j.result.warnings[0] : ''), j.result.warnings.length ? 'warning' : 'success', 5000);
  }

  // ── runs ──
  // "Check now" checks everything (discovery on the servers of the settings, targets first)
  function check(discoverOnly) {
    if (G.busy || jobRunning()) return;
    startJob({all: true, discover: true, refresh: !discoverOnly});
  }

  async function refresh() {
    if (G.busy) return;
    G.busy = true; render();
    try {
      const r = await fetch('/api/git/targets/refresh', json('POST', {}));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = (await r.json()).state;
    } finally { G.busy = false; render(); }
  }

  // ── deploy: update, upgrade, rollback, stamp, restart ──
  const D = {plan: null, kind: '', ref: '', items: [], body: null};
  const can = k => !!(G.state.deploy || {})[k];
  const kindName = k => t('git_k_' + k);
  function runBadge(s) { return `<span class="g-st s-${s === 'ok' ? 'rs-ok' : esc(s)}">${esc(t('git_rs_' + s))}</span>`; }
  function localISO(v) { const d = new Date(v); return isNaN(d) ? '' : d.toISOString(); }
  function defaultTime(h) { const d = new Date(Date.now() + h * 3600e3); d.setMinutes(0, 0, 0); const p = n => String(n).padStart(2, '0'); return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}T${p(d.getHours())}:00`; }
  function installById(id) { return (G.state.installs || []).find(x => x.id === id) || {}; }

  function actionButtons(id) {
    const a = id ? `[${id}]` : 'null';
    const b = (k, label, cls) => `<button class="${cls || 'btn-sec'} btn-sm" ${can(k === 'stamp' ? 'update' : k) ? '' : `disabled title="${esc(t('git_no_permission'))}"`} onclick="gitWorkspace.${k}(${a})">${esc(t(label))}</button>`;
    return b('update', 'git_update', 'btn-sm') + b('upgrade', 'git_upgrade') + (id ? b('rollback', 'git_rollback') : '') + b('restart', 'git_restart') + b('stamp', 'git_stamp') + b('transfer', 'git_transfer');
  }

  function pickIds(ids) {
    ids = ids || [...G.sel];
    if (!ids.length) showToast(t('git_select_installs'), 'warning');
    return ids;
  }

  async function plan(kind, ids, ref) {
    modal(`<h3>${esc(kindName(kind))}</h3><div class="hint-box">⏳ ${esc(t('git_plan_loading'))}</div>`);
    const r = await fetch('/api/git/plan', json('POST', {kind, ref: ref || '', install_ids: ids}));
    if (!r.ok) { modal(`<h3>${esc(kindName(kind))}</h3><div class="hint-box bad">${esc(await apiError(r))}</div><div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`); return null; }
    return (await r.json()).items;
  }

  function update(ids) { ids = pickIds(ids); if (ids.length) openDeploy('update', ids, ''); }

  function upgrade(ids) {
    ids = pickIds(ids);
    if (!ids.length) return;
    modal(`<h3>${esc(kindName('upgrade'))}</h3><div class="hint-box warn">${esc(t('git_upgrade_warn'))}</div>
      <div class="g-grid"><div><label>${esc(t('git_upgrade_ref'))}</label><input id="gu-ref" value="${esc(D.ref || '')}" placeholder="${esc(t('git_upgrade_ref_h'))}"></div></div>
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gu-go">${esc(t('git_show_changes'))}</button></div>`);
    const go = () => { const ref = $('gu-ref').value.trim(); if (ref) { D.ref = ref; openDeploy('upgrade', ids, ref); } };
    $('gu-go').onclick = go;
    $('gu-ref').onkeydown = e => { if (e.key === 'Enter') go(); };
  }

  async function openDeploy(kind, ids, ref) {
    const items = await plan(kind, ids, ref);
    if (!items) return;
    D.plan = items; D.kind = kind; D.ref = ref;
    const st = G.state.settings || {};
    const apps = [...new Set(items.map(i => i.app))];
    const units = items.filter(i => (i.units || []).length);
    const fileRow = (i, idx, f) => {
      const dis = f.state === 'ok' || f.state === 'protected';
      return `<tr class="${f.state === 'ok' ? 'g-okrow' : ''}"><td><input type="checkbox" class="gd-f" data-i="${idx}" data-p="${esc(f.path)}" data-local="${f.local ? 1 : 0}" ${f.selected ? 'checked' : ''} ${dis ? 'disabled' : ''} onchange="gitWorkspace.localWarn()"></td>
        <td class="g-mono">${esc(f.path)}</td><td>${fsBadge(f.state)}${f.local ? ` <span class="g-mut">${esc(t('git_local'))}</span>` : ''}${f.behind ? ` <span class="g-mut">${esc(tf('git_behind', {n: f.behind}))} (${esc(f.tag || '')} ${esc(f.date || '')})</span>` : ''}</td></tr>`;
    };
    const itemHTML = (i, idx) => {
      const changes = (i.files || []).filter(f => f.state !== 'ok' && f.state !== 'protected');
      const head = `<h4>${esc(i.conn_name)} · ${esc(i.app)} · <span class="g-mono">${esc(i.path)}</span>${i.env ? `<span class="g-env">${esc(i.env)}</span>` : ''}${i.prod ? `<span class="g-prod">${esc(t('git_prod'))}</span>` : ''}</h4>`;
      if (i.error) return `<div class="g-item">${head}<div class="hint-box bad">${esc(i.error)}</div></div>`;
      const ver = `<div class="g-mut">${esc(t('git_version_md'))}: <span class="g-mono">${esc(i.from_version || t('git_vmd_missing'))}</span> → ${esc(t('git_target'))}: <span class="g-mono">${esc(i.target ? i.target.version : '—')}</span> ${i.target ? `<span class="g-mut">${esc(i.target.branch)} · ${esc(i.target.commit_date || '')}</span>` : ''}${kind === 'upgrade' && i.current ? ` <span class="g-mut">(${esc(t('git_from'))} ${esc(i.current.version)} · ${esc(i.current.branch)})</span>` : ''}</div>`;
      if (i.pipeline) return `<div class="g-item">${head}${ver}<div class="hint-box">🔧 ${esc(tf('git_ci_item', {kind: i.pipeline.kind === 'gitlab' ? 'GitLab' : 'Jenkins'}))}<div class="g-mono g-mut">${esc(i.pipeline.url)}</div></div></div>`;
      if (kind === 'stamp') return `<div class="g-item">${head}${ver}<div>${i.eligible ? fsBadge('ok') : `<span class="g-st s-review">${esc(t('git_not_current'))}</span>`} ${i.has_version_md ? `<span class="g-mut">${esc(t('git_has_vmd'))}</span>` : ''}</div></div>`;
      const c = i.counts || {};
      const sum = kind === 'upgrade' ? `<div class="hint-box warn">${esc(tf('git_changes_summary', {changed: (c.changed || 0) + (c.old || 0), missing: c.missing || 0, local: c.modified || 0, removed: (i.removed || []).length}))}${(i.removed || []).length ? `<div class="g-mono g-mut">${(i.removed || []).map(esc).join(', ')}</div>` : ''}</div>` : '';
      return `<div class="g-item">${head}${ver}${sum}${changes.length ? '' : `<div class="g-mut">${esc(t('git_files_none'))}</div>`}
        <div class="g-files"><table>${(i.files || []).map(f => fileRow(i, idx, f)).join('')}</table></div></div>`;
    };
    const restartHTML = () => {
      if (kind === 'stamp') return '';
      if (!units.length) return `<div class="g-card"><h3>${esc(t('git_restart_choice'))}</h3><div class="g-mut">${esc(t('git_no_units_restart'))}</div></div>`;
      const dis = can('restart') ? '' : 'disabled';
      return `<div class="g-card"><h3>${esc(t('git_restart_choice'))}</h3>
        <label class="g-opt"><input type="radio" name="gd-rm" value="none" checked onchange="gitWorkspace.restartWarn()"><span>${esc(t('git_rm_none'))}</span></label>
        <label class="g-opt"><input type="radio" name="gd-rm" value="now" ${dis} onchange="gitWorkspace.restartWarn()"><span>${esc(t('git_rm_now'))}</span></label>
        <label class="g-opt"><input type="radio" name="gd-rm" value="at" ${dis} onchange="gitWorkspace.restartWarn()"><span>${esc(t('git_rm_at'))}</span> <input type="datetime-local" id="gd-rat" value="${defaultTime(2)}"></label>
        <label class="g-opt"><input type="radio" name="gd-rm" value="delay" ${dis} onchange="gitWorkspace.restartWarn()"><span>${esc(t('git_rm_delay'))}</span> <input type="number" id="gd-rdelay" min="1" max="10080" value="15" style="width:90px"></label>
        <div id="gd-rwarn" class="hint-box warn" style="display:none">${esc(tf('git_restart_warn', {list: units.map(i => i.conn_name + ': ' + i.units.map(u => u.name).join(', ')).join('; ')}))}</div>
        ${!can('restart') ? `<div class="g-mut">${esc(t('git_no_permission'))}</div>` : ''}</div>`;
    };
    modal(`<h3>${esc(kindName(kind))}${kind === 'upgrade' ? ' → <span class="g-mono">' + esc(ref) + '</span>' : ''}</h3>
      ${kind === 'upgrade' ? `<div class="hint-box warn">${esc(t('git_upgrade_warn'))} ${esc(t('git_after_upgrade'))}</div>` : ''}
      ${kind === 'stamp' ? `<div class="g-mut">${esc(t('git_stamp_h'))}</div>` : `<label class="g-opt"><input type="checkbox" onchange="document.getElementById('gd-items').classList.toggle('g-hideok', !this.checked)"><span>${esc(t('git_show_all'))}</span></label>`}
      <div id="gd-items" class="g-hideok">${items.map(itemHTML).join('')}</div>
      <div id="gd-lwarn" class="hint-box bad" style="display:none">${esc(t('git_local_warn'))}</div>
      ${kind === 'stamp' ? `<label class="g-opt"><input type="checkbox" id="gd-force"><span>${esc(t('git_force_stamp'))}</span></label>` : `<div class="g-card"><div class="g-grid">${apps.map(a => `<div><label>${esc(tf('git_check_cmd', {app: a}))}</label><input class="gd-check" data-app="${esc(a)}" value="${esc((st.check_commands || {})[a] || '')}" placeholder="./run.py --selftest"></div>`).join('')}</div>
        <label class="g-opt"><input type="checkbox" id="gd-noimports"><span>${esc(t('git_ignore_imports'))}</span></label>
        ${kind === 'upgrade' ? `<label class="g-opt"><input type="checkbox" id="gd-override"><span>${esc(t('git_set_override'))}</span></label>` : ''}</div>`}
      ${restartHTML()}
      ${kind === 'stamp' ? '' : `<div class="g-card"><h3>${esc(t('git_when'))}</h3>
        <label class="g-opt"><input type="radio" name="gd-when" value="now" checked><span>${esc(t('git_run_now'))}</span></label>
        <label class="g-opt"><input type="radio" name="gd-when" value="at"><span>${esc(t('git_run_at'))}</span> <input type="datetime-local" id="gd-at" value="${defaultTime(1)}"></label>
        <div class="g-mut">${esc(tf('git_grace_h', {n: G.state.grace_minutes || 30}))}</div></div>`}
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gd-next">${esc(t('git_continue'))}</button></div>`);
    $('gd-next').onclick = deployNext;
  }

  function localWarn() { $('gd-lwarn').style.display = [...document.querySelectorAll('.gd-f:checked')].some(x => x.dataset.local === '1') ? '' : 'none'; }
  function restartWarn() { const r = document.querySelector('input[name=gd-rm]:checked'); $('gd-rwarn').style.display = r && r.value !== 'none' ? '' : 'none'; }

  function deployNext() {
    const kind = D.kind;
    const items = [];
    D.plan.forEach((i, idx) => {
      if (i.error) return;
      if (kind === 'stamp') { if (i.eligible) items.push({install_id: i.install_id}); return; }
      if (i.pipeline) { items.push({install_id: i.install_id}); return; }
      const boxes = [...document.querySelectorAll(`.gd-f[data-i="${idx}"]:checked`)];
      if (!boxes.length) return;
      items.push({install_id: i.install_id, files: boxes.map(b => b.dataset.p), modified_ok: boxes.filter(b => b.dataset.local === '1').map(b => b.dataset.p)});
    });
    if (!items.length) { showToast(t('git_nothing_selected'), 'warning'); return; }
    const body = {kind, ref: D.ref || '', items};
    if (kind === 'stamp') body.force = $('gd-force').checked;
    else {
      body.checks = Object.fromEntries([...document.querySelectorAll('.gd-check')].map(x => [x.dataset.app, x.value.trim()]));
      body.ignore_imports = $('gd-noimports').checked;
      if ($('gd-override')) body.set_ref_override = $('gd-override').checked;
      const rm = document.querySelector('input[name=gd-rm]:checked');
      body.restart = {mode: rm ? rm.value : 'none'};
      if (body.restart.mode === 'at') body.restart.at = localISO($('gd-rat').value);
      if (body.restart.mode === 'delay') body.restart.delay_minutes = Number($('gd-rdelay').value) || 0;
      if (document.querySelector('input[name=gd-when]:checked').value === 'at') {
        body.schedule_at = localISO($('gd-at').value);
        if (!body.schedule_at || new Date(body.schedule_at) < new Date()) { showToast(t('git_time_past'), 'warning'); return; }
      }
      if (body.restart.mode === 'at' && (!body.restart.at || new Date(body.restart.at) < new Date())) { showToast(t('git_time_past'), 'warning'); return; }
    }
    confirmRun(body, D.plan.filter(i => items.some(x => x.install_id === i.install_id)));
  }

  // second confirmation: a summary, and the server name typed for production installations
  function confirmRun(body, list) {
    D.body = body;
    const servers = [...new Set(list.map(i => i.conn_name))];
    const prod = [...new Map(list.filter(i => i.prod).map(i => [i.conn_name, i])).values()];
    const rm = (body.restart || {}).mode;
    const units = list.filter(i => (i.units || []).length);
    modal(`<h3>${esc(tf('git_confirm_title', {kind: kindName(body.kind)}))}</h3>
      <div>${esc(tf('git_confirm_q', {kind: kindName(body.kind), n: list.length, m: servers.length}))}</div>
      <ul>${list.map(i => `<li><b>${esc(i.app)}</b> @ ${esc(i.conn_name)} <span class="g-mono">${esc(i.path)}</span>${i.prod ? `<span class="g-prod">${esc(t('git_prod'))}</span>` : ''}${body.backup_names ? ` ← <span class="g-mono">${esc(body.backup_names[i.install_id] || '')}</span>` : ''}</li>`).join('')}</ul>
      ${body.kind === 'upgrade' ? `<div class="hint-box warn">${esc(t('git_upgrade_warn'))} → <b class="g-mono">${esc(body.ref)}</b></div>` : ''}
      ${rm && rm !== 'none' && units.length ? `<div class="hint-box warn">${esc(tf('git_restart_warn', {list: units.map(i => i.conn_name + ': ' + i.units.map(u => u.name).join(', ')).join('; ')}))}${body.restart.at ? ' · ' + esc(fmtTime(body.restart.at)) : ''}${body.restart.delay_minutes ? ' · +' + body.restart.delay_minutes + ' min' : ''}</div>` : ''}
      ${body.schedule_at ? `<div class="hint-box">⏰ ${esc(tf('git_scheduled_ok', {t: fmtTime(body.schedule_at)}))}</div>` : ''}
      ${prod.map(i => `<div class="g-conf"><label>${esc(tf('git_type_host', {name: i.conn_name}))}</label><input class="gd-conf" data-name="${esc(i.conn_name)}" autocomplete="off" oninput="gitWorkspace.confirmCheck()"></div>`).join('')}
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gd-go" class="${body.kind === 'stamp' ? '' : 'btn-danger'}" ${prod.length ? 'disabled' : ''}>${esc(t(body.schedule_at || (body.kind === 'restart' && rm !== 'now') ? 'git_schedule' : 'git_start'))}</button></div>`);
    D.items = list;
    $('gd-go').onclick = startRun;
  }

  function confirmCheck() { $('gd-go').disabled = ![...document.querySelectorAll('.gd-conf')].every(x => x.value.trim().toLowerCase() === x.dataset.name.trim().toLowerCase()); }

  async function startRun() {
    const body = Object.assign({}, D.body);
    delete body.backup_names;
    const typed = Object.fromEntries([...document.querySelectorAll('.gd-conf')].map(x => [x.dataset.name, x.value.trim()]));
    body.confirm = Object.fromEntries(D.items.filter(i => i.prod).map(i => [String(i.install_id), typed[i.conn_name] || '']));
    $('gd-go').disabled = true;
    const r = await fetch('/api/git/runs', json('POST', body));
    if (!r.ok) { $('gd-go').disabled = false; showToast(await apiError(r), 'error', 7000); return; }
    const v = await r.json();
    if (v.state === 'scheduled') {
      showToast(tf('git_scheduled_ok', {t: fmtTime(v.scheduled_at)}), 'success', 5000);
      closeModal(); G.tab = 'runs'; render(); loadRuns();
      return;
    }
    showRun(v.id);
  }

  function rollback(ids) {
    ids = pickIds(ids);
    if (!ids.length) return;
    const i = installById(ids[0]);
    modal(`<h3>${esc(kindName('rollback'))}: ${esc(i.conn_name || '')} · <span class="g-mono">${esc(i.path || '')}</span></h3><div class="hint-box">⏳ ${esc(t('git_plan_loading'))}</div>`);
    fetch(`/api/git/installs/${ids[0]}/backups`).then(async r => {
      if (!r.ok) { modal(`<h3>${esc(kindName('rollback'))}</h3><div class="hint-box bad">${esc(await apiError(r))}</div><div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`); return; }
      const j = await r.json();
      const list = j.backups || [];
      modal(`<h3>${esc(kindName('rollback'))}: ${esc(i.conn_name || '')} · <span class="g-mono">${esc(i.path || '')}</span></h3>
        <div class="g-mut">${esc(t('git_restore_hint'))}</div><h3 style="margin-top:10px">${esc(t('git_backups'))} <span class="g-mut g-mono">${esc(j.dir || '')}</span></h3>
        ${list.length ? `<div class="g-files"><table>${list.map((b, k) => `<tr><td><input type="radio" name="gr-b" value="${esc(b.name)}" ${k === 0 ? 'checked' : ''}></td>
          <td><span class="g-mono">${esc(b.name)}</span><div class="g-mut">${esc(b.time || '')}${b.user ? ' · ' + esc(b.user) : ''}${b.version ? ' · VERSION.md ' + esc(b.version) : ''}${b.from_version ? ' · ' + esc(b.from_version) + ' → ' + esc(b.to_version || '') : ''}</div>
          <details><summary class="g-mut">${esc(tf('git_backup_files', {n: (b.log_files || b.files || []).length}))}</summary><div class="g-mono g-mut">${(b.log_files || b.files || []).map(esc).join('<br>')}</div></details>
          ${b.log_files ? '' : `<div class="g-mut">⚠ ${esc(t('git_backup_unknown'))}</div>`}</td></tr>`).join('')}</table></div>` : `<div class="hint-box">${esc(t('git_no_backups'))}</div>`}
        <div class="g-card" style="margin-top:10px"><label>${esc(tf('git_check_cmd', {app: i.app}))}</label><input id="gr-check" value="${esc(((G.state.settings || {}).check_commands || {})[i.app] || '')}"></div>
        <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button>${list.length ? `<button id="gr-next">${esc(t('git_continue'))}</button>` : ''}</div>`);
      if (!list.length) return;
      $('gr-next').onclick = () => {
        const b = document.querySelector('input[name=gr-b]:checked').value;
        confirmRun({kind: 'rollback', items: [{install_id: i.id, backup: b}], checks: {[i.app]: $('gr-check').value.trim()}, backup_names: {[i.id]: b}},
          [{install_id: i.id, app: i.app, conn_name: i.conn_name, path: i.path, prod: isProd(i), units: i.units}]);
      };
    });
  }

  function isProd(i) {
    const p = ['prod', 'production', 'prd', 'live'];
    return p.includes(String(i.env || '').toLowerCase()) || (i.tags || []).some(x => p.includes(x.toLowerCase()) || (x.toLowerCase().startsWith('env:') && p.includes(x.toLowerCase().slice(4))));
  }

  function restart(ids) {
    ids = pickIds(ids);
    if (!ids.length) return;
    const list = ids.map(installById).filter(i => i.id);
    const withUnits = list.filter(i => (i.units || []).length);
    modal(`<h3>${esc(kindName('restart'))}</h3>
      ${withUnits.length ? `<table><tr><th>${esc(t('git_server'))}</th><th>${esc(t('git_path'))}</th><th>${esc(t('git_unit_list'))}</th></tr>${withUnits.map(i => `<tr><td>${esc(i.conn_name)}${isProd(i) ? `<span class="g-prod">${esc(t('git_prod'))}</span>` : ''}</td><td class="g-mono">${esc(i.path)}</td><td class="g-mono">${i.units.map(u => esc(u.kind + ': ' + u.name)).join('<br>')}</td></tr>`).join('')}</table>`
        : `<div class="hint-box">${esc(t('git_no_units_restart'))}</div>`}
      ${withUnits.length ? `<div class="g-card" style="margin-top:10px"><label class="g-opt"><input type="radio" name="gx-rm" value="now" checked><span>${esc(t('git_run_now'))}</span></label>
        <label class="g-opt"><input type="radio" name="gx-rm" value="at"><span>${esc(t('git_rm_at'))}</span> <input type="datetime-local" id="gx-rat" value="${defaultTime(2)}"></label>
        <label class="g-opt"><input type="radio" name="gx-rm" value="delay"><span>${esc(t('git_rm_delay'))}</span> <input type="number" id="gx-rdelay" min="1" max="10080" value="15" style="width:90px"></label>
        <div class="g-mut">${esc(tf('git_grace_h', {n: G.state.grace_minutes || 30}))}</div></div>` : ''}
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button>${withUnits.length ? `<button id="gx-next">${esc(t('git_continue'))}</button>` : ''}</div>`);
    if (!withUnits.length) return;
    $('gx-next').onclick = () => {
      const mode = document.querySelector('input[name=gx-rm]:checked').value;
      const body = {kind: 'restart', items: withUnits.map(i => ({install_id: i.id})), restart: {mode}};
      if (mode === 'at') { body.restart.at = localISO($('gx-rat').value); if (!body.restart.at || new Date(body.restart.at) < new Date()) { showToast(t('git_time_past'), 'warning'); return; } }
      if (mode === 'delay') body.restart.delay_minutes = Number($('gx-rdelay').value) || 0;
      confirmRun(body, withUnits.map(i => ({install_id: i.id, app: i.app, conn_name: i.conn_name, path: i.path, prod: isProd(i), units: i.units})));
    };
  }

  function stamp(ids) { ids = pickIds(ids); if (ids.length) openDeploy('stamp', ids, ''); }

  // ── new server: install and transfer wizards ──
  const W = {};
  const kb = n => Math.ceil((n || 0) / 1024);
  const sshConns = () => (conns || []).filter(c => c.protocol === 'SSH').sort((a, b) => a.name.localeCompare(b.name));
  const errBox = async r => `<div class="hint-box bad">${esc(await apiError(r))}</div>`;
  const btns = (next, label) => `<div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="${next}">${esc(t(label))}</button></div>`;

  function install() {
    if (!can('install')) return;
    const s = G.state;
    const apps = [...s.catalog.apps.map(a => a.name), ...(s.bundle_apps || [])];
    const targets = Object.fromEntries(s.targets.map(x => [x.app, x]));
    const bundles = s.sources.filter(x => x.kind === 'bundle');
    W.kind = 'install';
    const refSel = a => `<select class="gw-ref" data-app="${esc(a)}" onchange="gitWorkspace.wSvcRef(this)" style="width:auto">
        <option value="">${esc(tf('git_w_ref_target', {v: targets[a] ? targets[a].version : t('git_no_target')}))}</option>
        ${s.sources.some(x => x.kind !== 'bundle') ? `<option value="?">${esc(t('git_w_ref_other'))}</option>` : ''}
        ${bundles.filter(b => (b.apps || []).includes(a)).map(b => `<option value="bundle:${b.id}">${esc(tf('git_w_ref_bundle', {b: b.bundle_id || b.name}))}</option>`).join('')}</select>
        <input class="gw-refx" data-app="${esc(a)}" placeholder="tag:v2 / main" style="display:none;width:130px">`;
    modal(`<h3>${esc(t('git_k_install'))}: ${esc(t('git_new_server'))}</h3>
      <div class="g-grid"><div><label>${esc(t('git_w_server'))}</label><select id="gw-conn">${sshConns().map(c => `<option value="${c.id}">${esc(c.name)} (${esc(c.host)})</option>`).join('')}</select></div></div>
      <div class="g-card"><h3>${esc(t('git_w_services'))}</h3><div class="g-files"><table>${apps.map(a => `<tr><td><input type="checkbox" class="gw-app" value="${esc(a)}"></td><td><b>${esc(a)}</b></td><td>${refSel(a)}</td></tr>`).join('')}</table></div></div>
      ${btns('gw-next', 'git_continue')}`);
    $('gw-next').onclick = () => {
      const svcs = [...document.querySelectorAll('.gw-app:checked')].map(x => {
        const a = x.value, sel = document.querySelector(`.gw-ref[data-app="${CSS.escape(a)}"]`);
        const ref = sel.value === '?' ? document.querySelector(`.gw-refx[data-app="${CSS.escape(a)}"]`).value.trim() : sel.value;
        return {app: a, ref};
      });
      if (!svcs.length || !$('gw-conn').value) { showToast(t('git_w_pick_service'), 'warning'); return; }
      W.conn = Number($('gw-conn').value); W.svcs = svcs;
      installDirs();
    };
  }
  function wSvcRef(sel) { const x = document.querySelector(`.gw-refx[data-app="${CSS.escape(sel.dataset.app)}"]`); x.style.display = sel.value === '?' ? '' : 'none'; if (sel.value === '?') x.focus(); }

  async function prepare(body) {
    const r = await fetch('/api/git/provision/prepare', json('POST', body));
    if (!r.ok) { modal(`<h3>${esc(t('git_k_' + W.kind))}</h3>${await errBox(r)}<div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`); return null; }
    return r.json();
  }

  function rootsHTML(p) {
    return `<div class="g-mut">${esc(t('git_w_roots'))}: ${(p.roots || []).map(r => `<span class="g-mono">${esc(r.path)}</span> ${r.exists ? (r.writable ? '✅' : '⚠ ' + esc(t('git_w_root_ro'))) : '✗ ' + esc(t('git_w_root_missing'))}`).join(' · ') || '—'}</div>`;
  }

  async function installDirs() {
    modal(`<h3>${esc(t('git_k_install'))}</h3><div class="hint-box">⏳ ${esc(t('git_plan_loading'))}</div>`);
    const p = await prepare({kind: 'install', conn_id: W.conn, services: W.svcs});
    if (!p) return;
    const root = ((p.roots || []).find(r => r.exists && r.writable) || (p.roots || [])[0] || {}).path || '/opt';
    modal(`<h3>${esc(t('git_k_install'))} → ${esc(p.conn_name)}${p.prod ? `<span class="g-prod">${esc(t('git_prod'))}</span>` : ''}</h3>${rootsHTML(p)}
      <div class="g-card"><h3>${esc(t('git_w_dirs'))}</h3>${p.services.map((x, i) => `<div class="g-item"><h4>${esc(x.app)} ${x.target ? `<span class="g-mono">${esc(x.target.version)}</span> <span class="g-mut">${esc(x.target.branch)} · ${esc(x.target.source)} · ${esc(tf('git_w_files', {n: x.files, kb: kb(x.bytes)}))}</span>` : ''}</h4>
        ${x.error ? `<div class="hint-box bad">${esc(x.error)}</div>` : `<div class="g-grid"><div><label>${esc(t('git_path'))}</label><input class="gw-path" data-i="${i}" value="${esc(root.replace(/\/$/, '') + '/' + x.app)}"></div>
        <div><label>${esc(t('git_w_env'))}</label><input class="gw-env" data-i="${i}" placeholder="test / prod"></div></div>`}</div>`).join('')}</div>
      ${btns('gw-next', 'git_w_check')}`);
    if (p.services.some(x => x.error)) { $('gw-next').disabled = true; return; }
    $('gw-next').onclick = () => {
      W.svcs = W.svcs.map((x, i) => Object.assign({}, x, {path: document.querySelector(`.gw-path[data-i="${i}"]`).value.trim(), env: document.querySelector(`.gw-env[data-i="${i}"]`).value.trim()}));
      installCheck();
    };
  }

  function checkHTML(c, idx) {
    if (!c) return '';
    const tools = Object.entries(c.tools || {}).filter(([k]) => k !== 'tar').map(([k, v]) => `${esc(k)} <span class="g-mono">${esc(v)}</span>`);
    const errs = c.errors || [];
    return `<div class="g-mut">${esc(t('git_w_checks'))}: ${esc(c.parent)} · ${esc(tf('git_w_free', {free: Math.floor(c.free_kb / 1024), need: Math.ceil(c.need_kb / 1024)}))}${tools.length ? ' · ' + tools.join(', ') : ''}${(c.needs || []).filter(n => !(c.tools || {})[n]).map(n => ` · ⚠ ${esc(n)} ${esc(t('git_w_missing_tool'))}`).join('')}</div>
      ${errs.map(e => `<div class="hint-box bad">${esc(e)}</div>`).join('')}${(c.warnings || []).map(w => `<div class="hint-box warn">${esc(w)}</div>`).join('')}
      ${c.existing ? `<div class="hint-box warn">${esc(tf('git_w_existing', {n: c.files, v: c.version ? ', VERSION.md ' + c.version : ''}))}<label class="g-opt"><input type="checkbox" class="gw-over" data-i="${idx}"><span>${esc(t('git_w_overwrite'))}</span></label></div>` : ''}`;
  }
  const blocking = c => c && (c.errors || []).length > 0; // prepared with overwrite: an existing installation is a tick, not an error

  function slotHTML(i, k, sl) {
    const tmpls = sl.templates.filter(x => x.available);
    const mode = tmpls.length ? 'template' : sl.copies.length ? 'copy' : 'skip';
    const radio = (m, label, dis) => `<label class="g-opt"><input type="radio" name="gw-m-${i}-${k}" value="${m}" ${m === mode ? 'checked' : ''} ${dis ? 'disabled' : ''} onchange="gitWorkspace.wSlotMode(${i}, ${k})"><span>${esc(t(label))}</span></label>`;
    const vaultOpts = (W.creds || []).map(c => `<option value="${c.id}:password">${esc(tf('git_w_vault_pw', {n: c.name}))}</option><option value="${c.id}:username">${esc(tf('git_w_vault_user', {n: c.name}))}</option>`).join('');
    const fields = (tp, ti) => `<div class="g-grid gw-fields" data-t="${ti}" ${ti ? 'style="display:none"' : ''}>${tp.fields.map(f => `<div><label class="g-mono">${esc(f.name)}</label>
      <input class="gw-fv" data-id="${esc(f.id)}" data-def="${esc(f.default)}" value="${esc(f.default)}" type="${f.secret ? 'password' : 'text'}" autocomplete="off">
      ${vaultOpts ? `<select class="gw-fs" data-id="${esc(f.id)}" onchange="this.previousElementSibling.style.display = this.value ? 'none' : ''"><option value="">${esc(t('git_w_typed'))}</option>${vaultOpts}</select>` : ''}</div>`).join('')}</div>`;
    return `<div class="g-item gw-slot" data-i="${i}" data-k="${k}" data-path="${esc(sl.path)}"><h4><span class="g-mono">${esc(sl.path)}</span> ${sl.protected ? fsBadge('protected') : ''}</h4>
      <div class="g-bar">${radio('template', 'git_w_m_template', !tmpls.length)}${radio('copy', 'git_w_m_copy', !sl.copies.length)}${radio('manual', 'git_w_m_manual')}${radio('skip', 'git_w_m_skip')}</div>
      <div class="gw-p gw-p-template" ${mode === 'template' ? '' : 'style="display:none"'}>${tmpls.length > 1 ? `<select class="gw-tsel" onchange="this.parentElement.querySelectorAll('.gw-fields').forEach(x => x.style.display = x.dataset.t === this.value ? '' : 'none')">${tmpls.map((x, ti) => `<option value="${ti}">${esc(x.path)}</option>`).join('')}</select>` : ''}
        ${tmpls.map((tp, ti) => `<div class="g-mut gw-tname" data-t="${ti}" data-path="${esc(tp.path)}" ${ti ? 'style="display:none"' : ''}>${esc(tp.path)}${sl.in_repo && tp.path === sl.path ? ' (' + esc(t('git_w_in_repo')) + ')' : ''}</div>`).join('')}
        ${tmpls.map(fields).join('')}${sl.templates.filter(x => !x.available).map(x => `<div class="g-mut">${esc(x.path)}: ${esc(t('git_w_tmpl_na'))}</div>`).join('')}</div>
      <div class="gw-p gw-p-copy" ${mode === 'copy' ? '' : 'style="display:none"'}><select class="gw-copy">${sl.copies.map(c => `<option value="${c.install_id}">${esc(c.conn_name)} · ${esc(c.path)}</option>`).join('')}</select></div>
      <div class="gw-p gw-p-manual" style="display:none"><textarea class="gw-manual" rows="5" spellcheck="false"></textarea></div></div>`;
  }
  function wSlotMode(i, k) {
    const el = document.querySelector(`.gw-slot[data-i="${i}"][data-k="${k}"]`);
    const m = el.querySelector('input[type=radio]:checked').value;
    el.querySelectorAll('.gw-p').forEach(x => { x.style.display = x.classList.contains('gw-p-' + m) ? '' : 'none'; });
  }

  function unitHTML(i, app, dir) {
    return `<div class="g-card"><label class="g-opt"><input type="checkbox" class="gw-unit-on" data-i="${i}" onchange="gitWorkspace.wUnit(${i})"><span><b>${esc(t('git_w_unit'))}</b></span></label>
      <div class="gw-unit" data-i="${i}" style="display:none"><div class="hint-box warn">${esc(t('git_w_unit_warn'))}</div><div class="g-grid">
        <div><label>${esc(t('git_w_unit_kind'))}</label><select class="gw-uk"><option value="systemd">systemd</option><option value="supervisor">supervisor</option></select></div>
        <div><label>${esc(t('git_w_unit_name'))}</label><input class="gw-un" value="${esc(app)}"></div>
        <div><label>${esc(t('git_w_unit_user'))}</label><input class="gw-uu" placeholder="svc-${esc(app)}"></div>
        <div><label>${esc(t('git_w_unit_cmd'))}</label><input class="gw-uc" placeholder="/usr/bin/python3 ${esc(dir)}/run.py"></div>
        <div><label>${esc(t('git_w_unit_dir'))}</label><input class="gw-ud" value="${esc(dir)}"></div></div>
        <label class="g-opt"><input type="checkbox" class="gw-ue" checked><span>${esc(t('git_w_unit_enable'))}</span></label></div></div>`;
  }
  function wUnit(i) { document.querySelector(`.gw-unit[data-i="${i}"]`).style.display = document.querySelector(`.gw-unit-on[data-i="${i}"]`).checked ? '' : 'none'; }
  function unitOf(i) {
    if (!document.querySelector(`.gw-unit-on[data-i="${i}"]`).checked) return null;
    const el = document.querySelector(`.gw-unit[data-i="${i}"]`), v = c => el.querySelector(c).value.trim();
    return {kind: v('.gw-uk'), name: v('.gw-un'), user: v('.gw-uu'), command: v('.gw-uc'), work_dir: v('.gw-ud'), enable: el.querySelector('.gw-ue').checked};
  }
  const startHTML = () => `<div class="g-card"><label class="g-opt"><input type="checkbox" id="gw-start" ${can('restart') ? '' : 'disabled'}><span>${esc(t('git_w_start'))}</span></label>
    <div class="g-mut">${esc(can('restart') ? t('git_w_start_h') : t('git_no_permission'))}</div></div>`;

  async function installCheck() {
    modal(`<h3>${esc(t('git_k_install'))}</h3><div class="hint-box">⏳ ${esc(t('git_plan_loading'))}</div>`);
    if (!W.creds) { const r = await fetch('/api/credentials'); W.creds = r.ok ? (await r.json()).filter(c => c.has_password || c.username) : []; }
    const p = await prepare({kind: 'install', conn_id: W.conn, services: W.svcs, overwrite: true});
    if (!p) return;
    W.prep = p;
    const st = G.state.settings || {};
    modal(`<h3>${esc(t('git_k_install'))} → ${esc(p.conn_name)}${p.prod ? `<span class="g-prod">${esc(t('git_prod'))}</span>` : ''}</h3>
      ${p.services.map((x, i) => `<div class="g-card"><h3>${esc(x.app)} <span class="g-mono">${esc(x.target ? x.target.version : '')}</span> → <span class="g-mono">${esc(x.path || '')}</span></h3>
        ${x.error ? `<div class="hint-box bad">${esc(x.error)}</div>` : checkHTML(x.check, i)}
        <h3 style="margin-top:10px">${esc(t('git_w_protected'))}</h3><div class="g-mut">${esc(t('git_w_protected_h'))}${W.creds.length ? ' ' + esc(t('git_w_vault_h')) : ''}</div>
        ${x.slots.length ? x.slots.map((sl, k) => slotHTML(i, k, sl)).join('') : `<div class="g-mut">${esc(t('git_w_no_slots'))}</div>`}
        <div style="margin-top:8px"><label class="g-mut">${esc(tf('git_check_cmd', {app: x.app}))}</label><input class="gw-chk" data-app="${esc(x.app)}" value="${esc((st.check_commands || {})[x.app] || '')}" style="width:100%"></div>
        ${unitHTML(i, x.app, x.path || '')}</div>`).join('')}
      ${startHTML()}<div id="gw-block" class="hint-box bad" style="display:none">${esc(t('git_w_blocked'))}</div>${btns('gw-next', 'git_continue')}`);
    $('gw-next').onclick = installNext;
  }

  function slotSpec(el) {
    const m = el.querySelector('input[type=radio]:checked').value, f = {path: el.dataset.path, mode: m};
    if (m === 'copy') f.from_install = Number(el.querySelector('.gw-copy').value);
    if (m === 'manual') f.content = el.querySelector('.gw-manual').value;
    if (m === 'template') {
      const ti = el.querySelector('.gw-tsel') ? el.querySelector('.gw-tsel').value : '0';
      f.template = el.querySelector(`.gw-tname[data-t="${ti}"]`).dataset.path;
      f.values = {};
      el.querySelectorAll(`.gw-fields[data-t="${ti}"] .gw-fv`).forEach(inp => {
        const vs = inp.parentElement.querySelector('.gw-fs');
        if (vs && vs.value) { const [id, part] = vs.value.split(':'); f.values[inp.dataset.id] = {cred_id: Number(id), part}; }
        else if (inp.value !== inp.dataset.def) f.values[inp.dataset.id] = {value: inp.value};
      });
    }
    return f;
  }

  function installNext() {
    const p = W.prep;
    let blocked = false;
    const services = p.services.map((x, i) => {
      const c = x.check;
      const over = document.querySelector(`.gw-over[data-i="${i}"]`);
      if (x.error || !c || blocking(c) || (c.existing && !(over && over.checked))) blocked = true;
      return {app: x.app, ref: W.svcs[i].ref, path: x.path, env: W.svcs[i].env, files: [...document.querySelectorAll(`.gw-slot[data-i="${i}"]`)].map(slotSpec), unit: unitOf(i)};
    });
    $('gw-block').style.display = blocked ? '' : 'none';
    if (blocked) return;
    const overwrite = p.services.some(x => x.check && x.check.existing);
    const body = {kind: 'install', provision: {conn_id: W.conn, services, overwrite}, restart: {mode: $('gw-start').checked ? 'now' : 'none'},
      checks: Object.fromEntries([...document.querySelectorAll('.gw-chk')].map(x => [x.dataset.app, x.value.trim()]))};
    provisionConfirm(body, p, services.map(x => `<li><b>${esc(x.app)}</b> → <span class="g-mono">${esc(x.path)}</span>${x.unit ? ` · ${esc(x.unit.kind)} <span class="g-mono">${esc(x.unit.name)}</span>` : ''}</li>`).join(''), overwrite);
  }

  function provisionConfirm(body, p, list, overwrite) {
    W.body = body;
    const need = p.prod || overwrite;
    modal(`<h3>${esc(tf('git_confirm_title', {kind: kindName(body.kind)}))}</h3>
      <div>${esc(tf('git_w_confirm', {kind: kindName(body.kind), n: (body.provision.services || [1]).length, server: p.conn_name}))}${p.prod ? `<span class="g-prod">${esc(t('git_prod'))}</span>` : ''}</div><ul>${list}</ul>
      ${overwrite ? `<div class="hint-box warn">${esc(t('git_w_overwrite'))}</div>` : ''}
      ${body.restart.mode === 'now' ? `<div class="hint-box warn">${esc(t('git_w_start'))}</div>` : ''}
      ${need ? `<div class="g-conf"><label>${esc(tf('git_w_type', {name: p.conn_name}))}</label><input id="gw-conf" autocomplete="off" oninput="gitWorkspace.wCheckGo()"></div>` : ''}
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gw-go" class="btn-danger" ${need ? 'disabled' : ''}>${esc(t('git_start'))}</button></div>`);
    W.confName = p.conn_name;
    $('gw-go').onclick = async () => {
      const b = JSON.parse(JSON.stringify(W.body));
      if ($('gw-conf')) b.provision.confirm = $('gw-conf').value.trim();
      $('gw-go').disabled = true;
      const r = await fetch('/api/git/runs', json('POST', b));
      if (!r.ok) { $('gw-go').disabled = false; showToast(await apiError(r), 'error', 7000); return; }
      showRun((await r.json()).id);
    };
  }
  function wCheckGo() { $('gw-go').disabled = $('gw-conf').value.trim().toLowerCase() !== String(W.confName).trim().toLowerCase(); }

  function transfer(ids) {
    ids = pickIds(ids);
    if (!ids.length || !can('transfer')) return;
    const src = installById(ids[0]);
    W.kind = 'transfer'; W.src = src;
    const others = sshConns();
    modal(`<h3>${esc(t('git_k_transfer'))}: ${esc(src.app)} · ${esc(src.conn_name)} · <span class="g-mono">${esc(src.path)}</span></h3>
      <div class="g-mut">${esc(t('git_w_tr_h'))}</div>
      <div class="g-grid" style="margin-top:8px"><div><label>${esc(t('git_w_server'))}</label><select id="gw-conn">${others.map(c => `<option value="${c.id}">${esc(c.name)} (${esc(c.host)})</option>`).join('')}</select></div>
        <div><label>${esc(t('git_path'))}</label><input id="gw-path" value="${esc(src.path)}"></div>
        <div><label>${esc(t('git_w_env'))}</label><input id="gw-env" value="${esc(src.env || '')}"></div></div>
      ${btns('gw-next', 'git_w_check')}`);
    const firstOther = others.find(c => c.id !== src.conn_id);
    if (firstOther) $('gw-conn').value = String(firstOther.id);
    $('gw-next').onclick = () => {
      W.conn = Number($('gw-conn').value); W.path = $('gw-path').value.trim(); W.env = $('gw-env').value.trim();
      if (W.conn === src.conn_id && W.path === src.path) { showToast(t('git_w_same'), 'warning'); return; }
      transferCheck();
    };
  }

  async function transferCheck() {
    modal(`<h3>${esc(t('git_k_transfer'))}</h3><div class="hint-box">⏳ ${esc(t('git_plan_loading'))}</div>`);
    const p = await prepare({kind: 'transfer', conn_id: W.conn, source_install: W.src.id, path: W.path, overwrite: true});
    if (!p) return;
    W.prep = p;
    const s = p.source, c = p.check;
    modal(`<h3>${esc(t('git_k_transfer'))}: ${esc(s.app)} → ${esc(p.conn_name)}${p.prod ? `<span class="g-prod">${esc(t('git_prod'))}</span>` : ''}</h3>
      <div class="g-card"><h3>${esc(t('git_w_source'))}: ${esc(s.conn_name)} · <span class="g-mono">${esc(s.path)}</span> ${s.version ? `<span class="g-mono">${esc(s.version)}</span>` : ''}</h3>
        <div class="g-mut">${esc(tf('git_w_files', {n: s.files, kb: kb(s.bytes)}))}</div>
        ${s.excluded ? `<details><summary class="g-mut">${esc(tf('git_w_excluded', {n: s.excluded}))}</summary><div class="g-mono g-mut">${s.excluded_list.map(x => esc(x.path) + ' <i>(' + esc(x.why) + ')</i>').join('<br>')}</div></details>` : ''}
        ${(s.units || []).length ? `<div class="g-mut">${esc(tf('git_w_units_not_copied', {list: s.units.map(u => u.name).join(', ')}))}</div>` : ''}</div>
      <div class="g-card"><h3>→ <span class="g-mono">${esc(W.path)}</span></h3>${checkHTML(c, 0)}</div>
      ${unitHTML(0, s.app, W.path)}${startHTML()}
      <div class="g-card"><label class="g-mut">${esc(tf('git_check_cmd', {app: s.app}))}</label><input class="gw-chk" data-app="${esc(s.app)}" value="${esc(((G.state.settings || {}).check_commands || {})[s.app] || '')}" style="width:100%"></div>
      <div id="gw-block" class="hint-box bad" style="display:none">${esc(t('git_w_blocked'))}</div>${btns('gw-next', 'git_continue')}`);
    $('gw-next').onclick = () => {
      const over = document.querySelector('.gw-over');
      const blocked = !c || blocking(c) || (c.existing && !(over && over.checked));
      $('gw-block').style.display = blocked ? '' : 'none';
      if (blocked) return;
      const body = {kind: 'transfer', provision: {conn_id: W.conn, source_install: W.src.id, path: W.path, env: W.env, unit: unitOf(0), overwrite: !!c.existing},
        restart: {mode: $('gw-start').checked ? 'now' : 'none'}, checks: {[s.app]: document.querySelector('.gw-chk').value.trim()}};
      provisionConfirm(body, p, `<li><b>${esc(s.app)}</b>: ${esc(s.conn_name)} <span class="g-mono">${esc(s.path)}</span> → ${esc(p.conn_name)} <span class="g-mono">${esc(W.path)}</span></li>`, !!c.existing);
    };
  }

  // ── .gitignore helper ──
  const I = {app: '', info: null, loading: false, preview: null, checked: null};
  const fmtSize = n => !n ? '—' : n < 1024 ? n + ' B' : n < 1048576 ? (n / 1024).toFixed(1) + ' KB' : (n / 1048576).toFixed(1) + ' MB';
  function allApps() { return [...G.state.catalog.apps.map(a => a.name), ...(G.state.bundle_apps || [])]; }

  async function ignoreLoad(app) {
    const apps = allApps();
    I.app = app || I.app || apps[0] || '';
    I.info = null; I.preview = null; I.checked = null;
    if (!I.app) { if (G.tab === 'gitignore') render(); return; }
    I.loading = true; if (G.tab === 'gitignore') render();
    try {
      const r = await fetch('/api/git/gitignore/' + encodeURIComponent(I.app));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      I.info = await r.json();
      // pre-selected: the patterns of the stacks found in the repository and the protected globs
      I.checked = new Set([...I.info.stacks.filter(x => x.detected).flatMap(x => x.patterns), ...I.info.protected]);
    } finally { I.loading = false; if (G.tab === 'gitignore') render(); }
  }

  function ignoreHTML() {
    const apps = allApps();
    const bar = `<div class="g-bar"><label class="g-mut">${esc(t('gi_service'))}</label><select onchange="gitWorkspace.ignoreLoad(this.value)">${apps.map(a => `<option value="${esc(a)}" ${a === I.app ? 'selected' : ''}>${esc(a)}</option>`).join('')}</select>
      <button class="btn-sec btn-sm" onclick="gitWorkspace.ignoreLoad()">↻ ${esc(t('gi_load'))}</button></div>`;
    if (!apps.length) return bar + `<div class="hint-box">${esc(t('git_empty_overview'))}</div>`;
    if (I.loading || !I.info) return bar + `<div class="hint-box">⏳ ${esc(t('gi_loading'))}</div>`;
    const x = I.info, chk = v => I.checked && I.checked.has(v) ? 'checked' : '';
    const pat = v => `<label><input type="checkbox" class="gi-p" value="${esc(v)}" ${chk(v)} onchange="gitWorkspace.ignoreMark(this)">${esc(v)}</label>`;
    let h = bar + `<div class="g-card"><div class="g-grid"><div><label>${esc(t('gi_repo'))}</label><span class="g-mono">${esc(x.repo || '—')}</span></div>
      <div><label>${esc(t('gi_branch'))}</label><span class="g-mono">${esc(x.branch || '—')}</span></div><div><label>${esc(t('gi_file'))}</label><span class="g-mono">${esc(x.file)}</span></div></div>
      ${x.note ? `<div class="hint-box warn" style="margin-top:8px">${esc(t('gi_note'))}: ${esc(x.note)}</div>` : ''}</div>`;
    if (x.warnings.length) h += `<div class="g-card g-warn"><h3>⚠ ${esc(t('gi_warnings'))}</h3>${x.warnings.map(w => `<div style="margin:6px 0"><span class="g-mono"><b>${esc(w.path)}</b></span> <span class="g-st s-review">${esc(t('gi_reason_' + w.reason))}</span>
      <div class="g-mut">${esc(tf('gi_warn_advice', {tmpl: w.template, path: w.path}))}</div></div>`).join('')}</div>`;
    h += `<div class="g-card"><h3>${esc(t('gi_extra'))}</h3>${x.extra.length ? `<div class="g-scroll" style="max-height:40vh"><table><tr><th></th><th>${esc(t('git_path'))}</th><th>${esc(t('gi_size'))}</th><th>${esc(t('git_server'))}</th><th>${esc(t('gi_pattern'))}</th></tr>
      ${x.extra.map(c => `<tr><td><input type="checkbox" class="gi-p" value="${esc(c.pattern)}" ${chk(c.pattern)} onchange="gitWorkspace.ignoreMark(this)"></td><td class="g-mono">${esc(c.path)}</td><td class="g-mut" style="white-space:nowrap">${esc(fmtSize(c.size))}</td>
        <td class="g-mut">${esc(c.servers.join(', '))}</td><td class="g-mono">${esc(c.pattern)}</td></tr>`).join('')}</table></div>` : `<div class="g-mut">${esc(t('gi_extra_none'))}</div>`}</div>`;
    h += `<div class="g-card"><h3>${esc(t('gi_stacks'))}</h3>${x.stacks.map(st => `<div style="margin:6px 0"><b>${esc(st.name)}</b>${st.detected ? ` <span class="g-st s-ok">${esc(tf('gi_detected', {f: st.reason}))}</span>` : ''}<div class="g-pats">${st.patterns.map(pat).join('')}</div></div>`).join('')}
      ${x.protected.length ? `<div style="margin:10px 0 6px"><b>${esc(t('gi_protected'))}</b><div class="g-pats">${x.protected.map(pat).join('')}</div></div>` : ''}
      <div class="g-grid" style="margin-top:8px"><div><label>${esc(t('gi_custom'))}</label><textarea id="gi-custom" placeholder="${esc(t('gi_custom_h'))}">${esc(I.custom || '')}</textarea></div></div></div>`;
    h += `<div class="g-card"><h3>${esc(t('gi_lists'))}</h3><div class="g-mut" style="margin-bottom:6px">${esc(t('gi_lists_h'))}</div>${x.in_catalog ? `<div class="g-grid">
      <div><label>${esc(t('git_f_exclude'))}</label><textarea id="gi-exclude">${esc(x.exclude.join('\n'))}</textarea></div>
      <div><label>${esc(t('git_f_protected'))}</label><textarea id="gi-protected">${esc(x.protected.join('\n'))}</textarea></div></div>
      <div class="g-bar" style="margin-top:8px"><button class="btn-sec btn-sm" onclick="gitWorkspace.ignoreSaveLists()">${esc(t('git_save'))}</button></div>` : `<div class="g-mut">${esc(t('gi_not_catalog'))}</div>`}</div>`;
    h += `<div class="g-card"><h3>${esc(t('gi_current'))}</h3>${x.has_current ? `<pre class="g-pre">${esc(x.current)}</pre>` : `<div class="g-mut">${esc(t('gi_no_current'))}</div>`}</div>`;
    const mrLabel = t(x.source_kind === 'github' ? 'gi_mr_gh' : 'gi_mr');
    h += `<div class="g-bar"><button onclick="gitWorkspace.ignorePreview()">${esc(t('gi_preview'))}</button><button class="btn-sec" onclick="gitWorkspace.ignoreDownload()">⬇ ${esc(t('gi_download'))}</button>
      <button class="btn-sec" ${x.can_mr ? '' : `disabled title="${esc(x.mr_note || '')}"`} onclick="gitWorkspace.ignoreMR()">${esc(mrLabel)}</button>${x.can_mr ? '' : ` <span class="g-mut">${esc(x.mr_note || '')}</span>`}</div>
      <div class="g-mut">${esc(t('gi_mr_default'))}</div><div id="gi-out">${I.preview ? previewHTML(I.preview) : ''}</div>`;
    return h;
  }

  function ignoreMark(el) { if (!I.checked) I.checked = new Set(); if (el.checked) I.checked.add(el.value); else I.checked.delete(el.value); document.querySelectorAll('.gi-p').forEach(b => { if (b.value === el.value) b.checked = el.checked; }); }
  function ignorePatterns() { I.custom = ($('gi-custom') || {}).value || ''; return [...new Set([...(I.checked || [])].concat(lines(I.custom.replace(/,/g, '\n'))))]; }

  function previewHTML(p) {
    if (!p.added.length) return `<div class="hint-box">${esc(t('gi_nothing'))}</div>`;
    return `<h3 style="margin:12px 0 6px">${esc(t('gi_diff'))} · ${esc(tf('gi_added', {n: p.added.length}))}</h3><div class="g-diff">` +
      p.lines.map(l => l.op === '@' ? `<div class="d-hunk">@@ -${l.a} +${l.b} @@</div>` : `<div class="${l.op === '+' ? 'd-add' : l.op === '-' ? 'd-del' : ''}">${esc(l.op + ' ' + l.s)}</div>`).join('') + '</div>';
  }

  async function ignorePreview() {
    const r = await fetch(`/api/git/gitignore/${encodeURIComponent(I.app)}/preview`, json('POST', {patterns: ignorePatterns()}));
    if (!r.ok) { showToast(await apiError(r), 'error'); return null; }
    I.preview = await r.json();
    $('gi-out').innerHTML = previewHTML(I.preview);
    $('gi-out').scrollIntoView({block: 'nearest'});
    return I.preview;
  }

  function ignoreDownload() { fetch(`/api/git/gitignore/${encodeURIComponent(I.app)}/download`, json('POST', {patterns: ignorePatterns()})).then(r => download(r, '.gitignore')); }

  async function ignoreSaveLists() {
    const a = G.state.catalog.apps.find(x => x.name === I.app);
    if (!a) return;
    const body = Object.assign({}, a, {exclude: lines($('gi-exclude').value), protected: lines($('gi-protected').value)});
    const r = await fetch('/api/git/catalog/apps/' + encodeURIComponent(I.app), json('PUT', body));
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    G.state = await r.json(); showToast(t('git_saved'), 'success'); ignoreLoad(I.app);
  }

  async function ignoreMR() {
    const x = I.info, p = await ignorePreview();
    if (!p || !p.added.length) return;
    const gh = x.source_kind === 'github';
    modal(`<h3>${esc(gh ? t('gi_mr_gh').replace('…', '') : tf('gi_mr_title', {repo: x.repo}))}</h3>
      <div>${esc(tf('gi_mr_q', {branch: 'wrm/gitignore-' + I.app.toLowerCase() + '-…', file: x.file, base: x.branch}))}</div>
      <div class="hint-box warn" style="margin:8px 0"><b class="g-mono">${esc(x.repo)}</b></div>
      <div class="g-mono g-mut">${p.added.map(esc).join('<br>')}</div>
      <div class="g-conf"><label>${esc(tf('gi_mr_type', {name: x.project}))}</label><input id="gi-conf" autocomplete="off"></div>
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gi-go" disabled>${esc(t('gi_mr_go'))}</button></div>`);
    $('gi-conf').oninput = () => { $('gi-go').disabled = $('gi-conf').value.trim() !== x.project; };
    $('gi-go').onclick = async () => {
      $('gi-go').disabled = true;
      const r = await fetch(`/api/git/gitignore/${encodeURIComponent(I.app)}/mr`, json('POST', {patterns: ignorePatterns(), confirm: $('gi-conf').value.trim()}));
      if (!r.ok) { $('gi-go').disabled = false; showToast(await apiError(r), 'error', 7000); return; }
      const res = await r.json();
      modal(`<h3>✅ ${esc(tf('gi_mr_done', {url: ''}))}</h3><div><a href="${esc(res.url)}" target="_blank" rel="noopener noreferrer" class="g-mono">${esc(res.url)}</a></div><div class="g-mut g-mono">${esc(res.branch)}</div>
        <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`);
    };
  }

  // ── CI: webhooks, artifact feeds, deploy method ──
  function hookURL(id) { return location.origin + '/api/hooks/git/' + id; }
  function webhook(id, secret) {
    const x = G.state.sources.find(s => s.id === id);
    if (!x) return;
    modal(`<h3>${esc(tf('git_wh_title', {name: x.name}))}</h3><div class="g-mut">${esc(t('git_wh_h'))}</div>
      ${x.hook_id ? `<div class="g-card"><label>${esc(t('git_wh_url'))}</label><input readonly value="${esc(hookURL(x.hook_id))}" onclick="this.select()" style="width:100%">
        ${secret ? `<label style="margin-top:8px">${esc(t('git_wh_secret'))}</label><input readonly value="${esc(secret)}" onclick="this.select()" style="width:100%">` : ''}
        <div class="g-mut" style="margin-top:6px">${esc(t('git_wh_last'))}: ${x.hook_at ? esc(fmtTime(x.hook_at)) : esc(t('git_never'))}</div></div>
        <ul class="g-mut"><li>${esc(t('git_wh_gitlab'))}</li><li>${esc(t('git_wh_github'))}</li><li>${esc(t('git_wh_generic'))}</li></ul>` : `<div class="hint-box">${esc(t('git_wh_off'))}</div>`}
      <div class="g-btns">${x.hook_id ? `<button class="btn-danger" onclick="gitWorkspace.hookSet(${id}, false)">${esc(t('git_wh_disable'))}</button><button class="btn-sec" onclick="gitWorkspace.hookSet(${id}, true)">${esc(t('git_wh_renew'))}</button>`
        : `<button onclick="gitWorkspace.hookSet(${id}, true)">${esc(t('git_wh_enable'))}</button>`}<button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`);
  }
  async function hookSet(id, on) {
    const r = await fetch(`/api/git/sources/${id}/hook`, {method: on ? 'POST' : 'DELETE'});
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const j = await r.json();
    G.state = on ? j.state : j; render();
    webhook(id, on ? j.secret : '');
  }

  function feedsHTML() {
    const feeds = G.state.feeds || [];
    return `<div class="g-card"><h3>${esc(t('git_feeds'))}</h3><div class="g-mut" style="margin-bottom:6px">${esc(t('git_feeds_h'))}</div>
      ${feeds.length ? `<table>${feeds.map(f => `<tr><td><b>${esc(f.name)}</b>${f.interval_minutes ? ` <span class="g-mut">⏱ ${f.interval_minutes} min</span>` : ''}<div class="g-mono g-mut">${esc(f.url)}</div>
        <div class="g-mut">${esc(t('git_feed_last'))}: ${f.last_ok_at ? esc(fmtTime(f.last_ok_at)) : esc(t('git_never'))}</div>${f.last_error ? `<div class="g-st s-error">${esc(f.last_error)}</div>` : ''}</td>
        <td style="white-space:nowrap"><button class="btn-sec btn-sm" onclick="gitWorkspace.feedFetch(${f.id})">${esc(t('git_feed_fetch'))}</button> <button class="btn-sec btn-sm" onclick="gitWorkspace.editFeed(${f.id})">${esc(t('git_edit'))}</button> <button class="btn-danger btn-sm" onclick="gitWorkspace.delFeed(${f.id})">✕</button></td></tr>`).join('')}</table>` : ''}
      <div class="g-bar" style="margin-top:8px"><button class="btn-sec" onclick="gitWorkspace.editFeed(0)">＋ ${esc(t('git_add_feed'))}</button></div></div>`;
  }
  function editFeed(id) {
    const f = (G.state.feeds || []).find(x => x.id === id) || {name: '', url: '', header_name: 'PRIVATE-TOKEN', interval_minutes: 0};
    modal(`<h3>${esc(id ? f.name : t('git_add_feed'))}</h3><div class="g-grid">
      <div><label>${esc(t('git_f_name'))}</label><input id="gf-name" value="${esc(f.name)}"></div>
      <div style="grid-column:1/-1"><label>${esc(t('git_feed_url'))}</label><input id="gf-url" value="${esc(f.url)}" placeholder="https://git.example.com/api/v4/projects/demo%2Fdemo-api/jobs/artifacts/main/raw/bundle.tar.gz?job=bundle"></div>
      <div><label>${esc(t('git_feed_header'))}</label><input id="gf-hname" value="${esc(f.header_name)}" placeholder="PRIVATE-TOKEN"></div>
      <div><label>${esc(t('git_feed_value'))}</label><input id="gf-hval" type="password" autocomplete="new-password" placeholder="${f.has_header ? esc(t('git_token_keep')) : ''}"></div>
      <div><label>${esc(t('git_feed_interval'))}</label><input id="gf-int" type="number" min="0" max="10080" value="${f.interval_minutes || 0}"></div></div>
      <div class="g-mut" style="margin-top:6px">${esc(t('git_feed_value_h'))}</div>
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gf-save">${esc(t('git_save'))}</button></div>`);
    $('gf-save').onclick = async () => {
      const body = {name: $('gf-name').value.trim(), url: $('gf-url').value.trim(), header_name: $('gf-hname').value.trim(), interval_minutes: Number($('gf-int').value) || 0};
      if ($('gf-hval').value || !id) body.header_value = $('gf-hval').value;
      const r = await fetch(id ? '/api/git/feeds/' + id : '/api/git/feeds', json(id ? 'PUT' : 'POST', body));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = await r.json(); closeModal(); render();
    };
  }
  async function feedFetch(id) {
    showToast('⏳ ' + t('git_feed_fetch'), 'info', 1500);
    const r = await fetch(`/api/git/feeds/${id}/fetch`, {method: 'POST'});
    if (!r.ok) { showToast(await apiError(r), 'error', 7000); load().then(ok => ok && render()); return; }
    const j = await r.json();
    G.state = j.state; render();
    showToast(t(j.status === 'imported' ? 'git_feed_imported' : 'git_feed_unchanged'), 'success');
  }
  async function delFeed(id) {
    if (!(await uiConfirm(t('git_delete') + '?', {danger: true, okText: t('git_delete')}))) return;
    const r = await fetch('/api/git/feeds/' + id, {method: 'DELETE'});
    if (r.ok) { G.state = await r.json(); render(); }
  }

  function pipelineForm(name) {
    if (!can('update')) return '';
    const p = (G.state.pipelines || []).find(x => x.app === name) || {kind: ''};
    const f = (k, label, v, ph, type) => `<div class="gp-f gp-${k}"><label>${esc(t(label))}</label><input id="gp-${k}" value="${esc(v || '')}" placeholder="${esc(ph || '')}" ${type ? `type="${type}" autocomplete="new-password"` : ''}></div>`;
    return `<div class="g-card"><h3>${esc(t('git_deploy_method'))}</h3><select id="gp-kind" onchange="gitWorkspace.pipelineKind()">
      <option value="">${esc(t('git_dm_files'))}</option><option value="jenkins" ${p.kind === 'jenkins' ? 'selected' : ''}>${esc(t('git_dm_jenkins'))}</option><option value="gitlab" ${p.kind === 'gitlab' ? 'selected' : ''}>${esc(t('git_dm_gitlab'))}</option></select>
      <div class="g-mut" style="margin:6px 0">${esc(t('git_dm_h'))}</div><div class="g-grid" id="gp-fields" data-had="${p.kind ? 1 : 0}">
      ${f('url', p.kind === 'gitlab' ? 'git_dm_gl_url' : 'git_dm_job', p.url, 'https://ci.example.com/job/deploy-demo')}${f('username', 'git_dm_user', p.username, 'ci-bot')}
      ${f('project', 'git_dm_project', p.project, 'group/project')}${f('ref', 'git_dm_ref', p.ref, 'main')}
      ${f('token', p.kind === 'gitlab' ? 'git_dm_trigger' : 'git_dm_token', '', p.has_token ? t('git_dm_keep') : '', 'password')}</div></div>`;
  }
  function pipelineKind() {
    const k = $('gp-kind').value;
    document.querySelectorAll('.gp-f').forEach(el => { el.style.display = !k ? 'none' : (el.classList.contains('gp-username') ? k === 'jenkins' : el.classList.contains('gp-project') || el.classList.contains('gp-ref') ? k === 'gitlab' : true) ? '' : 'none'; });
    if (k) { document.querySelector('.gp-url label').textContent = t(k === 'gitlab' ? 'git_dm_gl_url' : 'git_dm_job'); document.querySelector('.gp-token label').textContent = t(k === 'gitlab' ? 'git_dm_trigger' : 'git_dm_token'); }
  }
  async function savePipeline(name) {
    if (!$('gp-kind')) return true;
    const k = $('gp-kind').value, had = $('gp-fields').dataset.had === '1';
    if (!k) {
      if (!had) return true;
      const r = await fetch('/api/git/pipelines/' + encodeURIComponent(name), {method: 'DELETE'});
      if (r.ok) G.state = await r.json();
      return r.ok;
    }
    const body = {kind: k, url: $('gp-url').value.trim(), username: $('gp-username').value.trim(), project: $('gp-project').value.trim(), ref: $('gp-ref').value.trim()};
    if ($('gp-token').value || !had) body.token = $('gp-token').value;
    const r = await fetch('/api/git/pipelines/' + encodeURIComponent(name), json('PUT', body));
    if (!r.ok) { showToast(await apiError(r), 'error'); return false; }
    G.state = await r.json();
    return true;
  }

  // ── progress and history ──
  let runTimer = null;
  async function showRun(id) {
    clearTimeout(runTimer);
    const r = await fetch('/api/git/runs/' + id);
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const v = await r.json();
    const open = $('git-modal');
    if (open && open.dataset.run && open.dataset.run !== String(id)) return;
    const icon = s => ({ok: '✅', fail: '❌', skip: '⏭', info: 'ℹ'}[s] || '•');
    const html = `<h3>${esc(tf('git_progress', {label: v.label}))} · ${esc(kindName(v.kind))}${v.ref ? ' → <span class="g-mono">' + esc(v.ref) + '</span>' : ''} ${runBadge(v.state)}</h3>
      ${v.message ? `<div class="hint-box ${v.state === 'failed' ? 'bad' : ''}">${esc(v.message)}</div>` : ''}
      ${v.scheduled_at ? `<div class="g-mut">${esc(t('git_scheduled_for'))}: ${esc(fmtTime(v.scheduled_at))}</div>` : v.parent_id ? `<div class="g-mut">${esc(t('git_scheduled_for'))}: ${esc(t('git_waits_update'))}</div>` : ''}
      ${v.items.map(it => `<div class="g-item"><h4>${esc(it.conn_name)} · ${esc(it.app)} · <span class="g-mono">${esc(it.path)}</span> ${runBadge(it.state)}${it.step ? ` <span class="g-mut">⏳ ${esc(t('git_step_' + it.step))}</span>` : ''}</h4>
        ${it.source_path ? `<div class="g-mut">${esc(t('git_w_source'))}: ${esc(it.source_conn_name)} · <span class="g-mono">${esc(it.source_path)}</span></div>` : ''}
        ${it.from_version || it.to_version ? `<div class="g-mut g-mono">${esc(it.from_version || '')} → ${esc(it.to_version || '')}</div>` : ''}
        ${it.via === 'ci' ? `<div>🔧 ${esc(t('git_step_ci'))} (${esc(it.ci_kind || '')}): <b>${esc(it.ci_status || '…')}</b>${it.ci_url ? ` · <a href="${esc(it.ci_url)}" target="_blank" rel="noopener noreferrer">${esc(t('git_ci_link'))}</a>` : ''}</div>` : ''}
        ${it.error ? `<div class="hint-box bad">${esc(it.error)}</div>` : ''}
        ${(it.dangling || []).length ? `<div class="hint-box bad"><b>${esc(t('git_dangling'))}</b><div class="g-mono">${it.dangling.map(esc).join('<br>')}</div></div>` : ''}
        ${(it.health || []).length ? `<div><b>${esc(t('git_health'))}:</b> ${it.health.map(h => `${h.ok ? '✅' : '❌'} <span class="g-mono">${esc(h.name)}</span> ${esc(h.status)}`).join(' · ')}</div>` : ''}
        ${it.made_backup ? `<div class="g-mut">${esc(t('git_backup'))}: <span class="g-mono">${esc(it.made_backup)}</span></div>` : ''}
        ${(it.written || []).length || (it.deleted || []).length ? `<div class="g-mut">${esc(tf('git_written', {n: (it.written || []).length}))}${(it.deleted || []).length ? ' · ' + esc(tf('git_deleted', {n: it.deleted.length})) : ''}</div>` : ''}
        <ul class="g-log">${(it.log || []).map(l => `<li>${icon(l.status)} <b>${esc(t('git_step_' + l.step))}</b> ${esc(l.msg || '')}</li>`).join('')}</ul></div>`).join('')}
      <div class="g-btns">${v.running ? `<button class="btn-danger" onclick="gitWorkspace.cancelRun(${v.id})">${esc(t('git_stop'))}</button>` : v.state === 'scheduled' ? `<button class="btn-danger" onclick="gitWorkspace.cancelRun(${v.id})">${esc(t('git_cancel_run'))}</button>` : ''}
        <button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`;
    if (open && open.dataset.run === String(id)) open.querySelector('.g-dlg').innerHTML = html;
    else { modal(html); $('git-modal').dataset.run = String(id); }
    if (v.running || v.state === 'starting') runTimer = setTimeout(() => { if ($('git-modal') && $('git-modal').dataset.run === String(id)) showRun(id); }, 1000);
    else { load().then(ok => { if (ok && $('git-ws') && $('git-ws').style.display !== 'none') { render(); if (G.tab === 'runs') loadRuns(); } }); }
  }

  async function cancelRun(id) {
    if (!(await uiConfirm(t('git_cancel_q'), {danger: true, okText: t('git_cancel_run')}))) return;
    const r = await fetch(`/api/git/runs/${id}/cancel`, {method: 'POST'});
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    if ($('git-modal') && $('git-modal').dataset.run === String(id)) showRun(id);
    loadRuns();
  }

  async function loadRuns() {
    const r = await fetch('/api/git/runs');
    if (!r.ok) return;
    G.runs = (await r.json()).runs || [];
    if (G.tab === 'runs' && $('git-body')) $('git-body').innerHTML = runsHTML();
  }

  function runsHTML() {
    const list = G.runs;
    if (!list) return `<div class="hint-box">⏳</div>`;
    return `<div class="g-bar"><button class="btn-sec btn-sm" onclick="gitWorkspace.loadRuns()">↻ ${esc(t('git_refresh'))}</button></div>
      <div class="g-card g-scroll">${list.length ? `<table><tr><th>${esc(t('git_run'))}</th><th>${esc(t('git_kind_col'))}</th><th>${esc(t('git_state'))}</th><th>${esc(t('git_items'))}</th><th>${esc(t('git_versions'))}</th><th>${esc(t('git_created'))}</th><th>${esc(t('git_scheduled_for'))}</th><th></th></tr>
      ${list.map(v => `<tr><td class="g-mono">${esc(v.label)}<div class="g-mut">#${v.id}${v.parent_id ? ' ← #' + v.parent_id : ''}</div></td><td>${esc(kindName(v.kind))}${v.ref ? ` <span class="g-mono">${esc(v.ref)}</span>` : ''}</td>
        <td>${runBadge(v.state)}${v.message ? `<div class="g-mut">${esc(v.message)}</div>` : ''}</td>
        <td>${v.items.map(it => `<div>${runBadge(it.state)} ${esc(it.app)} @ ${esc(it.conn_name)}</div>`).join('')}</td>
        <td class="g-mono">${Object.entries(v.versions || {}).map(([a, x]) => esc(a + ' ' + x)).join('<br>')}</td>
        <td class="g-mut">${esc(fmtTime(v.created_at))}</td><td class="g-mut">${v.scheduled_at ? esc(fmtTime(v.scheduled_at)) : v.parent_id && v.state === 'scheduled' ? esc(t('git_waits_update')) : '—'}</td>
        <td style="white-space:nowrap"><button class="btn-sec btn-sm" onclick="gitWorkspace.showRun(${v.id})">${esc(t('git_details'))}</button>${v.state === 'scheduled' || v.running ? ` <button class="btn-danger btn-sm" onclick="gitWorkspace.cancelRun(${v.id})">✕</button>` : ''}</td></tr>`).join('')}</table>` : `<div class="g-mut">${esc(t('git_runs_empty'))}</div>`}</div>`;
  }

  function modal(html) {
    closeModal();
    const m = document.createElement('div');
    m.id = 'git-modal';
    m.innerHTML = `<div class="g-dlg">${html}</div>`;
    m.addEventListener('mousedown', e => { if (e.target === m) closeModal(); });
    m.addEventListener('keydown', e => { if (e.key === 'Escape') { e.stopPropagation(); closeModal(); } });
    document.body.appendChild(m);
    const first = m.querySelector('input, select, textarea');
    if (first) first.focus();
  }
  function closeModal() { const m = $('git-modal'); if (m) m.remove(); clearTimeout(runTimer); }

  let qTimer = null;
  window.gitWorkspace = {
    open, close, tab: setTab, check, refresh, details, diff, forget, sel, clearSel, setStale, checkIds, checkSel, checkVisible, checkStale, cancelCheck,
    viewFile, fileFilter, fileSearch, diffMode: setMode, addInstall, editApp, delApp, suggest, importCatalog, exportCatalog, exportBundle,
    update, upgrade, rollback, restart, stamp, install, transfer, wSvcRef, wSlotMode, wUnit, wCheckGo, localWarn, restartWarn, confirmCheck, showRun, cancelRun, loadRuns,
    ignoreLoad, ignoreMark, ignorePreview, ignoreDownload, ignoreSaveLists, ignoreMR, webhook, hookSet, editFeed, feedFetch, delFeed, pipelineKind,
    editSource, testSource, delSource, importBundle, saveSettings, saveCatalogSettings, setRef, filterConns, selConns, closeModal,
    filter(k, v, debounce) { G.f[k] = v; clearTimeout(qTimer); if (debounce) qTimer = setTimeout(() => { render(); const i = document.querySelector('#git-body input[type=search]'); if (i) { i.focus(); i.setSelectionRange(i.value.length, i.value.length); } }, 250); else render(); },
    toggleChanges(on) { G.onlyChanges = on; if (G.detail) details(G.detail.id); },
  };
})();
