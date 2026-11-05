/* Copyright 2026 vedranius and WRM contributors
   SPDX-License-Identifier: Apache-2.0 */

// Chains: live monitoring of a service chain (chains.go). Loaded on demand by openChains() in
// index.html; uses its helpers (t, esc, showToast). A live SVG diagram (automatic layout, dragged
// positions saved, animated flows, node/link details, a scenario band over time) and a list
// view, generic starter templates, and a JSON editor with validation, import and export.
'use strict';

Object.assign(LANGS.en, {
  chn_title: 'Chains', chn_new: 'New chain', chn_none: 'No chains yet. A chain links the servers of one system and checks them together.',
  chn_pick: 'Choose a chain on the left, or create one.', chn_check_now: 'Check now', chn_checking: 'Checking…', chn_edit: 'Edit JSON',
  chn_export: 'Export', chn_import: 'Import', chn_delete: 'Delete', chn_delete_q: 'Delete the chain {name}? Its timeline is deleted too; the servers are not touched.',
  chn_scenario: 'Scenario', chn_no_scenario: 'no scenario matches', chn_not_checked: 'not checked yet', chn_last: 'Last check',
  chn_nodes: 'Nodes', chn_checks: 'Checks', chn_links: 'Links', chn_scenarios: 'Scenarios', chn_timeline: 'Changes',
  chn_timeline_none: 'No changes recorded yet.', chn_up: 'up', chn_down: 'down', chn_disabled: 'paused',
  chn_kind_server: 'server', chn_kind_storage: 'storage', chn_kind_external: 'outside endpoint',
  chn_editor_t: 'Chain definition (JSON)', chn_validate: 'Validate', chn_valid: 'The chain is valid.', chn_save: 'Save', chn_cancel: 'Cancel',
  chn_saved: 'Saved. The first check runs within a minute.', chn_enabled: 'Checked in the background', chn_help: 'Reference',
  chn_help_conns: 'Your SSH connections (conn_id)', chn_help_storages: 'Your storages (storage_id)', chn_help_types: 'Check types',
  chn_help_rules: 'Rules: {"all":[…]}, {"any":[…]}, {"check":"id","field":"value|ok","op":"==|!=|in|contains|exists","value":"…"} or {"node":"id","op":"==","value":"up"}. The first scenario that matches is active; an empty rule {} matches always.',
  chn_import_bad: 'The file is not a chain export.', chn_failed: 'The request failed.', chn_ok_count: '{ok}/{n} checks OK',
  chn_nodes_count: '{up}/{n} servers up',
  chn_type_systemd: "systemd unit is active (target: unit name)",
  chn_type_timer: "systemd timer is active (target: timer name)",
  chn_type_docker: "Docker container is running (target: container name)",
  chn_type_process: "number of processes whose command line matches (target: pattern)",
  chn_type_port: "a TCP port is listening on the server (port)",
  chn_type_file: "file exists, optionally not older than max_age_s (target: absolute path)",
  chn_type_dir_count: "number of files in a directory (target: absolute path)",
  chn_type_config: "one value from a configuration file (format: ini, json, yaml, env, kv or regex; key: section.key, a.b.c or a regular expression)",
  chn_type_mount: "a mount point is mounted, read-write or read-only (target: mount point)",
  chn_type_cron: "a crontab line with the text is enabled, disabled (commented out) or missing (target: text)",
  chn_type_monitor: "the state of WRM's own monitoring for the server",
  chn_type_dns: "DNS answer for a name (target: name; record: A, AAAA, CNAME, TXT, MX, NS; resolver: optional ip[:port])",
  chn_type_tcp: "a TCP port answers, checked from WRM (target: host:port)",
  chn_type_http: "HTTP(S) status, or a JSON value with key (target: URL)",
  chn_type_tls: "days until the TLS certificate expires (target: host[:port]; warn_days, default 14)",
  chn_type_storage: "the storage's bucket answers; with prefix, the newest object under it (and max_age_s)",
});
Object.assign(LANGS.hr, {
  chn_title: 'Lanci', chn_new: 'Novi lanac', chn_none: 'Nema lanaca. Lanac povezuje servere jednog sustava i provjerava ih zajedno.',
  chn_pick: 'Odaberite lanac lijevo ili napravite novi.', chn_check_now: 'Provjeri sada', chn_checking: 'Provjeravam…', chn_edit: 'Uredi JSON',
  chn_export: 'Izvezi', chn_import: 'Uvezi', chn_delete: 'Obriši', chn_delete_q: 'Obrisati lanac {name}? Briše se i njegova povijest promjena; serveri se ne diraju.',
  chn_scenario: 'Scenarij', chn_no_scenario: 'nijedan scenarij ne odgovara', chn_not_checked: 'još nije provjereno', chn_last: 'Zadnja provjera',
  chn_nodes: 'Čvorovi', chn_checks: 'Provjere', chn_links: 'Veze', chn_scenarios: 'Scenariji', chn_timeline: 'Promjene',
  chn_timeline_none: 'Još nema zabilježenih promjena.', chn_up: 'radi', chn_down: 'ne radi', chn_disabled: 'pauzirano',
  chn_kind_server: 'server', chn_kind_storage: 'storage', chn_kind_external: 'vanjski endpoint',
  chn_editor_t: 'Definicija lanca (JSON)', chn_validate: 'Provjeri', chn_valid: 'Lanac je ispravan.', chn_save: 'Spremi', chn_cancel: 'Odustani',
  chn_saved: 'Spremljeno. Prva provjera kreće unutar minute.', chn_enabled: 'Provjerava se u pozadini', chn_help: 'Pomoć',
  chn_help_conns: 'Vaše SSH veze (conn_id)', chn_help_storages: 'Vaši storagei (storage_id)', chn_help_types: 'Vrste provjera',
  chn_help_rules: 'Pravila: {"all":[…]}, {"any":[…]}, {"check":"id","field":"value|ok","op":"==|!=|in|contains|exists","value":"…"} ili {"node":"id","op":"==","value":"up"}. Aktivan je prvi scenarij koji odgovara; prazno pravilo {} uvijek odgovara.',
  chn_import_bad: 'Datoteka nije izvoz lanca.', chn_failed: 'Zahtjev nije uspio.', chn_ok_count: '{ok}/{n} provjera u redu',
  chn_nodes_count: '{up}/{n} servera radi',
  chn_type_systemd: "systemd unit je aktivan (target: naziv unita)",
  chn_type_timer: "systemd timer je aktivan (target: naziv timera)",
  chn_type_docker: "Docker kontejner radi (target: naziv kontejnera)",
  chn_type_process: "broj procesa čija naredba odgovara uzorku (target: uzorak)",
  chn_type_port: "TCP port sluša na serveru (port)",
  chn_type_file: "datoteka postoji, po želji ne starija od max_age_s (target: apsolutna putanja)",
  chn_type_dir_count: "broj datoteka u direktoriju (target: apsolutna putanja)",
  chn_type_config: "jedna vrijednost iz konfiguracijske datoteke (format: ini, json, yaml, env, kv ili regex; key: section.key, a.b.c ili regularni izraz)",
  chn_type_mount: "mount point je montiran, za čitanje i pisanje ili samo za čitanje (target: mount point)",
  chn_type_cron: "crontab linija s tekstom je uključena, isključena (zakomentirana) ili ne postoji (target: tekst)",
  chn_type_monitor: "stanje WRM-ovog nadzora za server",
  chn_type_dns: "DNS odgovor za ime (target: ime; record: A, AAAA, CNAME, TXT, MX, NS; resolver: po želji ip[:port])",
  chn_type_tcp: "TCP port odgovara, provjereno iz WRM-a (target: host:port)",
  chn_type_http: "HTTP(S) status ili JSON vrijednost s key (target: URL)",
  chn_type_tls: "dana do isteka TLS certifikata (target: host[:port]; warn_days, zadano 14)",
  chn_type_storage: "bucket storagea odgovara; s prefix, najnoviji objekt ispod njega (i max_age_s)",
});

// Step 2: the live diagram, templates and the scenario timeline.
Object.assign(LANGS.en, {
  chn_view_diagram: 'Diagram', chn_view_list: 'List', chn_zoom_in: 'Zoom in', chn_zoom_out: 'Zoom out', chn_fit: 'Fit to window',
  chn_auto_layout: 'Arrange automatically', chn_auto_layout_q: 'Forget the saved positions and arrange the diagram automatically?',
  chn_drag_hint: 'Drag nodes to arrange them; the positions are saved. Click a node or a link for details.',
  chn_pos_failed: 'The positions were not saved.', chn_summary: 'Overview', chn_link: 'Link',
  chn_link_idle: 'idle in this scenario', chn_link_active_in: 'Carries traffic in', chn_link_check: 'Coloured by check',
  chn_history: 'History', chn_no_history: 'No changes recorded for this yet.', chn_definition: 'Definition',
  chn_conn: 'SSH connection', chn_storage: 'Storage', chn_band_t: 'Scenario over time', chn_band_none: 'No scenario recorded yet.',
  chn_now: 'now', chn_since: 'since {at}', chn_flow: 'Traffic flows', chn_legend_idle: 'idle', chn_legend_bad: 'failing',
  chn_lk_data: 'data', chn_lk_read: 'reads', chn_lk_write: 'writes', chn_lk_failover: 'failover', chn_lk_dependency: 'depends on',
  chn_tpl_pick: 'Start from a template', chn_tpl_empty: 'Empty chain', chn_tpl_empty_d: 'Write the definition in JSON yourself.',
  chn_tpl_single: 'One service', chn_tpl_single_d: 'One server with a service and its port, and a public HTTPS check.',
  chn_tpl_failover: 'Primary and backup path', chn_tpl_failover_d: 'Ingest on a primary and a backup server, a content server writing to a bucket, and a CDN in front.',
  chn_tpl_name: 'Name', chn_tpl_servers: 'Choose the server for each node', chn_tpl_storages: 'Choose the storage for each storage node',
  chn_tpl_choose: '— choose —', chn_tpl_no_ssh: 'You have no SSH connections yet. Add one first.',
  chn_tpl_no_storage: 'You have no storages yet; the storage nodes will simply be left out. A chain works without one.',
  chn_tpl_missing: 'Choose a server for every server node.', chn_tpl_after: 'Service names, paths and addresses are examples: adjust them in the guided editor after creating the chain.',
  chn_create: 'Create', chn_back: 'Back', chn_nodes_n: '{n} nodes',
  chn_help_links: 'Links: {"from":"a","to":"b","kind":"data|read|write|failover|dependency","check":"id","scenarios":["primary"]}. With scenarios, the link carries traffic only in those scenarios; otherwise the diagram shows it idle.',
});
Object.assign(LANGS.hr, {
  chn_view_diagram: 'Dijagram', chn_view_list: 'Popis', chn_zoom_in: 'Povećaj', chn_zoom_out: 'Smanji', chn_fit: 'Prilagodi prozoru',
  chn_auto_layout: 'Rasporedi automatski', chn_auto_layout_q: 'Zaboraviti spremljene položaje i automatski rasporediti dijagram?',
  chn_drag_hint: 'Povucite čvorove da ih rasporedite; položaji se spremaju. Kliknite čvor ili vezu za detalje.',
  chn_pos_failed: 'Položaji nisu spremljeni.', chn_summary: 'Pregled', chn_link: 'Veza',
  chn_link_idle: 'neaktivna u ovom scenariju', chn_link_active_in: 'Promet nosi u', chn_link_check: 'Boju određuje provjera',
  chn_history: 'Povijest', chn_no_history: 'Za ovo još nema zabilježenih promjena.', chn_definition: 'Definicija',
  chn_conn: 'SSH veza', chn_storage: 'Storage', chn_band_t: 'Scenarij kroz vrijeme', chn_band_none: 'Još nema zabilježenog scenarija.',
  chn_now: 'sada', chn_since: 'od {at}', chn_flow: 'Promet teče', chn_legend_idle: 'neaktivno', chn_legend_bad: 'ne radi',
  chn_lk_data: 'podaci', chn_lk_read: 'čita', chn_lk_write: 'piše', chn_lk_failover: 'rezervni put', chn_lk_dependency: 'ovisi o',
  chn_tpl_pick: 'Počnite od predloška', chn_tpl_empty: 'Prazan lanac', chn_tpl_empty_d: 'Definiciju napišite sami u JSON-u.',
  chn_tpl_single: 'Jedan servis', chn_tpl_single_d: 'Jedan server sa servisom i njegovim portom te javna HTTPS provjera.',
  chn_tpl_failover: 'Primarni i rezervni put', chn_tpl_failover_d: 'Ingest na primarnom i rezervnom serveru, content server koji piše u bucket i CDN ispred.',
  chn_tpl_name: 'Naziv', chn_tpl_servers: 'Odaberite server za svaki čvor', chn_tpl_storages: 'Odaberite storage za svaki storage čvor',
  chn_tpl_choose: '— odaberite —', chn_tpl_no_ssh: 'Još nemate SSH veza. Najprije dodajte jednu.',
  chn_tpl_no_storage: 'Još nemate storagea; storage čvorovi će jednostavno biti izostavljeni. Lanac radi i bez njih.',
  chn_tpl_missing: 'Odaberite server za svaki serverski čvor.', chn_tpl_after: 'Nazivi servisa, putanje i adrese su primjeri: prilagodite ih u vođenom uređivaču nakon stvaranja lanca.',
  chn_create: 'Napravi', chn_back: 'Natrag', chn_nodes_n: '{n} čvorova',
  chn_help_links: 'Veze: {"from":"a","to":"b","kind":"data|read|write|failover|dependency","check":"id","scenarios":["primary"]}. Sa scenarios veza nosi promet samo u tim scenarijima; inače je dijagram prikazuje neaktivnom.',
});

// Step 3: drafts written by the AI assistant and approved by the user.
Object.assign(LANGS.en, {
  chn_drafts: 'Drafts from the AI', chn_draft_badge: 'Draft — not monitored yet',
  chn_draft_note_t: 'Review this draft.',
  chn_draft_note: 'The assistant wrote it from your description. Nothing is checked in the background and no command runs until you create the chain; press Edit to adjust a server, check or link first.',
  chn_approve: 'Create chain', chn_approve_q: 'Create the chain {name} from this draft and start its background checks?',
  chn_discard: 'Discard', chn_discard_q: 'Discard the draft {name}? The assistant can write it again.',
  chn_draft_approved: 'The chain was created. The first check runs within a minute.',
  chn_draft_untitled: 'Untitled draft', chn_draft_counts: '{n} nodes · {c} checks',
});
Object.assign(LANGS.hr, {
  chn_drafts: 'Skice iz AI asistenta', chn_draft_badge: 'Skica — još se ne provjerava',
  chn_draft_note_t: 'Pregledajte ovu skicu.',
  chn_draft_note: 'Asistent ju je napisao prema vašem opisu. Ništa se ne provjerava u pozadini i nijedna naredba se ne izvršava dok ne stvorite lanac; ako je neki servis, putanja ili veza pogrešna, pritisnite Uredi i prilagodite ju.',
  chn_approve: 'Stvori lanac', chn_approve_q: 'Stvoriti lanac {name} iz ove skice i pokrenuti provjere u pozadini?',
  chn_discard: 'Odbaci', chn_discard_q: 'Odbaciti skicu {name}? Asistent ju može ponovno napisati.',
  chn_draft_approved: 'Lanac je stvoren. Prva provjera kreće unutar minute.',
  chn_draft_untitled: 'Skica bez naziva', chn_draft_counts: 'Čvorovi: {n} · Provjere: {c}',
});

// Step 4: the guided builder, the visual connections and the AI wizard.
Object.assign(LANGS.en, {
  chn_build: 'Build a chain', chn_build_ai: 'Build with AI', chn_build_ai_d: 'Choose the servers, describe the system in one sentence; WRM prepares a draft you review on the diagram.',
  chn_build_manual: 'Build step by step', chn_build_manual_d: 'Pick servers, add checks and connect them yourself. No JSON.',
  chn_build_tpl: 'Start from a template', chn_build_adv: 'Advanced: write JSON',
  chn_how: 'How a chain works', chn_how_t: 'Nodes are your servers, storages and outside services. Checks test each one (a service, a port, a file, a certificate). Links show what talks to what — a failover link is the backup path. Scenarios name the state: the first rule that matches wins, and the last empty rule is the fallback.',
  chn_ai_note: 'WRM first reads the selected servers (OS, services, containers, ports and mounts) with fixed read-only checks, then prepares a draft. Nothing is checked until you press Create chain.',
  chn_ai_analyze: 'Look at the servers', chn_ai_analyzing: 'Looking…', chn_ai_found_services: 'Services', chn_ai_found_ports: 'Ports', chn_ai_found_mounts: 'Mounts', chn_ai_no_facts: 'Nothing found (or the server did not answer).', chn_ai_storage_optional: 'Storages are optional. Leave them out and WRM only observes what is on the servers and any storage they mount; it never attaches to a storage or asks for its credentials.',
  chn_ai_pick_servers: 'Servers in this system', chn_ai_pick_storages: 'Storages (optional)',
  chn_ai_desc: 'What does the system do?', chn_ai_desc_ph: 'For example: ingest on this server writes recordings to the bucket; the backup server takes over if the primary fails; a CDN serves the content.',
  chn_ai_name: 'Name (optional)', chn_ai_make: 'Prepare the draft', chn_ai_working: 'Preparing…',
  chn_ai_servers_none: 'No SSH connections yet. Add a connection first.',
  chn_ai_ready: 'The draft is ready. Review it on the diagram, adjust it, then create the chain.',
  chn_ai_fallback: 'AI is not configured for this account, so WRM built a basic chain from your servers. Adjust the checks on the diagram.',
  chn_ai_failed: 'The AI did not return a usable chain, so WRM built a basic one. Adjust it on the diagram.',
  chn_g_name: 'Name', chn_g_desc: 'Description', chn_g_interval: 'Check every', chn_g_confirm: 'A change counts after',
  chn_g_seconds: '{n} s', chn_g_confirm_n: '{n} checks', chn_g_enabled: 'Check in the background',
  chn_g_basics: 'Name and behaviour', chn_g_nodes: 'Servers, storages and outside services', chn_g_add: 'Add',
  chn_g_kind: 'Kind', chn_g_pick_server: 'Choose a server', chn_g_pick_storage: 'Choose a storage',
  chn_g_label: 'Name', chn_g_role: 'Role', chn_g_role_ph: 'app, backup, database, cdn…', chn_g_remove: 'Remove',
  chn_g_node_empty: 'Add at least one server, storage or outside service.', chn_g_checks: 'What to check',
  chn_g_add_check: 'Add a check', chn_g_no_checks: 'No checks yet. A server with no check still shows whether SSH answers.',
  chn_g_edit: 'Edit', chn_g_type: 'Check', chn_g_target: 'Target', chn_g_port: 'Port', chn_g_format: 'Format',
  chn_g_key: 'Key', chn_g_record: 'DNS record', chn_g_resolver: 'Resolver (optional)', chn_g_prefix: 'Prefix (optional)',
  chn_g_max_age: 'Max age (s, optional)', chn_g_warn_days: 'Warn before expiry (days)', chn_g_expect: 'Also OK when the value',
  chn_g_expect_none: 'Use the default meaning of OK', chn_g_check_label: 'Label', chn_g_add_btn: 'Add', chn_g_save: 'Save',
  chn_g_links: 'What talks to what', chn_g_add_link: 'Add a link', chn_g_from: 'From', chn_g_to: 'To',
  chn_g_link_kind: 'Meaning', chn_g_link_check: 'Colour by check (optional)',
  chn_g_link_none: 'No links yet. Links are for the picture and the backup path; checks work without them.',
  chn_g_scenarios: 'States (scenarios)', chn_g_add_scenario: 'Add a state', chn_g_auto_scenarios: 'Create the usual states automatically',
  chn_g_auto_scenarios_d: 'Normal: every check is OK. Backup: the first server is down and the second is up. Degraded: everything else.',
  chn_g_when: 'Active when', chn_g_cond_add: 'Add a condition', chn_g_cond_none: 'Always (fallback)',
  chn_g_subject: 'Subject', chn_g_value: 'Value', chn_g_up: 'Up', chn_g_down: 'Down', chn_g_ok: 'OK', chn_g_op: 'Operator',
  chn_g_save_chain: 'Save chain', chn_g_save_draft: 'Save draft', chn_g_json: 'Advanced JSON', chn_g_cancel: 'Cancel',
  chn_g_errs: 'Fix these problems:', chn_g_required: 'Give the chain a name and add at least one node.',
  chn_g_any_note: 'This state uses "any" rules; edit it in Advanced JSON.',
  chn_g_edit_chain: 'Edit', chn_g_edit_draft: 'Edit draft', chn_g_new: 'New chain',
  chn_add_node: 'Add node', chn_connect: 'Connect to…', chn_delete_node: 'Delete node', chn_delete_link: 'Delete link',
  chn_add_check: 'Add check', chn_edit_chain: 'Edit chain', chn_pick_target: 'Click the target node',
  chn_link_saved: 'Saved.', chn_link_exists: 'These two nodes are already linked.', chn_check_bad: 'not OK',
});
Object.assign(LANGS.hr, {
  chn_build: 'Izgradi lanac', chn_build_ai: 'Izgradi uz AI', chn_build_ai_d: 'Odaberite servere i jednom rečenicom opišite sustav; WRM pripremi skicu koju pregledate na dijagramu.',
  chn_build_manual: 'Izgradi korak po korak', chn_build_manual_d: 'Sami odaberite servere, dodajte provjere i povežite ih. Bez JSON-a.',
  chn_build_tpl: 'Počni od predloška', chn_build_adv: 'Napredno: napiši JSON',
  chn_how: 'Kako lanac radi', chn_how_t: 'Čvorovi su vaši serveri, storagei i vanjske usluge. Provjere testiraju svaki (servis, port, datoteku, certifikat). Veze pokazuju što s čim razgovara — failover veza je rezervni put. Scenariji imenuju stanje: prvo pravilo koje odgovara pobjeđuje, a zadnje prazno pravilo je rezerva.',
  chn_ai_note: 'WRM najprije pročita odabrane servere (OS, servise, kontejnere, portove i mountove) fiksnim provjerama samo za čitanje, zatim pripremi skicu. Ništa se ne provjerava dok ne pritisnete Stvori lanac.',
  chn_ai_analyze: 'Pregledaj servere', chn_ai_analyzing: 'Pregledavam…', chn_ai_found_services: 'Servisi', chn_ai_found_ports: 'Portovi', chn_ai_found_mounts: 'Mountovi', chn_ai_no_facts: 'Ništa nije pronađeno (ili se server nije javio).', chn_ai_storage_optional: 'Storagei nisu obavezni. Ako ih izostavite, WRM samo promatra što je na serverima i što oni mountaju; nikada se ne spaja na storage i ne traži njegove vjerodajnice.',
  chn_ai_pick_servers: 'Serveri u ovom sustavu', chn_ai_pick_storages: 'Storagei (nije obavezno)',
  chn_ai_desc: 'Što sustav radi?', chn_ai_desc_ph: 'Na primjer: ingest na ovom serveru piše snimke u bucket; rezervni server preuzima ako primarni padne; CDN poslužuje sadržaj.',
  chn_ai_name: 'Naziv (nije obavezno)', chn_ai_make: 'Pripremi skicu', chn_ai_working: 'Pripremam…',
  chn_ai_servers_none: 'Još nemate SSH veza. Najprije dodajte vezu.',
  chn_ai_ready: 'Skica je gotova. Pregledajte ju na dijagramu, prilagodite i zatim stvorite lanac.',
  chn_ai_fallback: 'AI nije konfiguriran za ovaj račun, pa je WRM sastavio osnovni lanac od odabranih servera. Prilagodite provjere na dijagramu.',
  chn_ai_failed: 'AI nije vratio ispravan lanac, pa je WRM sastavio osnovni. Prilagodite ga na dijagramu.',
  chn_g_name: 'Naziv', chn_g_desc: 'Opis', chn_g_interval: 'Provjeravaj svakih', chn_g_confirm: 'Promjena se broji nakon',
  chn_g_seconds: '{n} s', chn_g_confirm_n: '{n} provjera', chn_g_enabled: 'Provjeravaj u pozadini',
  chn_g_basics: 'Naziv i ponašanje', chn_g_nodes: 'Serveri, storagei i vanjske usluge', chn_g_add: 'Dodaj',
  chn_g_kind: 'Vrsta', chn_g_pick_server: 'Odaberite server', chn_g_pick_storage: 'Odaberite storage',
  chn_g_label: 'Naziv', chn_g_role: 'Uloga', chn_g_role_ph: 'app, backup, database, cdn…', chn_g_remove: 'Ukloni',
  chn_g_node_empty: 'Dodajte barem jedan server, storage ili vanjsku uslugu.', chn_g_checks: 'Što provjeravati',
  chn_g_add_check: 'Dodaj provjeru', chn_g_no_checks: 'Još nema provjera. Server bez provjere i dalje pokazuje odgovara li SSH.',
  chn_g_edit: 'Uredi', chn_g_type: 'Provjera', chn_g_target: 'Cilj', chn_g_port: 'Port', chn_g_format: 'Format',
  chn_g_key: 'Ključ', chn_g_record: 'DNS zapis', chn_g_resolver: 'Resolver (nije obavezno)', chn_g_prefix: 'Prefiks (nije obavezno)',
  chn_g_max_age: 'Najveća starost (s, nije obavezno)', chn_g_warn_days: 'Upozori prije isteka (dana)', chn_g_expect: 'Također u redu kad je vrijednost',
  chn_g_expect_none: 'Koristi zadano značenje "u redu"', chn_g_check_label: 'Naziv', chn_g_add_btn: 'Dodaj', chn_g_save: 'Spremi',
  chn_g_links: 'Što s čim razgovara', chn_g_add_link: 'Dodaj vezu', chn_g_from: 'Od', chn_g_to: 'Prema',
  chn_g_link_kind: 'Značenje', chn_g_link_check: 'Boju određuje provjera (nije obavezno)',
  chn_g_link_none: 'Još nema veza. Veze služe slici i rezervnom putu; provjere rade i bez njih.',
  chn_g_scenarios: 'Stanja (scenariji)', chn_g_add_scenario: 'Dodaj stanje', chn_g_auto_scenarios: 'Automatski stvori uobičajena stanja',
  chn_g_auto_scenarios_d: 'Normal: sve provjere su u redu. Backup: prvi server ne radi, drugi radi. Degraded: sve ostalo.',
  chn_g_when: 'Aktivno kad', chn_g_cond_add: 'Dodaj uvjet', chn_g_cond_none: 'Uvijek (rezerva)',
  chn_g_subject: 'Predmet', chn_g_value: 'Vrijednost', chn_g_up: 'Radi', chn_g_down: 'Ne radi', chn_g_ok: 'U redu', chn_g_op: 'Operator',
  chn_g_save_chain: 'Spremi lanac', chn_g_save_draft: 'Spremi skicu', chn_g_json: 'Napredni JSON', chn_g_cancel: 'Odustani',
  chn_g_errs: 'Ispravite sljedeće:', chn_g_required: 'Dajte lancu naziv i dodajte barem jedan čvor.',
  chn_g_any_note: 'Ovo stanje koristi "any" pravila; uredite ga u naprednom JSON-u.',
  chn_g_edit_chain: 'Uredi', chn_g_edit_draft: 'Uredi skicu', chn_g_new: 'Novi lanac',
  chn_add_node: 'Dodaj čvor', chn_connect: 'Poveži s…', chn_delete_node: 'Obriši čvor', chn_delete_link: 'Obriši vezu',
  chn_add_check: 'Dodaj provjeru', chn_edit_chain: 'Uredi lanac', chn_pick_target: 'Kliknite ciljni čvor',
  chn_link_saved: 'Spremljeno.', chn_link_exists: 'Ova dva čvora su već povezana.', chn_check_bad: 'nije u redu',
});

(function () {
  const state = { list: [], drafts: [], current: null, timeline: [], styleAdded: false, mode: 'diagram', sel: null, views: {}, pos: {}, openCheck: null, conns: [], storages: [], catalog: null };
  const $ = (sel) => document.querySelector(sel);
  const el = (tag, cls, text) => { const e = document.createElement(tag); if (cls) e.className = cls; if (text != null) e.textContent = text; return e; };
  const fmt = (s, v) => s.replace(/\{(\w+)\}/g, (_, k) => (v[k] != null ? String(v[k]) : ''));
  const when = (at) => {
    if (!at) return t('chn_not_checked');
    const d = new Date(at);
    return isNaN(d) ? at.slice(0, 19).replace('T', ' ') : d.toLocaleString();
  };
  const trunc = (s, n) => { s = String(s || ''); return s.length > n ? s.slice(0, n - 1) + '…' : s; };
  const NW = 184, NH = 78, GX = 250, GY = 112;
  const SCN_COLORS = ['#3b82f6', '#f59e0b', '#8b5cf6', '#10b981', '#ec4899', '#14b8a6', '#eab308', '#6366f1'];
  try { if (localStorage.getItem('wrm.chains.view') === 'list') state.mode = 'list'; } catch (_) { /* storage blocked */ }

  function css() {
    if (state.styleAdded) return;
    state.styleAdded = true;
    const s = document.createElement('style');
    s.textContent = `
.chn-wrap { display:flex; width:min(1320px,97vw); height:88vh; padding:0; overflow:hidden; }
.chn-side { width:250px; min-width:200px; border-right:1px solid var(--border2); display:flex; flex-direction:column; }
.chn-side-head { padding:14px 14px 8px; display:flex; align-items:center; justify-content:space-between; gap:8px; flex-wrap:wrap; }
.chn-side-head h3 { margin:0; font-size:14px; }
.chn-list { flex:1; overflow:auto; padding:6px; }
.chn-item { padding:9px 10px; border-radius:8px; cursor:pointer; border:1px solid transparent; margin-bottom:4px; }
.chn-item:hover, .chn-item.sel { background:var(--bg3); }
.chn-item.sel { border-color:var(--accent); }
.chn-item .n { font-size:12.5px; font-weight:600; }
.chn-item .m { font-size:10.5px; color:var(--text2); margin-top:3px; }
.chn-side-h { font-size:9.5px; text-transform:uppercase; letter-spacing:.06em; color:var(--text3); padding:8px 8px 3px; }
.chn-item.draft { border-style:dashed; border-color:var(--accent); }
.chn-item.draft .n { font-weight:600; }
.chn-dot.draft { background:var(--accent); opacity:1; }
.chn-badge.draft { background:var(--bg3); color:var(--accent); border:1px solid var(--accent); }
.chn-draft-note { font-size:11.5px; line-height:1.45; padding:8px 10px; border:1px solid var(--accent); border-radius:8px; background:var(--bg3); margin-bottom:8px; }
.chn-draft-note b { color:var(--accent); }
.chn-dot { display:inline-block; width:8px; height:8px; border-radius:50%; background:var(--text3); opacity:.6; margin-right:7px; vertical-align:middle; flex:none; }
.chn-dot.ok { background:#3fb950; opacity:1; } .chn-dot.bad { background:#f85149; opacity:1; } .chn-dot.warn { background:#d29922; opacity:1; }
.chn-main { flex:1; display:flex; flex-direction:column; min-width:0; position:relative; }
.chn-head { padding:10px 14px; border-bottom:1px solid var(--border2); display:flex; align-items:center; gap:8px; flex-wrap:wrap; }
.chn-head h3 { margin:0; font-size:14px; flex:1; min-width:140px; }
.chn-seg { display:inline-flex; border:1px solid var(--border2); border-radius:7px; overflow:hidden; }
.chn-wrap .chn-seg button { border:0; border-radius:0; }
.chn-wrap .chn-seg button + button { border-left:1px solid var(--border2); }
.chn-body { flex:1; overflow:auto; padding:12px 14px; min-height:0; }
.chn-body.dg { padding:0; overflow:hidden; display:flex; flex-direction:column; }
.chn-empty { padding:30px 14px; text-align:center; color:var(--text2); font-size:12px; }
.chn-badge { font-size:11px; padding:2px 8px; border-radius:10px; background:var(--bg3); color:var(--text2); }
.chn-badge.active { background:var(--accent); color:#fff; }
.chn-sec { margin:4px 0 16px; }
.chn-sec h4 { margin:0 0 8px; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-nodes { display:grid; grid-template-columns:repeat(auto-fill,minmax(280px,1fr)); gap:10px; }
.chn-node { border:1px solid var(--border2); border-radius:10px; padding:10px 12px; background:var(--bg2); }
.chn-node .h { display:flex; align-items:center; gap:6px; font-weight:600; font-size:12.5px; }
.chn-node .r { font-size:10.5px; color:var(--text2); margin:2px 0 6px; }
.chn-node .e, .chn-det .e { font-size:10.5px; color:#f85149; margin-bottom:6px; overflow-wrap:anywhere; }
.chn-chk { display:flex; flex-wrap:wrap; align-items:baseline; gap:2px 6px; font-size:11.5px; padding:3px 0; border-top:1px dashed var(--border2); }
.chn-chk .l { flex:1; min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .v { font-family:var(--font-mono); font-size:11px; max-width:45%; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .d { flex-basis:100%; font-size:10px; color:var(--text2); padding-left:14px; overflow-wrap:anywhere; }
.chn-chk.click { cursor:pointer; } .chn-chk.click:hover { background:var(--bg3); }
.chn-chk-more { font-size:10.5px; padding:4px 0 6px 14px; color:var(--text2); }
.chn-chk-more pre { font-family:var(--font-mono); font-size:10.5px; white-space:pre-wrap; overflow-wrap:anywhere; background:var(--bg); border-radius:6px; padding:6px 8px; margin:4px 0; color:var(--text); }
.chn-row { display:flex; gap:8px; align-items:center; font-size:11.5px; padding:4px 0; flex-wrap:wrap; }
.chn-row .k { color:var(--text2); font-size:10.5px; }
.chn-tl { font-size:11px; padding:3px 0; display:flex; gap:8px; align-items:baseline; }
.chn-tl .a { color:var(--text2); width:150px; flex:none; font-family:var(--font-mono); font-size:10.5px; }
.chn-dg-top { flex:1; display:flex; min-height:0; }
.chn-canvas { flex:1; position:relative; min-width:0; background:var(--bg); background-image:radial-gradient(var(--border2) 1px, transparent 1px); background-size:22px 22px; }
.chn-canvas svg { width:100%; height:100%; display:block; touch-action:none; cursor:grab; user-select:none; }
.chn-canvas svg.panning { cursor:grabbing; }
.chn-tools { position:absolute; top:8px; left:8px; display:flex; gap:4px; flex-wrap:wrap; }
.chn-tools button { min-width:30px; }
.chn-hint { position:absolute; bottom:6px; left:10px; right:10px; font-size:10.5px; color:var(--text2); pointer-events:none; }
.chn-legend { position:absolute; top:8px; right:10px; font-size:10.5px; color:var(--text2); display:flex; gap:10px; background:var(--bg2); border:1px solid var(--border2); border-radius:7px; padding:3px 8px; }
.chn-legend i { display:inline-block; width:18px; height:0; border-top:2px solid; vertical-align:middle; margin-right:4px; }
.chn-det { width:310px; flex:none; border-left:1px solid var(--border2); overflow:auto; padding:12px; font-size:11.5px; background:var(--bg2); }
.chn-det h4 { margin:0 0 4px; font-size:13px; display:flex; align-items:center; gap:6px; }
.chn-det h5 { margin:14px 0 6px; font-size:10.5px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-det .r { font-size:10.5px; color:var(--text2); margin-bottom:8px; }
.chn-det .act { display:flex; gap:6px; flex-wrap:wrap; margin:8px 0; }
.chn-band { flex:none; max-height:30%; overflow:auto; border-top:1px solid var(--border2); padding:8px 14px 10px; }
.chn-band h4 { margin:0 0 6px; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-bar { display:flex; height:16px; border-radius:5px; overflow:hidden; background:var(--bg3); }
.chn-bar span { min-width:3px; height:100%; }
.chn-marks { position:relative; height:10px; margin-top:2px; }
.chn-marks span { position:absolute; top:1px; width:3px; height:8px; border-radius:1px; transform:translateX(-1px); }
.chn-axis { display:flex; justify-content:space-between; font-size:10px; color:var(--text2); margin:2px 0 4px; }
.chn-keys { display:flex; flex-wrap:wrap; gap:4px 12px; font-size:10.5px; color:var(--text2); margin-bottom:6px; }
.chn-keys b { display:inline-block; width:10px; height:10px; border-radius:2px; vertical-align:middle; margin-right:4px; }
.chn-node-g { cursor:pointer; }
.chn-node-g text { font-family:var(--font-ui); pointer-events:none; }
.chn-link-hit { stroke:transparent; stroke-width:14; fill:none; cursor:pointer; }
.chn-link { fill:none; stroke-width:2.2; }
.chn-link.ok { stroke:#3fb950; } .chn-link.bad { stroke:#f85149; } .chn-link.unknown { stroke:var(--text2); opacity:.7; } .chn-link.idle { stroke:var(--text2); opacity:.45; stroke-dasharray:4 5; }
.chn-link.k-failover:not(.idle), .chn-link.k-dependency:not(.idle) { stroke-dasharray:8 5; }
.chn-link.sel { stroke-width:3.4; }
.chn-flow { fill:none; stroke-width:2.4; stroke:#7ee787; stroke-dasharray:3 15; animation:chn-flow 1.1s linear infinite; pointer-events:none; }
.chn-flow.rev { animation-direction:reverse; }
@keyframes chn-flow { to { stroke-dashoffset:-18; } }
.chn-ll { font-size:10px; fill:var(--text2); paint-order:stroke; stroke:var(--bg); stroke-width:3px; pointer-events:none; }
.chn-editor { position:absolute; inset:0; background:var(--bg2); display:flex; flex-direction:column; z-index:5; }
.chn-editor-head { padding:10px 14px; border-bottom:1px solid var(--border2); display:flex; gap:8px; align-items:center; flex-wrap:wrap; }
.chn-editor-head h3 { margin:0; font-size:13px; flex:1; min-width:160px; }
.chn-editor-body { flex:1; display:flex; min-height:0; }
.chn-editor textarea { flex:1; border:0; border-radius:0; font-family:var(--font-mono); font-size:12px; resize:none; padding:12px 14px; background:var(--bg); color:var(--text); min-width:0; }
.chn-help { width:320px; overflow:auto; border-left:1px solid var(--border2); padding:10px 12px; font-size:11px; color:var(--text2); }
.chn-help h5 { margin:10px 0 4px; font-size:10.5px; text-transform:uppercase; color:var(--text2); }
.chn-help code { font-family:var(--font-mono); font-size:10.5px; }
.chn-errs { max-height:120px; overflow:auto; padding:6px 14px; font-size:11px; color:#f85149; border-top:1px solid var(--border2); }
.chn-errs:empty { display:none; }
.chn-tpls { flex:1; overflow:auto; padding:14px; }
.chn-tpls > h4 { margin:0 0 10px; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-tpl-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(240px,1fr)); gap:10px; }
.chn-tpl { border:1px solid var(--border2); border-radius:10px; padding:12px; cursor:pointer; background:var(--bg); text-align:left; color:var(--text); font:inherit; }
.chn-tpl:hover, .chn-tpl:focus-visible { border-color:var(--accent); }
.chn-tpl b { display:block; font-size:12.5px; margin-bottom:4px; } .chn-tpl span { font-size:11px; color:var(--text2); }
.chn-form { max-width:560px; display:flex; flex-direction:column; gap:10px; font-size:12px; }
.chn-form label { display:flex; flex-direction:column; gap:4px; font-size:11.5px; color:var(--text2); }
.chn-form input, .chn-form select { font-size:12.5px; background:var(--bg); border:1px solid var(--border2); color:var(--text); border-radius:7px; padding:6px 8px; outline:none; font-family:inherit; }
.chn-form input:focus, .chn-form select:focus { border-color:var(--accent); }
.chn-form h4 { font-size:13px; margin:0; } .chn-form h5 { font-size:10.5px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); margin:6px 0 0; }
.chn-form .note { font-size:11px; color:var(--text2); }

/* Guided builder, wizard and dialogs (step 4) */
.chn-build { padding:14px; overflow:auto; }
.chn-build h4 { margin:0 0 4px; font-size:13px; }
.chn-build-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(235px,1fr)); gap:10px; margin-top:12px; }
.chn-build-card { text-align:left; padding:14px; border:1px solid var(--border2); border-radius:10px; background:var(--bg2); color:var(--text); cursor:pointer; display:flex; flex-direction:column; gap:6px; }
.chn-build-card:hover { border-color:var(--accent); }
.chn-build-card .ic { font-size:22px; }
.chn-build-card b { font-size:13px; }
.chn-build-card span { font-size:11px; color:var(--text2); }
.chn-wizard { flex:1; overflow:auto; padding:14px; }
.chn-wizard h4 { margin:0 0 6px; font-size:13px; }
.chn-wizard h5 { margin:12px 0 4px; font-size:11px; color:var(--text2); }
.chn-wizard textarea, .chn-wizard input, .chn-build input, .chn-build textarea { width:100%; }
.chn-wizard textarea { min-height:74px; resize:vertical; }
.chn-picklist { display:flex; flex-wrap:wrap; gap:6px; margin:4px 0 8px; }
.chn-pick { display:inline-flex; align-items:center; gap:6px; border:1px solid var(--border2); border-radius:8px; padding:6px 9px; font-size:11.5px; background:var(--bg2); cursor:pointer; }
.chn-pick input { margin:0; }
.chn-facts { display:flex; flex-direction:column; gap:6px; margin:6px 0 10px; }
.chn-fact { border:1px solid var(--border2); border-radius:8px; padding:7px 9px; background:var(--bg1); font-size:11px; }
.chn-fact b { color:var(--text); }
.chn-fact .sub { color:var(--text3); font-size:10px; }
.chn-fact-row { display:flex; flex-wrap:wrap; gap:4px; align-items:center; margin-top:3px; }
.chn-fact-row .k { color:var(--text3); font-size:10px; }
.chn-guided { flex:1; display:flex; gap:14px; overflow:auto; padding:12px 14px; }
.chn-guided-main { flex:1; min-width:0; }
.chn-guided-help { width:270px; flex:none; font-size:11.5px; color:var(--text2); border-left:1px solid var(--border2); padding-left:12px; }
.chn-guided-help h4 { margin:0 0 6px; font-size:12px; }
.chn-gsec { border:1px solid var(--border2); border-radius:10px; background:var(--bg2); padding:10px 12px; margin-bottom:10px; }
.chn-gsec h4 { margin:0 0 8px; font-size:12.5px; }
.chn-glist { display:flex; flex-direction:column; gap:6px; }
.chn-grow { display:flex; align-items:center; gap:8px; padding:6px 8px; border:1px solid var(--border2); border-radius:8px; background:var(--bg1); }
.chn-grow-main { flex:1; min-width:0; display:flex; flex-direction:column; gap:3px; }
.chn-grow-top { display:flex; gap:6px; align-items:center; }
.chn-grow-top input { flex:1; min-width:0; }
.chn-grow-top input.chn-role { flex:none; width:130px; }
.chn-grow-sub { font-size:10.5px; color:var(--text3); }
.chn-gicon { font-size:15px; width:20px; text-align:center; flex:none; }
.chn-gchips { display:flex; gap:4px; flex-wrap:wrap; padding-left:34px; margin:-2px 0 4px; }
.chn-chip { font-size:9.5px; padding:1px 6px; border-radius:9px; background:var(--bg3); color:var(--text2); }
.chn-gadd { display:flex; flex-wrap:wrap; gap:6px; align-items:center; margin-top:8px; padding-top:8px; border-top:1px dashed var(--border2); }
.chn-gadd select, .chn-gadd input { flex:1; min-width:110px; }
.chn-cgroup { margin-bottom:10px; }
.chn-cgroup h5 { margin:0 0 5px; font-size:11.5px; color:var(--text2); }
.chn-f { display:flex; flex-direction:column; gap:3px; font-size:10.5px; color:var(--text3); margin-bottom:8px; }
.chn-f small { font-size:9.5px; color:var(--text3); }
.chn-grid2 { display:grid; grid-template-columns:repeat(auto-fit,minmax(160px,1fr)); gap:8px 10px; }
.chn-checkline { display:flex; align-items:center; gap:6px; font-size:11.5px; color:var(--text2); margin:4px 0; }
.chn-cond-row { display:flex; flex-wrap:wrap; gap:5px; align-items:center; margin:4px 0; }
.chn-cond-row select, .chn-cond-row input { font-size:11px; padding:3px 5px; }
.chn-cond { font-size:11px; color:var(--text2); padding:3px 0; }
.chn-scn-box { border:1px dashed var(--border2); border-radius:8px; padding:6px 8px; margin-bottom:8px; }
.chn-inline-err { color:#f85149; font-size:11px; flex-basis:100%; }
.chn-dlg-wrap { position:absolute; inset:0; background:rgba(0,0,0,.45); display:flex; align-items:center; justify-content:center; z-index:20; padding:12px; }
.chn-dlg { width:min(470px,96%); max-height:90%; display:flex; flex-direction:column; background:var(--bg2); border:1px solid var(--border2); border-radius:12px; overflow:hidden; }
.chn-dlg-head { display:flex; align-items:center; gap:8px; padding:10px 12px; border-bottom:1px solid var(--border2); }
.chn-dlg-head h4 { margin:0; font-size:13px; flex:1; }
.chn-dlg-body { padding:12px; overflow:auto; }
.chn-svg .chn-port { fill:var(--bg1); stroke:var(--accent); stroke-width:2; opacity:0; cursor:crosshair; }
.chn-svg .chn-connect-line { fill:none; stroke:var(--accent); stroke-width:2; stroke-dasharray:6 4; pointer-events:none; }
.chn-node-g:hover .chn-port, .chn-node-g.sel .chn-port { opacity:1; }
@media (max-width: 900px) { .chn-guided-help { display:none; } }
@media (max-width: 900px) { .chn-det { width:250px; } }
@media (max-width: 760px) {
  .chn-wrap { flex-direction:column; height:92vh; }
  .chn-side { width:auto; max-height:24vh; border-right:0; border-bottom:1px solid var(--border2); }
  .chn-help { display:none; }
  .chn-body.dg { overflow:auto; }
  .chn-dg-top { flex-direction:column; flex:none; }
  .chn-canvas { height:46vh; flex:none; }
  .chn-det { width:auto; border-left:0; border-top:1px solid var(--border2); }
  .chn-band { max-height:none; }
  .chn-legend { display:none; }
}
`;
    document.head.appendChild(s);
  }

  function api(method, path, body) {
    const opts = { method, headers: { 'X-WRM-Request': '1' }, credentials: 'same-origin' };
    if (body !== undefined) { opts.body = JSON.stringify(body); opts.headers['Content-Type'] = 'application/json'; }
    return fetch(path, opts).then(async (res) => {
      const text = await res.text();
      let data = null;
      try { data = text ? JSON.parse(text) : null; } catch (_) { data = null; }
      if (!res.ok) {
        const err = new Error((data && data.error) || t('chn_failed'));
        err.status = res.status; err.errors = (data && data.errors) || [];
        throw err;
      }
      return data;
    });
  }

  function overlay() { return document.getElementById('chn-overlay'); }

  function shell() {
    let ov = overlay();
    if (ov) return ov;
    ov = el('div', 'modal-overlay');
    ov.id = 'chn-overlay';
    ov.innerHTML = `<div class="modal chn-wrap" role="dialog" aria-label="${esc(t('chn_title'))}">
      <div class="chn-side"><div class="chn-side-head"><h3>⛓ ${esc(t('chn_title'))}</h3>
        <span><button class="btn-sec btn-sm" id="chn-import">${esc(t('chn_import'))}</button>
        <button class="btn-sec btn-sm" id="chn-new">＋ ${esc(t('chn_new'))}</button></span></div>
        <div class="chn-list" id="chn-list"></div></div>
      <div class="chn-main"><div class="chn-head"><h3 id="chn-name"></h3>
          <span class="chn-badge" id="chn-scn" style="display:none;"></span>
          <span class="chn-seg" id="chn-mode" role="group" style="display:none;"><button class="btn-sec btn-sm" data-m="diagram">${esc(t('chn_view_diagram'))}</button><button class="btn-sec btn-sm" data-m="list">${esc(t('chn_view_list'))}</button></span>
          <button class="btn-sec btn-sm" id="chn-check" style="display:none;">↻ ${esc(t('chn_check_now'))}</button>
          <button class="btn-sec btn-sm" id="chn-edit" style="display:none;">✎ ${esc(t('chn_edit_chain'))}</button>
          <button class="btn-sec btn-sm" id="chn-json" style="display:none;" title="${esc(t('chn_edit'))}">{ }</button>
          <button class="btn-sec btn-sm" id="chn-export" style="display:none;">⬇ ${esc(t('chn_export'))}</button>
          <button class="btn-sec btn-sm" id="chn-del" style="display:none;">${esc(t('chn_delete'))}</button>
          <button class="btn-pri btn-sm" id="chn-approve" style="display:none;">✓ ${esc(t('chn_approve'))}</button>
          <button class="btn-sec btn-sm" id="chn-discard" style="display:none;">${esc(t('chn_discard'))}</button>
          <button class="btn-sec btn-sm" id="chn-close" aria-label="${esc(t('chn_cancel'))}">✕</button>
        </div>
        <div class="chn-body" id="chn-body"></div></div></div>`;
    document.body.appendChild(ov);
    ov.addEventListener('mousedown', (e) => { if (e.target === ov) close(); });
    ov.querySelector('#chn-close').onclick = close;
    ov.querySelector('#chn-new').onclick = openNew;
    ov.querySelector('#chn-import').onclick = importFile;
    ov.querySelector('#chn-check').onclick = (e) => checkNow(e.currentTarget);
    ov.querySelector('#chn-edit').onclick = editCurrent;
    ov.querySelector('#chn-json').onclick = () => state.current && openEditor(state.current.id, { definition: state.current.definition, enabled: state.current.enabled });
    ov.querySelector('#chn-export').onclick = () => { if (state.current) location.href = `/api/chains/${state.current.id}/export`; };
    ov.querySelector('#chn-del').onclick = deleteChain;
    ov.querySelector('#chn-approve').onclick = approveDraft;
    ov.querySelector('#chn-discard').onclick = discardDraft;
    ov.querySelectorAll('#chn-mode button').forEach((b) => { b.onclick = () => setMode(b.dataset.m); });
    return ov;
  }

  function setMode(m) {
    state.mode = m === 'list' ? 'list' : 'diagram';
    try { localStorage.setItem('wrm.chains.view', state.mode); } catch (_) { /* storage blocked */ }
    renderDetail();
  }

  function open() {
    css();
    shell().classList.add('open');
    Promise.all([
      api('GET', '/api/connections').catch(() => []),
      api('GET', '/api/storages').catch(() => []),
      api('GET', '/api/chains/catalog').catch(() => ({ types: {} })),
    ]).then(([conns, storages, cat]) => { state.conns = conns || []; state.storages = storages || []; state.catalog = cat || { types: {} }; renderSide(); }).catch(() => {});
    loadList();
  }

  function close() {
    const ov = overlay();
    if (ov) ov.classList.remove('open');
  }

  function isOpen() { const ov = overlay(); return !!(ov && ov.classList.contains('open')); }

  function loadList() {
    return Promise.all([
      api('GET', '/api/chains'),
      api('GET', '/api/chains/drafts').catch(() => []),
    ]).then(([list, drafts]) => {
      state.list = list || [];
      state.drafts = drafts || [];
      renderList();
      if (state.current && state.current.draft) {
        if (state.drafts.some((d) => d.id === state.current.draftId)) return selectDraft(state.current.draftId, true);
        state.current = null;
      }
      if (state.current && state.list.some((c) => c.id === state.current.id)) return select(state.current.id, true);
      if (state.list.length) return select(state.list[0].id);
      if (state.drafts.length) return selectDraft(state.drafts[0].id);
      state.current = null;
      renderDetail();
      return null;
    }).catch((e) => showToast(e.message, 'error'));
  }

  function statusClass(c) {
    if (!c.enabled || !c.last_at) return '';
    if (c.nodes_up < c.nodes_total || c.checks_ok < c.checks_total) return c.nodes_up < c.nodes_total ? 'bad' : 'warn';
    return 'ok';
  }

  function renderList() {
    const box = $('#chn-list');
    if (!box) return;
    box.textContent = '';
    if (state.drafts.length) {
      box.appendChild(el('div', 'chn-side-h', t('chn_drafts')));
      state.drafts.forEach((d) => {
        const item = el('div', 'chn-item draft' + (state.current && state.current.draftId === d.id ? ' sel' : ''));
        item.tabIndex = 0;
        const n = el('div', 'n');
        n.appendChild(el('span', 'chn-dot draft'));
        n.appendChild(document.createTextNode(d.name || t('chn_draft_untitled')));
        item.appendChild(n);
        item.appendChild(el('div', 'm', fmt(t('chn_draft_counts'), { n: d.nodes, c: d.checks })));
        item.onclick = () => selectDraft(d.id);
        item.onkeydown = (e) => { if (e.key === 'Enter') selectDraft(d.id); };
        box.appendChild(item);
      });
    }
    if (!state.list.length && !state.drafts.length) { box.appendChild(el('div', 'chn-empty', t('chn_none'))); return; }
    state.list.forEach((c) => {
      const item = el('div', 'chn-item' + (state.current && !state.current.draft && state.current.id === c.id ? ' sel' : ''));
      item.tabIndex = 0;
      const n = el('div', 'n');
      n.appendChild(el('span', 'chn-dot ' + statusClass(c)));
      n.appendChild(document.createTextNode(c.name));
      const m = el('div', 'm');
      m.textContent = !c.enabled ? t('chn_disabled')
        : [c.scenario_label || t('chn_no_scenario'), c.nodes_total ? fmt(t('chn_nodes_count'), { up: c.nodes_up, n: c.nodes_total }) : '',
          fmt(t('chn_ok_count'), { ok: c.checks_ok, n: c.checks_total })].filter(Boolean).join(' · ');
      item.appendChild(n); item.appendChild(m);
      item.onclick = () => select(c.id);
      item.onkeydown = (e) => { if (e.key === 'Enter') select(c.id); };
      box.appendChild(item);
    });
  }

  // selectDraft opens the AI's draft on the same live diagram. It is not monitored and its
  // checks show "not checked yet" until the user approves it.
  function selectDraft(id, keep) {
    return api('GET', `/api/chains/drafts/${id}`).then((d) => {
      if (!keep || !state.current || state.current.draftId !== id) {
        state.sel = null; state.openCheck = null; state.timeline = [];
      }
      state.current = { id: 'draft:' + id, draft: true, draftId: id, definition: d.definition,
        enabled: false, state: {}, scenario: '', last_at: '', created_at: d.created_at };
      state.mode = 'diagram'; // a draft is reviewed on the diagram, whatever view was last used
      state.pos = {};
      renderList();
      renderDetail();
      return d;
    }).catch((e) => showToast(e.message, 'error'));
  }

  // keep: a live refresh of the same chain keeps the selection and the zoom.
  function select(id, keep) {
    return api('GET', `/api/chains/${id}`).then((c) => {
      if (!keep || !state.current || state.current.id !== c.id) { state.sel = null; state.openCheck = null; state.timeline = []; }
      state.current = c;
      state.pos = {};
      renderList();
      renderDetail();
      return api('GET', `/api/chains/${id}/timeline`).then((tl) => { state.timeline = tl || []; renderTimeline(state.timeline); renderSide(); }).catch(() => {});
    }).catch((e) => showToast(e.message, 'error'));
  }

  // ─── status of nodes and links ──────────────────────────────────────────────────────────

  function checkRes(id) { const st = (state.current && state.current.state) || {}; return (st.checks || {})[id]; }

  function nodeStatus(d, n) {
    const ns = ((state.current.state || {}).nodes || {})[n.id];
    if (ns && ns.state === 'down') return 'bad';
    const res = (d.checks || []).filter((x) => x.node === n.id).map((x) => checkRes(x.id)).filter(Boolean);
    if (res.length) {
      const ok = res.filter((r) => r.ok).length;
      if (ok === res.length) return 'ok';
      return n.kind !== 'server' && ok === 0 ? 'bad' : 'warn';
    }
    return ns ? 'ok' : '';
  }

  function linkStatus(d, l) {
    const c = state.current;
    if (!c.enabled || !c.last_at) return 'unknown';
    if ((l.scenarios || []).length && l.scenarios.indexOf(c.scenario) < 0) return 'idle';
    if (l.check) { const r = checkRes(l.check); return r ? (r.ok ? 'ok' : 'bad') : 'unknown'; }
    const a = nodeStatus(d, node(d, l.from)), b = nodeStatus(d, node(d, l.to));
    if (a === 'bad' || b === 'bad') return 'bad';
    return a && b ? 'ok' : 'unknown';
  }

  function node(d, id) { return (d.nodes || []).find((x) => x.id === id) || { id, kind: 'external' }; }

  // ─── layout ─────────────────────────────────────────────────────────────────────────────

  // autoLayout places nodes in columns by their longest path along the links (sources left),
  // and orders each column by the average row of its predecessors to keep crossings low.
  function autoLayout(d) {
    const ids = (d.nodes || []).map((n) => n.id);
    const layer = {};
    ids.forEach((id) => { layer[id] = 0; });
    const links = (d.links || []).filter((l) => l.from in layer && l.to in layer && l.from !== l.to);
    for (let i = 0; i < ids.length; i++) {
      let moved = false;
      links.forEach((l) => {
        if (layer[l.to] < layer[l.from] + 1 && layer[l.from] + 1 < ids.length) { layer[l.to] = layer[l.from] + 1; moved = true; }
      });
      if (!moved) break;
    }
    const cols = [];
    ids.forEach((id) => { (cols[layer[id]] = cols[layer[id]] || []).push(id); });
    const row = {};
    cols.forEach((col, li) => {
      if (li > 0) {
        const bary = (id) => {
          const ps = links.filter((l) => l.to === id && layer[l.from] < li).map((l) => row[l.from]);
          return ps.length ? ps.reduce((a, b) => a + b, 0) / ps.length : 1e6;
        };
        const order = col.map((id, i) => ({ id, i, b: bary(id) }));
        order.sort((a, b) => (a.b - b.b) || (a.i - b.i));
        cols[li] = order.map((o) => o.id);
      }
      cols[li].forEach((id, i) => { row[id] = i; });
    });
    const tall = Math.max(1, ...cols.map((c) => (c ? c.length : 0)));
    const out = {};
    cols.forEach((col, li) => {
      if (!col) return;
      const off = (tall - col.length) * GY / 2;
      col.forEach((id, i) => { out[id] = { x: 40 + li * GX, y: 40 + off + i * GY }; });
    });
    return out;
  }

  function positions(d) {
    const auto = autoLayout(d);
    const out = {};
    (d.nodes || []).forEach((n) => {
      out[n.id] = state.pos[n.id] || ((n.x || n.y) ? { x: n.x, y: n.y } : auto[n.id]);
    });
    return out;
  }

  // ─── the diagram ────────────────────────────────────────────────────────────────────────

  function anchor(p, q) {
    const cx = p.x + NW / 2, cy = p.y + NH / 2, dx = q.x - p.x, dy = q.y - p.y;
    if (Math.abs(dx) * NH >= Math.abs(dy) * NW) return { x: dx >= 0 ? p.x + NW : p.x, y: cy, nx: dx >= 0 ? 1 : -1, ny: 0 };
    return { x: cx, y: dy >= 0 ? p.y + NH : p.y, nx: 0, ny: dy >= 0 ? 1 : -1 };
  }

  function linkPath(p, q) {
    const a = anchor(p, q), b = anchor(q, p);
    const k = Math.max(40, Math.hypot(b.x - a.x, b.y - a.y) / 2.5);
    return {
      d: `M${a.x},${a.y} C${a.x + a.nx * k},${a.y + a.ny * k} ${b.x + b.nx * k},${b.y + b.ny * k} ${b.x},${b.y}`,
      mx: (a.x + b.x) / 2 + (a.nx + b.nx) * k * 0.375, my: (a.y + b.y) / 2 + (a.ny + b.ny) * k * 0.375,
    };
  }

  const ARROW = { ok: '#3fb950', bad: '#f85149', unknown: '#8892a0', idle: '#8892a0' };

  function nodeShape(n, cls, selected) {
    const stroke = { ok: '#3fb950', warn: '#d29922', bad: '#f85149' }[cls] || 'var(--border2)';
    const sw = selected ? 3 : 1.6;
    const sel = selected ? ' stroke-dasharray="none"' : '';
    if (n.kind === 'storage') {
      return `<path d="M0,10 a${NW / 2},10 0 0,0 ${NW},0 v${NH - 20} a${NW / 2},10 0 0,1 ${-NW},0 z" fill="var(--bg2)" stroke="${stroke}" stroke-width="${sw}"${sel}/>
        <path d="M0,10 a${NW / 2},10 0 0,0 ${NW},0 a${NW / 2},10 0 0,0 ${-NW},0 z" fill="var(--bg3)" stroke="${stroke}" stroke-width="${sw}"/>`;
    }
    if (n.kind === 'external') {
      return `<rect width="${NW}" height="${NH}" rx="${NH / 2}" fill="var(--bg2)" stroke="${stroke}" stroke-width="${sw}"${selected ? '' : ' stroke-dasharray="6 4"'}/>`;
    }
    return `<rect width="${NW}" height="${NH}" rx="10" fill="var(--bg2)" stroke="${stroke}" stroke-width="${sw}"/>`;
  }

  function renderDiagram(wrap) {
    const c = state.current, d = c.definition;
    const pos = positions(d);
    const sel = state.sel || {};
    let defs = '<defs>';
    Object.keys(ARROW).forEach((k) => {
      defs += `<marker id="chn-a-${k}" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0,0 L10,5 L0,10 z" fill="${ARROW[k]}"/></marker>`;
    });
    defs += '</defs>';
    let links = '', flows = '', labels = '';
    (d.links || []).forEach((l, i) => {
      const p = pos[l.from], q = pos[l.to];
      if (!p || !q || l.from === l.to) return;
      const st = linkStatus(d, l);
      const g = linkPath(p, q);
      const kind = l.kind || 'data';
      const isSel = sel.link === i;
      const rev = kind === 'read';
      const mark = kind === 'dependency' ? '' : (rev ? ` marker-start="url(#chn-a-${st})"` : ` marker-end="url(#chn-a-${st})"`);
      links += `<path class="chn-link ${st} k-${kind}${isSel ? ' sel' : ''}" d="${g.d}"${mark}/>`;
      links += `<path class="chn-link-hit" data-link="${i}" d="${g.d}"><title>${esc(label(d, l.from) + ' → ' + label(d, l.to))}</title></path>`;
      if (st === 'ok' && kind !== 'dependency') flows += `<path class="chn-flow${rev ? ' rev' : ''}" d="${g.d}"/>`;
      const txt = l.label || (kind !== 'data' ? t('chn_lk_' + kind) : '');
      if (txt) labels += `<text class="chn-ll" x="${g.mx}" y="${g.my - 4}" text-anchor="middle">${esc(trunc(txt, 28))}</text>`;
    });
    let nodes = '';
    (d.nodes || []).forEach((n) => {
      const p = pos[n.id];
      const cls = nodeStatus(d, n);
      const checks = (d.checks || []).filter((x) => x.node === n.id);
      const res = checks.map((x) => checkRes(x.id));
      const okN = res.filter((r) => r && r.ok).length;
      const icon = { server: '🖥', storage: '🪣', external: '🌐' }[n.kind] || '•';
      const ns = ((c.state || {}).nodes || {})[n.id];
      const sub = [n.role || t('chn_kind_' + n.kind), ns && ns.state === 'down' ? t('chn_down') : ''].filter(Boolean).join(' · ');
      const top = n.kind === 'storage' ? 30 : 24;
      let dots = '';
      res.slice(0, 12).forEach((r, i) => {
        const col = r ? (r.ok ? '#3fb950' : '#f85149') : '#8892a0';
        dots += `<circle cx="${18 + i * 11}" cy="${NH - 13}" r="3.6" fill="${col}"/>`;
      });
      const count = checks.length ? `<text x="${NW - 12}" y="${NH - 9}" text-anchor="end" font-size="10.5" fill="var(--text2)">${okN}/${checks.length}</text>` : '';
      nodes += `<g class="chn-node-g" data-node="${esc(n.id)}" transform="translate(${p.x},${p.y})" tabindex="0" role="button" aria-label="${esc(n.label || n.id)}">
        ${nodeShape(n, cls, sel.node === n.id)}
        <text x="12" y="${top}" font-size="13" fill="var(--text)" font-weight="600">${icon} ${esc(trunc(n.label || n.id, 20))}</text>
        <text x="12" y="${top + 16}" font-size="10.5" fill="var(--text2)">${esc(trunc(sub, 30))}</text>
        ${dots}${count}<circle class="chn-port" data-port="${esc(n.id)}" cx="${NW}" cy="${NH / 2}" r="7"><title>${esc(t('chn_connect'))}</title></circle><title>${esc([n.label || n.id, n.role, n.id].filter(Boolean).join(' · '))}</title></g>`;
    });
    let svg = wrap.querySelector('svg');
    if (!svg) {
      svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
      svg.setAttribute('preserveAspectRatio', 'xMidYMid meet');
      svg.setAttribute('role', 'img');
      wrap.insertBefore(svg, wrap.firstChild);
      wireSvg(svg);
    }
    svg.setAttribute('aria-label', d.name);
    svg.innerHTML = defs + `<g class="chn-links">${links}</g><g>${flows}</g><g>${labels}</g><g class="chn-nodes-g">${nodes}</g>`;
    if (!state.views[c.id]) state.views[c.id] = fitView(pos);
    applyView(svg);
  }

  function fitView(pos) {
    const pts = Object.values(pos).filter(Boolean);
    if (!pts.length) return { x: 0, y: 0, w: 800, h: 400 };
    const minX = Math.min(...pts.map((p) => p.x)) - 40, minY = Math.min(...pts.map((p) => p.y)) - 50;
    const maxX = Math.max(...pts.map((p) => p.x + NW)) + 40, maxY = Math.max(...pts.map((p) => p.y + NH)) + 40;
    return { x: minX, y: minY, w: Math.max(420, maxX - minX), h: Math.max(240, maxY - minY) };
  }

  function applyView(svg) {
    const v = state.views[state.current.id];
    svg.setAttribute('viewBox', `${v.x} ${v.y} ${v.w} ${v.h}`);
  }

  function zoom(factor, cx, cy) {
    const c = state.current;
    const svg = $('.chn-canvas svg');
    if (!c || !svg) return;
    const v = state.views[c.id];
    const nw = Math.min(10000, Math.max(200, v.w * factor));
    const f = nw / v.w;
    const px = cx == null ? v.x + v.w / 2 : cx, py = cy == null ? v.y + v.h / 2 : cy;
    state.views[c.id] = { x: px - (px - v.x) * f, y: py - (py - v.y) * f, w: v.w * f, h: v.h * f };
    applyView(svg);
  }

  function svgPoint(svg, e) {
    const m = svg.getScreenCTM();
    if (!m) return { x: 0, y: 0 };
    const p = new DOMPoint(e.clientX, e.clientY).matrixTransform(m.inverse());
    return { x: p.x, y: p.y };
  }

  // Dragging a node moves it (and saves every position, so the picture stays as seen); a
  // press that does not move selects. Dragging the background pans; the wheel zooms.
  function wireSvg(svg) {
    let drag = null;
    svg.addEventListener('pointerdown', (e) => {
      if (e.button !== 0 || !state.current) return;
      const g = e.target.closest('[data-node]');
      const lk = e.target.closest('[data-link]');
      const port = e.target.closest('[data-port]');
      const p = svgPoint(svg, e);
      const d = state.current.definition;
      if (port) {
        drag = { kind: 'connect', from: port.getAttribute('data-port'), p };
        svg.setPointerCapture(e.pointerId);
        return;
      }
      if (g) {
        const id = g.getAttribute('data-node');
        const start = positions(d)[id];
        drag = { kind: 'node', id, p, start, moved: false };
      } else if (lk) {
        drag = { kind: 'link', i: Number(lk.getAttribute('data-link')), p, moved: false };
      } else {
        drag = { kind: 'pan', cx: e.clientX, cy: e.clientY, v: Object.assign({}, state.views[state.current.id]), moved: false };
        svg.classList.add('panning');
      }
      svg.setPointerCapture(e.pointerId);
    });
    svg.addEventListener('pointermove', (e) => {
      if (!drag || !state.current) return;
      if (drag.kind === 'pan') {
        const r = svg.getBoundingClientRect();
        const v = drag.v;
        const s = Math.max(v.w / r.width, v.h / r.height);
        if (Math.abs(e.clientX - drag.cx) + Math.abs(e.clientY - drag.cy) > 3) drag.moved = true;
        state.views[state.current.id] = { x: v.x - (e.clientX - drag.cx) * s, y: v.y - (e.clientY - drag.cy) * s, w: v.w, h: v.h };
        applyView(svg);
        return;
      }
      if (drag.kind === 'connect') {
        state.connectPoint = svgPoint(svg, e);
        drawConnectLine(svg, drag.from, state.connectPoint);
        return;
      }
      if (drag.kind !== 'node') return;
      const p = svgPoint(svg, e);
      const dx = p.x - drag.p.x, dy = p.y - drag.p.y;
      if (!drag.moved && Math.hypot(dx, dy) < 4) return;
      if (!drag.moved) {
        drag.moved = true;
        // freeze the current picture: every node keeps where it is now
        state.pos = positions(state.current.definition);
      }
      state.pos[drag.id] = { x: Math.max(1, Math.round(drag.start.x + dx)), y: Math.max(1, Math.round(drag.start.y + dy)) };
      renderDiagram(svg.parentNode);
    });
    const end = (e) => {
      if (!drag) return;
      const was = drag;
      drag = null;
      svg.classList.remove('panning');
      try { svg.releasePointerCapture(e.pointerId); } catch (_) { /* already released */ }
      if (was.kind === 'connect') {
        const line = svg.querySelector('.chn-connect-line');
        if (line) line.remove();
        if (e.type === 'pointercancel') return;
        const el2 = document.elementFromPoint(e.clientX, e.clientY);
        const target = el2 && el2.closest ? el2.closest('[data-node]') : null;
        const to = target && target.getAttribute('data-node');
        if (to && to !== was.from) {
          openLinkDialog(state.current.definition, { from: was.from, to }, (l) => applyDefinition((d) => { d.links.push(l); }));
        }
        return;
      }
      if (e.type === 'pointercancel') return;
      if (was.kind === 'node' && was.moved) { saveLayout(); return; }
      if (was.kind === 'node' && state.pickTarget) {
        const from = state.pickTarget;
        state.pickTarget = null;
        if (was.id !== from) openLinkDialog(state.current.definition, { from, to: was.id }, (l) => applyDefinition((d) => { d.links.push(l); }));
        renderDiagram(svg.parentNode);
        renderSide();
        return;
      }
      if (was.kind === 'node') state.sel = { node: was.id };
      else if (was.kind === 'link') state.sel = { link: was.i };
      else if (!was.moved) state.sel = null;
      else return;
      state.openCheck = null;
      renderDiagram(svg.parentNode);
      renderSide();
    };
    svg.addEventListener('pointerup', end);
    svg.addEventListener('pointercancel', end);
    svg.addEventListener('wheel', (e) => {
      e.preventDefault();
      const p = svgPoint(svg, e);
      zoom(e.deltaY > 0 ? 1.12 : 1 / 1.12, p.x, p.y);
    }, { passive: false });
    svg.addEventListener('keydown', (e) => {
      const g = e.target.closest && e.target.closest('[data-node]');
      if (g && (e.key === 'Enter' || e.key === ' ')) {
        e.preventDefault();
        state.sel = { node: g.getAttribute('data-node') };
        state.openCheck = null;
        renderDiagram(svg.parentNode);
        renderSide();
      }
    });
  }

  function drawConnectLine(svg, from, to) {
    let line = svg.querySelector('.chn-connect-line');
    if (!line) {
      line = document.createElementNS('http://www.w3.org/2000/svg', 'path');
      line.setAttribute('class', 'chn-connect-line');
      svg.appendChild(line);
    }
    const p = positions(state.current.definition)[from];
    if (!p) return;
    line.setAttribute('d', `M${p.x + NW},${p.y + NH / 2} L${to.x},${to.y}`);
  }

  function saveLayout() {
    const c = state.current;
    if (!c) return;
    const positionsNow = Object.assign({}, state.pos);
    api('PUT', `/api/chains/${c.id}/layout`, { positions: positionsNow }).then((got) => {
      if (state.current && state.current.id === got.id) {
        state.current.definition = got.definition;
      }
    }).catch((e) => showToast(t('chn_pos_failed') + ' ' + e.message, 'error'));
  }

  function resetLayout() {
    const c = state.current;
    if (!c || !confirm(t('chn_auto_layout_q'))) return;
    api('PUT', `/api/chains/${c.id}/layout`, { reset: true }).then((got) => {
      state.current.definition = got.definition;
      state.pos = {};
      delete state.views[c.id];
      renderDetail();
    }).catch((e) => showToast(e.message, 'error'));
  }

  // ─── the detail panel next to the diagram ──────────────────────────────────────────────

  function renderSide() {
    const box = $('#chn-det');
    const c = state.current;
    if (!box || !c) return;
    const d = c.definition;
    box.textContent = '';
    const sel = state.sel || {};
    if (sel.node && (d.nodes || []).some((n) => n.id === sel.node)) { sideNode(box, d, node(d, sel.node)); return; }
    if (sel.link != null && (d.links || [])[sel.link]) { sideLink(box, d, d.links[sel.link]); return; }
    sideSummary(box, d);
  }

  function sideSummary(box, d) {
    const c = state.current;
    const h = el('h4');
    const sum = state.list.find((x) => x.id === c.id) || {};
    h.appendChild(el('span', 'chn-dot ' + statusClass(sum)));
    h.appendChild(document.createTextNode(t('chn_summary')));
    box.appendChild(h);
    box.appendChild(el('div', 'r', t('chn_last') + ': ' + when(c.last_at)));
    if (!c.enabled) box.appendChild(el('div', 'r', t('chn_disabled')));
    if (d.description) box.appendChild(el('div', '', d.description));
    if (sum.checks_total != null) {
      const row = el('div', 'chn-row');
      if (sum.nodes_total) row.appendChild(el('span', 'chn-badge', fmt(t('chn_nodes_count'), { up: sum.nodes_up, n: sum.nodes_total })));
      row.appendChild(el('span', 'chn-badge', fmt(t('chn_ok_count'), { ok: sum.checks_ok, n: sum.checks_total })));
      box.appendChild(row);
    }
    if ((d.scenarios || []).length) {
      box.appendChild(el('h5', '', t('chn_scenarios')));
      d.scenarios.forEach((s, i) => {
        const row = el('div', 'chn-row');
        const b = el('span', 'chn-badge' + (s.id === c.scenario ? ' active' : ''), s.label || s.id);
        if (s.id !== c.scenario) b.style.borderLeft = '3px solid ' + SCN_COLORS[i % SCN_COLORS.length];
        row.appendChild(b);
        ruleText(row, s.when, d);
        box.appendChild(row);
      });
    }
    box.appendChild(el('h5', '', t('chn_nodes')));
    (d.nodes || []).forEach((n) => {
      const row = el('div', 'chn-chk click');
      row.appendChild(el('span', 'chn-dot ' + nodeStatus(d, n)));
      row.appendChild(el('span', 'l', n.label || n.id));
      row.appendChild(el('span', 'v', n.role || t('chn_kind_' + n.kind)));
      row.onclick = () => { state.sel = { node: n.id }; renderDiagramOnly(); renderSide(); };
      box.appendChild(row);
    });
  }

  function ruleText(parent, when, def) {
    const k = el('span', 'k', ruleSummary(def || state.current.definition, when));
    k.style.overflowWrap = 'anywhere';
    parent.appendChild(k);
  }

  // ruleSummary turns a scenario rule into a sentence a person can check at a glance.
  function ruleSummary(def, when) {
    if (!when) return t('chn_g_cond_none');
    const one = (c) => {
      if (c.all) return c.all.map(one).filter(Boolean).join(' · ');
      if (c.any) return c.any.map(one).filter(Boolean).join(' | ');
      if (c.node) {
        const val = c.value === 'up' ? t('chn_g_up') : c.value === 'down' ? t('chn_g_down') : c.value;
        return nodeName(def, c.node) + ' ' + (c.op || '==') + ' ' + val;
      }
      if (c.check) {
        const subj = checkName(def, c.check);
        if (c.field === 'ok') {
          const yes = c.value === 'true' || c.value === '';
          const op = c.op || '==';
          const val = yes ? t('chn_g_ok') : t('chn_check_bad');
          return subj + ' ' + op + ' ' + val;
        }
        return subj + ' ' + (c.op || '==') + ' ' + (c.value || '');
      }
      return '';
    };
    const parts = (when.all || when.any || [when]).map(one).filter(Boolean);
    return parts.length ? parts.join(' · ') : t('chn_g_cond_none');
  }

  function renderDiagramOnly() {
    const wrap = $('.chn-canvas');
    if (wrap && state.current) renderDiagram(wrap);
  }

  function sideNode(box, d, n) {
    const c = state.current;
    const ns = ((c.state || {}).nodes || {})[n.id];
    const h = el('h4');
    h.appendChild(el('span', 'chn-dot ' + nodeStatus(d, n)));
    h.appendChild(document.createTextNode(n.label || n.id));
    box.appendChild(h);
    box.appendChild(el('div', 'r', [t('chn_kind_' + n.kind), n.role, n.id].filter(Boolean).join(' · ')
      + (ns ? ' · ' + t(ns.state === 'up' ? 'chn_up' : 'chn_down') : '')));
    if (n.kind === 'server' && n.conn_id) {
      const conn = (state.conns || []).find((x) => x.id === n.conn_id);
      box.appendChild(el('div', 'r', t('chn_conn') + ': ' + (conn ? conn.name : '#' + n.conn_id)));
    }
    if (n.kind === 'storage' && n.storage_id) box.appendChild(el('div', 'r', t('chn_storage') + ': #' + n.storage_id));
    if (ns && ns.error) box.appendChild(el('div', 'e', ns.error));
    const act = el('div', 'act');
    if (!c.draft) {
      const b = el('button', 'btn-sec btn-sm', '↻ ' + t('chn_check_now'));
      b.onclick = () => checkNow(b);
      act.appendChild(b);
    }
    const edit = el('button', 'btn-sec btn-sm', '✎ ' + t('chn_edit_chain'));
    edit.onclick = editCurrent;
    act.appendChild(edit);
    const addChk = el('button', 'btn-sec btn-sm', '＋ ' + t('chn_add_check'));
    addChk.onclick = () => openCheckDialog(d, n, null, (chk) => applyDefinition((dd) => { dd.checks.push(chk); }));
    act.appendChild(addChk);
    const conn = el('button', 'btn-sec btn-sm', '⇢ ' + t('chn_connect'));
    conn.onclick = () => { state.pickTarget = n.id; showToast(t('chn_pick_target')); };
    act.appendChild(conn);
    const delN = el('button', 'btn-sec btn-sm', '✕');
    delN.title = t('chn_delete_node');
    delN.onclick = () => {
      if (!confirm(t('chn_delete_node') + ' · ' + (n.label || n.id) + '?')) return;
      applyDefinition((dd) => stripNode(dd, n.id));
    };
    act.appendChild(delN);
    box.appendChild(act);
    const checks = (d.checks || []).filter((x) => x.node === n.id);
    box.appendChild(el('h5', '', t('chn_checks')));
    checks.forEach((x) => box.appendChild(checkRow(x, true)));
    const links = (d.links || []).map((l, i) => ({ l, i })).filter((o) => o.l.from === n.id || o.l.to === n.id);
    if (links.length) {
      box.appendChild(el('h5', '', t('chn_links')));
      links.forEach(({ l, i }) => {
        const row = el('div', 'chn-chk click');
        row.appendChild(el('span', 'chn-dot ' + dotOf(linkStatus(d, l))));
        row.appendChild(el('span', 'l', `${label(d, l.from)} → ${label(d, l.to)}`));
        row.appendChild(el('span', 'v', t('chn_lk_' + (l.kind || 'data'))));
        row.onclick = () => { state.sel = { link: i }; renderDiagramOnly(); renderSide(); };
        box.appendChild(row);
      });
    }
    const ids = new Set(checks.map((x) => x.id));
    history(box, state.timeline.filter((e) => (e.kind === 'node' && e.subject === n.id) || (e.kind === 'check' && ids.has(e.subject))));
  }

  function stripNode(def, id) {
    const gone = (def.checks || []).filter((x) => x.node === id).map((x) => x.id);
    def.nodes = (def.nodes || []).filter((n) => n.id !== id);
    def.checks = (def.checks || []).filter((x) => x.node !== id);
    def.links = (def.links || []).filter((l) => l.from !== id && l.to !== id);
    (def.scenarios || []).forEach((s) => cleanConds(s.when, id, gone));
    (def.links || []).forEach((l) => { if (gone.includes(l.check)) delete l.check; });
  }

  function cleanConds(c, nodeID, checkIDs) {
    if (!c) return;
    if (c.node === nodeID) { delete c.node; delete c.op; delete c.value; }
    if (c.check && checkIDs.includes(c.check)) { delete c.check; delete c.field; delete c.op; delete c.value; }
    (c.all || []).forEach((x) => cleanConds(x, nodeID, checkIDs));
    (c.any || []).forEach((x) => cleanConds(x, nodeID, checkIDs));
    const empty = (x) => !x.check && !x.node && !(x.all && x.all.length) && !(x.any && x.any.length);
    if (c.all) c.all = c.all.filter((x) => !empty(x));
    if (c.any) c.any = c.any.filter((x) => !empty(x));
  }

  function dotOf(st) { return { ok: 'ok', bad: 'bad', warn: 'warn' }[st] || ''; }

  function checkRow(x, expandable) {
    const wrap = el('div');
    const r = checkRes(x.id);
    const row = el('div', 'chn-chk' + (expandable ? ' click' : ''));
    row.appendChild(el('span', 'chn-dot ' + (r ? (r.ok ? 'ok' : 'bad') : '')));
    const l = el('span', 'l', (x.label || x.id) + ' · ' + x.type);
    l.title = x.target || '';
    row.appendChild(l);
    row.appendChild(el('span', 'v', r ? r.value : '—'));
    if (r && (r.detail || r.error)) row.appendChild(el('span', 'd', r.error || r.detail));
    wrap.appendChild(row);
    if (expandable) {
      row.tabIndex = 0;
      row.setAttribute('role', 'button');
      row.setAttribute('aria-expanded', state.openCheck === x.id ? 'true' : 'false');
      const toggle = () => { state.openCheck = state.openCheck === x.id ? null : x.id; renderSide(); };
      row.onclick = toggle;
      row.onkeydown = (e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); toggle(); } };
      if (state.openCheck === x.id) {
        const more = el('div', 'chn-chk-more');
        more.appendChild(el('div', '', t('chn_definition')));
        const def = Object.assign({}, x);
        delete def.node;
        more.appendChild(el('pre', '', JSON.stringify(def, null, 2)));
        if (r && r.at) more.appendChild(el('div', '', t('chn_last') + ': ' + when(r.at)));
        const del = el('button', 'btn-sec btn-sm', '✕ ' + t('chn_g_remove'));
        del.onclick = () => applyDefinition((dd) => {
          dd.checks = (dd.checks || []).filter((y) => y.id !== x.id);
          (dd.links || []).forEach((l) => { if (l.check === x.id) delete l.check; });
        });
        more.appendChild(del);
        history(more, state.timeline.filter((e) => e.kind === 'check' && e.subject === x.id), true);
        wrap.appendChild(more);
      }
    }
    return wrap;
  }

  function history(box, list, small) {
    box.appendChild(el(small ? 'div' : 'h5', '', t('chn_history')));
    if (!list.length) { box.appendChild(el('div', 'k', t('chn_no_history'))); return; }
    list.slice(0, 30).forEach((e) => box.appendChild(timelineRow(e)));
  }

  function sideLink(box, d, l) {
    const st = linkStatus(d, l);
    const h = el('h4');
    h.appendChild(el('span', 'chn-dot ' + dotOf(st)));
    h.appendChild(document.createTextNode(`${label(d, l.from)} → ${label(d, l.to)}`));
    box.appendChild(h);
    box.appendChild(el('div', 'r', [t('chn_link'), t('chn_lk_' + (l.kind || 'data')), l.label].filter(Boolean).join(' · ')
      + (st === 'idle' ? ' · ' + t('chn_link_idle') : '')));
    if ((l.scenarios || []).length) {
      const row = el('div', 'chn-row');
      row.appendChild(el('span', 'k', t('chn_link_active_in') + ':'));
      l.scenarios.forEach((id) => {
        const s = (d.scenarios || []).find((x) => x.id === id);
        row.appendChild(el('span', 'chn-badge' + (id === state.current.scenario ? ' active' : ''), s ? (s.label || s.id) : id));
      });
      box.appendChild(row);
    }
    if (l.check) {
      const x = (d.checks || []).find((c) => c.id === l.check);
      if (x) {
        box.appendChild(el('h5', '', t('chn_link_check')));
        box.appendChild(checkRow(x, true));
      }
    }
    const act = el('div', 'act');
    [l.from, l.to].forEach((id) => {
      const b = el('button', 'btn-sec btn-sm', label(d, id));
      b.onclick = () => { state.sel = { node: id }; renderDiagramOnly(); renderSide(); };
      act.appendChild(b);
    });
    const del = el('button', 'btn-sec btn-sm', '✕ ' + t('chn_delete_link'));
    del.onclick = () => applyDefinition((dd) => {
      const at = (dd.links || []).indexOf(l);
      if (at >= 0) dd.links.splice(at, 1);
    });
    act.appendChild(del);
    box.appendChild(act);
  }

  // ─── the main area ──────────────────────────────────────────────────────────────────────

  function renderDetail() {
    const c = state.current;
    const isDraft = !!(c && c.draft);
    const show = (id, on) => { const b = $(id); if (b) b.style.display = on ? '' : 'none'; };
    ['#chn-check', '#chn-export', '#chn-del'].forEach((id) => show(id, !!c && !isDraft));
    show('#chn-edit', !!c);
    show('#chn-json', !!c && !isDraft);
    show('#chn-approve', isDraft);
    show('#chn-discard', isDraft);
    ['#chn-scn', '#chn-mode'].forEach((id) => show(id, !!c));
    document.querySelectorAll('#chn-mode button').forEach((b) => { const on = b.dataset.m === state.mode; b.className = (on ? 'btn-pri' : 'btn-sec') + ' btn-sm'; b.setAttribute('aria-pressed', on ? 'true' : 'false'); });
    const body = $('#chn-body');
    if (!body) return;
    body.textContent = '';
    body.classList.toggle('dg', !!c && state.mode === 'diagram');
    $('#chn-name').textContent = c ? c.definition.name : t('chn_title');
    if (!c) { body.appendChild(el('div', 'chn-empty', state.list.length || state.drafts.length ? t('chn_pick') : t('chn_none'))); return; }
    const eb = $('#chn-edit');
    if (eb) eb.textContent = '✎ ' + (isDraft ? t('chn_g_edit_draft') : t('chn_edit_chain'));
    const scn = $('#chn-scn');
    if (isDraft) {
      scn.textContent = t('chn_draft_badge');
      scn.className = 'chn-badge draft';
      const note = el('div', 'chn-draft-note');
      note.appendChild(el('b', '', t('chn_draft_note_t')));
      note.appendChild(document.createTextNode(' ' + t('chn_draft_note')));
      body.appendChild(note);
    } else {
      scn.textContent = t('chn_scenario') + ': ' + (c.scenario ? c.scenario_label : t('chn_no_scenario'));
      scn.className = 'chn-badge' + (c.scenario ? ' active' : '');
    }
    if (state.mode === 'diagram') renderDiagramView(body); else renderListView(body);
    if (!isDraft) renderTimeline(state.timeline);
  }

  function renderDiagramView(body) {
    const top = el('div', 'chn-dg-top');
    const canvas = el('div', 'chn-canvas');
    const tools = el('div', 'chn-tools');
    const tool = (txt, title, fn) => { const b = el('button', 'btn-sec btn-sm', txt); b.title = title; b.setAttribute('aria-label', title); b.onclick = fn; tools.appendChild(b); };
    tool('−', t('chn_zoom_out'), () => zoom(1.25));
    tool('+', t('chn_zoom_in'), () => zoom(1 / 1.25));
    tool('⤢', t('chn_fit'), () => { delete state.views[state.current.id]; renderDiagramOnly(); });
    tool('⟲', t('chn_auto_layout'), resetLayout);
    tool('＋', t('chn_add_node'), () => openNodeDialog(state.current.definition, (n) => applyDefinition((d) => { d.nodes.push(n); })));
    tool('⇢', t('chn_g_add_link'), () => openLinkDialog(state.current.definition, null, (l) => applyDefinition((d) => { d.links.push(l); })));
    canvas.appendChild(tools);
    const legend = el('div', 'chn-legend');
    [['#3fb950', t('chn_flow')], ['#f85149', t('chn_legend_bad')], ['#8892a0', t('chn_legend_idle')]].forEach(([col, txt]) => {
      const s = el('span'); const i = el('i'); i.style.borderColor = col; s.appendChild(i); s.appendChild(document.createTextNode(txt)); legend.appendChild(s);
    });
    canvas.appendChild(legend);
    canvas.appendChild(el('div', 'chn-hint', t('chn_drag_hint')));
    top.appendChild(canvas);
    const det = el('div', 'chn-det');
    det.id = 'chn-det';
    top.appendChild(det);
    body.appendChild(top);
    const band = el('div', 'chn-band');
    band.id = 'chn-band';
    body.appendChild(band);
    renderDiagram(canvas);
    renderSide();
    const tl = el('div', 'chn-sec');
    tl.id = 'chn-tl';
    band.appendChild(tl);
  }

  function renderListView(body) {
    const c = state.current, d = c.definition;
    const meta = el('div', 'chn-row');
    meta.appendChild(el('span', 'k', t('chn_last') + ': ' + when(c.last_at)));
    if (!c.enabled) meta.appendChild(el('span', 'chn-badge', t('chn_disabled')));
    if (d.description) meta.appendChild(el('span', '', d.description));
    body.appendChild(meta);
    const nodesSec = section(body, t('chn_nodes'));
    const grid = el('div', 'chn-nodes');
    nodesSec.appendChild(grid);
    (d.nodes || []).forEach((n) => {
      const ns = ((c.state || {}).nodes || {})[n.id];
      const card = el('div', 'chn-node');
      const h = el('div', 'h');
      h.appendChild(el('span', 'chn-dot ' + nodeStatus(d, n)));
      h.appendChild(document.createTextNode(n.label || n.id));
      card.appendChild(h);
      card.appendChild(el('div', 'r', [t('chn_kind_' + n.kind), n.role, n.id].filter(Boolean).join(' · ')
        + (ns ? ' · ' + t(ns.state === 'up' ? 'chn_up' : 'chn_down') : '')));
      if (ns && ns.error) card.appendChild(el('div', 'e', ns.error));
      (d.checks || []).filter((x) => x.node === n.id).forEach((x) => card.appendChild(checkRow(x, false)));
      grid.appendChild(card);
    });
    if ((d.links || []).length) {
      const ls = section(body, t('chn_links'));
      d.links.forEach((l) => {
        const row = el('div', 'chn-row');
        row.appendChild(el('span', 'chn-dot ' + dotOf(linkStatus(d, l))));
        row.appendChild(el('span', '', `${label(d, l.from)} → ${label(d, l.to)}`));
        row.appendChild(el('span', 'k', [t('chn_lk_' + (l.kind || 'data')), l.label].filter(Boolean).join(' · ')));
        ls.appendChild(row);
      });
    }
    if ((d.scenarios || []).length) {
      const ss = section(body, t('chn_scenarios'));
      d.scenarios.forEach((s) => {
        const row = el('div', 'chn-row');
        row.appendChild(el('span', 'chn-badge' + (s.id === c.scenario ? ' active' : ''), s.label || s.id));
        ruleText(row, s.when, d);
        ss.appendChild(row);
      });
    }
    const tl = section(body, t('chn_timeline'));
    tl.id = 'chn-tl';
  }

  function label(d, id) {
    const n = (d.nodes || []).find((x) => x.id === id);
    return n ? (n.label || n.id) : id;
  }

  function section(parent, title) {
    const s = el('div', 'chn-sec');
    s.appendChild(el('h4', '', title));
    parent.appendChild(s);
    return s;
  }

  // ─── timeline: the scenario band and the list of changes ───────────────────────────────

  function scenarioLabel(d, id) { const s = (d.scenarios || []).find((x) => x.id === id); return s ? (s.label || s.id) : (id || t('chn_no_scenario')); }
  function scenarioColor(d, id) { const i = (d.scenarios || []).findIndex((x) => x.id === id); return i < 0 ? '#8892a0' : SCN_COLORS[i % SCN_COLORS.length]; }

  function renderBand(d, list) {
    const band = $('#chn-band');
    if (!band) return;
    band.querySelectorAll('.chn-band-head').forEach((x) => x.remove());
    const head = el('div', 'chn-band-head');
    head.appendChild(el('h4', '', t('chn_band_t')));
    band.insertBefore(head, band.firstChild);
    const scn = list.filter((e) => e.kind === 'scenario').slice().reverse();
    if (!scn.length) { head.appendChild(el('div', 'chn-empty', t('chn_band_none'))); return; }
    const now = Date.now();
    const start = new Date(scn[0].at).getTime();
    const span = Math.max(1, now - start);
    const bar = el('div', 'chn-bar');
    bar.setAttribute('role', 'img');
    const seen = [];
    scn.forEach((e, i) => {
      const from = new Date(e.at).getTime();
      const to = i + 1 < scn.length ? new Date(scn[i + 1].at).getTime() : now;
      const seg = el('span');
      seg.style.flexGrow = String(Math.max(0, to - from));
      seg.style.background = scenarioColor(d, e.new);
      seg.title = scenarioLabel(d, e.new) + ' · ' + fmt(t('chn_since'), { at: when(e.at) });
      bar.appendChild(seg);
      if (seen.indexOf(e.new) < 0) seen.push(e.new);
    });
    bar.setAttribute('aria-label', t('chn_band_t') + ': ' + scn.map((e) => scenarioLabel(d, e.new)).join(' → '));
    head.appendChild(bar);
    const marks = el('div', 'chn-marks');
    list.filter((e) => e.kind !== 'scenario').slice(0, 200).forEach((e) => {
      const at = new Date(e.at).getTime();
      if (!(at >= start)) return;
      const m = el('span');
      m.style.left = (100 * (at - start) / span).toFixed(2) + '%';
      m.style.background = e.kind === 'node' ? (e.new === 'up' ? '#3fb950' : '#f85149') : (e.ok ? '#3fb950' : '#f85149');
      m.title = when(e.at) + ' · ' + changeText(d, e);
      marks.appendChild(m);
    });
    head.appendChild(marks);
    const axis = el('div', 'chn-axis');
    axis.appendChild(el('span', '', when(scn[0].at)));
    axis.appendChild(el('span', '', t('chn_now')));
    head.appendChild(axis);
    const keys = el('div', 'chn-keys');
    seen.forEach((id) => {
      const k = el('span'); const b = el('b'); b.style.background = scenarioColor(d, id);
      k.appendChild(b); k.appendChild(document.createTextNode(scenarioLabel(d, id))); keys.appendChild(k);
    });
    head.appendChild(keys);
  }

  function changeText(d, e) {
    let what = e.subject;
    if (e.kind === 'node') what = label(d, e.subject);
    if (e.kind === 'scenario') what = t('chn_scenario');
    if (e.kind === 'check') { const x = (d.checks || []).find((c) => c.id === e.subject); if (x) what = (x.label || x.id) + ' (' + label(d, x.node) + ')'; }
    const from = e.kind === 'scenario' ? (e.old ? scenarioLabel(d, e.old) : '—') : (e.old || '—');
    const to = e.kind === 'scenario' ? scenarioLabel(d, e.new) : e.new;
    return `${what}: ${from} → ${to}`;
  }

  function timelineRow(e) {
    const d = state.current.definition;
    const row = el('div', 'chn-tl');
    row.appendChild(el('span', 'a', when(e.at)));
    if (e.kind === 'check') row.appendChild(el('span', 'chn-dot ' + (e.ok ? 'ok' : 'bad')));
    if (e.kind === 'node') row.appendChild(el('span', 'chn-dot ' + (e.new === 'up' ? 'ok' : 'bad')));
    row.appendChild(el('span', '', changeText(d, e)));
    return row;
  }

  function renderTimeline(list) {
    if (!state.current) return;
    list = list || [];
    renderBand(state.current.definition, list);
    const box = $('#chn-tl');
    if (!box) return;
    box.textContent = '';
    box.appendChild(el('h4', '', t('chn_timeline')));
    if (!list.length) { box.appendChild(el('div', 'chn-empty', t('chn_timeline_none'))); return; }
    list.forEach((e) => box.appendChild(timelineRow(e)));
  }

  // ─── actions ────────────────────────────────────────────────────────────────────────────

  function checkNow(btn) {
    const c = state.current;
    if (!c) return;
    const label0 = btn.textContent;
    btn.disabled = true; btn.textContent = t('chn_checking');
    api('POST', `/api/chains/${c.id}/check`).then(() => loadList())
      .catch((e) => showToast(e.message, 'error'))
      .finally(() => { btn.disabled = false; btn.textContent = label0; });
  }

  function deleteChain() {
    const c = state.current;
    if (!c || !confirm(fmt(t('chn_delete_q'), { name: c.definition.name }))) return;
    api('DELETE', `/api/chains/${c.id}`).then(() => { delete state.views[c.id]; state.current = null; return loadList(); })
      .catch((e) => showToast(e.message, 'error'));
  }

  // Approving a draft is the only step that creates a chain and starts its background checks.
  function approveDraft() {
    const c = state.current;
    if (!c || !c.draft) return;
    if (!confirm(fmt(t('chn_approve_q'), { name: c.definition.name }))) return;
    const b = $('#chn-approve');
    b.disabled = true;
    api('POST', `/api/chains/drafts/${c.draftId}/approve`).then((chain) => {
      showToast(t('chn_draft_approved'), 'success');
      state.current = null;
      return loadList().then(() => select(chain.id));
    }).catch((e) => showToast((e.errors && e.errors[0]) || e.message, 'error'))
      .finally(() => { b.disabled = false; });
  }

  function discardDraft() {
    const c = state.current;
    if (!c || !c.draft) return;
    if (!confirm(fmt(t('chn_discard_q'), { name: c.definition.name }))) return;
    api('DELETE', `/api/chains/drafts/${c.draftId}`).then(() => { state.current = null; return loadList(); })
      .catch((e) => showToast(e.message, 'error'));
  }


  // ─── the guided builder, dialogs and the AI wizard (step 4) ──────────────────────────────

  function sshConns() { return (state.conns || []).filter((c) => (c.protocol || 'SSH').toUpperCase() === 'SSH'); }
  function catalogTypes() { return (state.catalog && state.catalog.types) || {}; }
  function nodeOf(def, id) { return (def.nodes || []).find((n) => n.id === id); }
  function nodeName(def, id) { const n = nodeOf(def, id); return n ? (n.label || n.id) : id; }
  function checkName(def, id) { const x = (def.checks || []).find((y) => y.id === id); return x ? (x.label || x.id) : id; }
  function slugify(s) {
    s = String(s || '').toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '');
    s = s.replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '').slice(0, 32);
    return s || 'node';
  }
  function uniqueID(base, used) {
    let id = slugify(base), i = 2;
    while (used[id]) { id = slugify(base).slice(0, 28) + '-' + i; i++; }
    used[id] = true;
    return id;
  }
  function usedNodeIDs(def) { const u = {}; (def.nodes || []).forEach((n) => { u[n.id] = true; }); return u; }
  function usedCheckIDs(def) { const u = {}; (def.checks || []).forEach((x) => { u[x.id] = true; }); return u; }
  function clone(v) { return JSON.parse(JSON.stringify(v)); }

  function defaultDef() {
    return { name: '', description: '', interval_s: 30, confirm: 2, nodes: [], checks: [], links: [], scenarios: [] };
  }

  function nodeSubtitle(n) {
    if (n.kind === 'server') {
      const c = (state.conns || []).find((x) => x.id === n.conn_id);
      return c ? c.name : t('chn_g_pick_server');
    }
    if (n.kind === 'storage') {
      const s = (state.storages || []).find((x) => x.id === n.storage_id);
      return s ? s.name : t('chn_g_pick_storage');
    }
    return t('chn_kind_external');
  }

  // autoScenarios builds the three states most chains need, from the actual checks and roles.
  function autoScenarios(def) {
    const servers = (def.nodes || []).filter((n) => n.kind === 'server');
    const out = [];
    if ((def.checks || []).length) {
      out.push({ id: 'normal', label: 'Normal', when: { all: (def.checks || []).map((x) => ({ check: x.id, field: 'ok', op: '==', value: 'true' })) } });
    }
    if (servers.length > 1) {
      out.push({ id: 'backup', label: 'Backup', when: { all: [
        { node: servers[0].id, op: '==', value: 'down' },
        { node: servers[1].id, op: '==', value: 'up' },
      ] } });
    }
    out.push({ id: 'degraded', label: 'Degraded', when: {} });
    return out;
  }

  function looksAuto(def) {
    const s = def.scenarios || [];
    if (!s.length) return true;
    const auto = autoScenarios(def);
    const shape = (list) => JSON.stringify(list.map((x) => ({ id: x.id, when: x.when })));
    return shape(s) === shape(auto);
  }

  function applyAuto(def) {
    const s = autoScenarios(def);
    // keep custom labels for the ids they still use
    const old = {};
    (def.scenarios || []).forEach((x) => { old[x.id] = x.label; });
    def.scenarios = s.map((x) => ({ id: x.id, label: old[x.id] || x.label, when: x.when }));
  }

  // ─── the guided editor ─────────────────────────────────────────────────────────────────

  function editCurrent() {
    const c = state.current;
    if (!c) return;
    if (c.draft) openGuided(null, { definition: c.definition, enabled: false }, c.draftId);
    else openGuided(c.id, { definition: c.definition, enabled: c.enabled });
  }

  function openGuided(id, value, draftId) {
    const ed = editorPane(draftId ? t('chn_g_edit_draft') : (id ? t('chn_edit_chain') : t('chn_g_new')));
    if (!ed) return;
    const def = clone(value.definition || defaultDef());
    if (!def.name) def.name = t('chn_g_new');
    def.nodes = def.nodes || []; def.checks = def.checks || []; def.links = def.links || []; def.scenarios = def.scenarios || [];
    let enabled = value.enabled !== false;
    let auto = looksAuto(def);
    const head = ed.querySelector('.chn-editor-head');
    const errBox = el('div', 'chn-errs');
    errBox.id = 'chn-ed-errs';
    errBox.setAttribute('role', 'status');
    errBox.setAttribute('aria-live', 'polite');
    const body = el('div', 'chn-guided');
    const main = el('div', 'chn-guided-main');
    const help = el('aside', 'chn-guided-help');
    help.appendChild(el('h4', '', t('chn_how')));
    help.appendChild(el('p', '', t('chn_how_t')));
    body.appendChild(main);
    body.appendChild(help);
    ed.appendChild(body);
    ed.appendChild(errBox);

    const cancel = el('button', 'btn-sec btn-sm', t('chn_g_cancel'));
    cancel.onclick = () => ed.remove();
    const adv = el('button', 'btn-sec btn-sm', '{ } ' + t('chn_g_json'));
    adv.onclick = () => { ed.remove(); openEditor(id, { definition: def, enabled }); };
    const save = el('button', 'btn-pri btn-sm', draftId ? t('chn_g_save_draft') : t('chn_g_save_chain'));
    head.appendChild(cancel); head.appendChild(adv); head.appendChild(save);

    const refresh = () => { errBox.textContent = ''; renderGuided(main, def, draftId, { enabled: () => enabled, setEnabled: (v) => { enabled = v; }, auto: () => auto, setAuto: (v) => { auto = v; }, refresh, close: () => ed.remove() }); };
    // Re-render keeps the scroll position of the pane, so adding a row does not jump to the top.
    const renderGuided = (root, d, isDraft, ctx) => {
      const scroll = root.scrollTop;
      root.textContent = '';
      guidedBasics(root, d, ctx);
      guidedNodes(root, d, ctx);
      guidedChecks(root, d, ctx);
      guidedLinks(root, d, ctx);
      guidedScenarios(root, d, ctx);
      root.scrollTop = scroll;
    };
    refresh();

    save.onclick = () => {
      errBox.textContent = '';
      if (!String(def.name || '').trim()) { showErrors(ed, [t('chn_g_required')]); return; }
      if (!(def.nodes || []).length) { showErrors(ed, [t('chn_g_required')]); return; }
      for (const n of def.nodes) {
        if (n.kind === 'server' && !n.conn_id) { showErrors(ed, [t('chn_g_node_empty')]); return; }
        if (n.kind === 'storage' && !n.storage_id) { showErrors(ed, [t('chn_g_node_empty')]); return; }
      }
      if (auto) applyAuto(def);
      saveGuided(id, draftId, def, enabled, ed, () => {});
    };
  }

  function fieldRow(parent, label, input, hint) {
    const l = el('label', 'chn-f');
    l.appendChild(el('span', '', label));
    l.appendChild(input);
    if (hint) l.appendChild(el('small', '', hint));
    parent.appendChild(l);
    return input;
  }

  function guidedBasics(root, def, ctx) {
    const sec = guidedSection(root, t('chn_g_basics'));
    const grid = el('div', 'chn-grid2');
    sec.appendChild(grid);
    const name = el('input');
    name.value = def.name || '';
    name.maxLength = 120;
    name.oninput = () => { def.name = name.value; };
    fieldRow(grid, t('chn_g_name'), name);
    const interval = el('select');
    [15, 30, 60, 120, 300, 600].forEach((n) => { const o = el('option', '', fmt(t('chn_g_seconds'), { n })); o.value = String(n); interval.appendChild(o); });
    interval.value = String(def.interval_s || 30);
    interval.onchange = () => { def.interval_s = Number(interval.value); };
    fieldRow(grid, t('chn_g_interval'), interval);
    const confirm = el('select');
    [1, 2, 3, 5].forEach((n) => { const o = el('option', '', fmt(t('chn_g_confirm_n'), { n })); o.value = String(n); confirm.appendChild(o); });
    confirm.value = String(def.confirm || 2);
    confirm.onchange = () => { def.confirm = Number(confirm.value); };
    fieldRow(grid, t('chn_g_confirm'), confirm);
    const en = el('label', 'chn-checkline');
    const cb = el('input'); cb.type = 'checkbox'; cb.checked = ctx.enabled();
    cb.onchange = () => ctx.setEnabled(cb.checked);
    en.appendChild(cb); en.appendChild(document.createTextNode(' ' + t('chn_g_enabled')));
    grid.appendChild(en);
    const desc = el('textarea');
    desc.rows = 2;
    desc.value = def.description || '';
    desc.oninput = () => { def.description = desc.value; };
    fieldRow(sec, t('chn_g_desc'), desc);
  }

  function guidedSection(root, title) {
    const s = el('section', 'chn-gsec');
    s.appendChild(el('h4', '', title));
    root.appendChild(s);
    return s;
  }

  function guidedNodes(root, def, ctx) {
    const sec = guidedSection(root, t('chn_g_nodes'));
    const list = el('div', 'chn-glist');
    sec.appendChild(list);
    if (!(def.nodes || []).length) list.appendChild(el('div', 'note', t('chn_g_node_empty')));
    def.nodes.forEach((n, idx) => {
      const row = el('div', 'chn-grow');
      const icon = el('span', 'chn-gicon', n.kind === 'server' ? '🖥' : n.kind === 'storage' ? '🪣' : '🌐');
      rows: {
        row.appendChild(icon);
        const bits = el('div', 'chn-grow-main');
        const top = el('div', 'chn-grow-top');
        const lab = el('input');
        lab.value = n.label || n.id;
        lab.maxLength = 120;
        lab.oninput = () => { n.label = lab.value; };
        top.appendChild(lab);
        const role = el('input');
        role.value = n.role || '';
        role.maxLength = 60;
        role.placeholder = t('chn_g_role_ph');
        role.oninput = () => { n.role = role.value; };
        role.className = 'chn-role';
        top.appendChild(role);
        bits.appendChild(top);
        bits.appendChild(el('div', 'chn-grow-sub', t('chn_kind_' + n.kind) + ' · ' + nodeSubtitle(n)));
        row.appendChild(bits);
      }
      const del = el('button', 'btn-sec btn-sm', '✕');
      del.title = t('chn_g_remove');
      del.onclick = () => {
        def.nodes.splice(idx, 1);
        def.checks = def.checks.filter((x) => x.node !== n.id);
        def.links = def.links.filter((l) => l.from !== n.id && l.to !== n.id);
        def.scenarios = def.scenarios.filter((s) => !condUsesNode(s.when, n.id));
        ctx.refresh();
      };
      row.appendChild(del);
      list.appendChild(row);
      const checks = (def.checks || []).filter((x) => x.node === n.id);
      if (checks.length) {
        const chips = el('div', 'chn-gchips');
        checks.forEach((x) => chips.appendChild(el('span', 'chn-chip', x.label || x.type)));
        list.appendChild(chips);
      }
    });
    guidedAddNode(sec, def, ctx);
  }

  function condUsesNode(when, id) {
    if (!when) return false;
    if (when.node === id) return true;
    return (when.all || []).some((x) => condUsesNode(x, id)) || (when.any || []).some((x) => condUsesNode(x, id));
  }

  function guidedAddNode(sec, def, ctx) {
    const box = el('div', 'chn-gadd');
    const kind = el('select');
    [['server', t('chn_g_kind') + ': ' + t('chn_kind_server')], ['storage', t('chn_g_kind') + ': ' + t('chn_kind_storage')], ['external', t('chn_g_kind') + ': ' + t('chn_kind_external')]]
      .forEach(([v, l]) => { const o = el('option', '', l); o.value = v; kind.appendChild(o); });
    const pick = el('select');
    const label = el('input');
    label.placeholder = t('chn_g_label');
    const role = el('input');
    role.placeholder = t('chn_g_role_ph');
    const err = el('span', 'chn-inline-err');
    const fillPick = () => {
      pick.textContent = '';
      if (kind.value === 'server') {
        pick.appendChild(el('option', '', t('chn_g_pick_server'))).value = '';
        sshConns().forEach((c) => { const o = el('option', '', c.name + (c.host ? ' (' + c.host + ')' : '')); o.value = String(c.id); pick.appendChild(o); });
        pick.style.display = '';
      } else if (kind.value === 'storage') {
        pick.appendChild(el('option', '', t('chn_g_pick_storage'))).value = '';
        (state.storages || []).forEach((s) => { const o = el('option', '', s.name); o.value = String(s.id); pick.appendChild(o); });
        pick.style.display = '';
      } else {
        pick.style.display = 'none';
      }
    };
    kind.onchange = fillPick;
    fillPick();
    const add = el('button', 'btn-sec btn-sm', '＋ ' + t('chn_g_add'));
    add.onclick = () => {
      err.textContent = '';
      const id = uniqueID(label.value || (kind.value === 'server' ? 'server' : kind.value === 'storage' ? 'storage' : 'external'), usedNodeIDs(def));
      const n = { id, kind: kind.value, label: label.value.trim() || id, role: role.value.trim() };
      if (kind.value === 'server') {
        const v = Number(pick.value) || 0;
        if (!v) { err.textContent = t('chn_g_pick_server'); return; }
        n.conn_id = v;
      } else if (kind.value === 'storage') {
        const v = Number(pick.value) || 0;
        if (!v) { err.textContent = t('chn_g_pick_storage'); return; }
        n.storage_id = v;
      }
      def.nodes.push(n);
      ctx.refresh();
    };
    [kind, pick, label, role, add].forEach((x) => box.appendChild(x));
    box.appendChild(err);
    sec.appendChild(box);
  }

  function guidedChecks(root, def, ctx) {
    const sec = guidedSection(root, t('chn_g_checks'));
    if (!(def.checks || []).length) sec.appendChild(el('div', 'note', t('chn_g_no_checks')));
    (def.nodes || []).forEach((n) => {
      const group = el('div', 'chn-cgroup');
      group.appendChild(el('h5', '', n.label || n.id));
      (def.checks || []).filter((x) => x.node === n.id).forEach((x, i) => {
        const row = el('div', 'chn-grow');
        row.appendChild(el('span', 'chn-gicon', '✓'));
        const main = el('div', 'chn-grow-main');
        main.appendChild(el('div', 'chn-grow-top', (x.label || x.id) + ' · ' + x.type));
        main.appendChild(el('div', 'chn-grow-sub', checkSummary(x)));
        row.appendChild(main);
        const editB = el('button', 'btn-sec btn-sm', '✎');
        editB.title = t('chn_g_edit');
        editB.onclick = () => openCheckDialog(def, n, x, (chk) => {
          const at = def.checks.findIndex((y) => y.id === x.id);
          chk.id = x.id;
          if (at >= 0) def.checks[at] = chk;
          ctx.refresh();
        });
        const delB = el('button', 'btn-sec btn-sm', '✕');
        delB.title = t('chn_g_remove');
        delB.onclick = () => { def.checks = def.checks.filter((y) => y.id !== x.id); def.links.forEach((l) => { if (l.check === x.id) delete l.check; }); ctx.refresh(); };
        row.appendChild(editB); row.appendChild(delB);
        group.appendChild(row);
      });
      const add = el('button', 'btn-sec btn-sm', '＋ ' + t('chn_g_add_check'));
      add.onclick = () => openCheckDialog(def, n, null, (chk) => { def.checks.push(chk); ctx.refresh(); });
      group.appendChild(add);
      sec.appendChild(group);
    });
  }

  function checkSummary(x) {
    const parts = [x.target, x.port ? ':' + x.port : '', x.key ? x.key : '', x.prefix ? x.prefix : ''];
    return parts.filter(Boolean).join(' · ') || '—';
  }

  function guidedLinks(root, def, ctx) {
    const sec = guidedSection(root, t('chn_g_links'));
    if (!(def.links || []).length) sec.appendChild(el('div', 'note', t('chn_g_link_none')));
    (def.links || []).forEach((l, i) => {
      const row = el('div', 'chn-grow');
      row.appendChild(el('span', 'chn-gicon', '⇢'));
      const main = el('div', 'chn-grow-main');
      main.appendChild(el('div', 'chn-grow-top', nodeName(def, l.from) + ' → ' + nodeName(def, l.to)));
      main.appendChild(el('div', 'chn-grow-sub', [t('chn_lk_' + (l.kind || 'data')), l.label, l.check ? checkName(def, l.check) : ''].filter(Boolean).join(' · ')));
      row.appendChild(main);
      const del = el('button', 'btn-sec btn-sm', '✕');
      del.title = t('chn_g_remove');
      del.onclick = () => { def.links.splice(i, 1); ctx.refresh(); };
      row.appendChild(del);
      sec.appendChild(row);
    });
    const add = el('button', 'btn-sec btn-sm', '＋ ' + t('chn_g_add_link'));
    add.onclick = () => openLinkDialog(def, null, (l) => { def.links.push(l); ctx.refresh(); });
    sec.appendChild(add);
  }

  function guidedScenarios(root, def, ctx) {
    const sec = guidedSection(root, t('chn_g_scenarios'));
    const autoL = el('label', 'chn-checkline');
    const cb = el('input'); cb.type = 'checkbox'; cb.checked = ctx.auto();
    cb.onchange = () => { ctx.setAuto(cb.checked); ctx.refresh(); };
    autoL.appendChild(cb); autoL.appendChild(document.createTextNode(' ' + t('chn_g_auto_scenarios')));
    sec.appendChild(autoL);
    if (ctx.auto()) {
      sec.appendChild(el('div', 'note', t('chn_g_auto_scenarios_d')));
      applyAuto(def);
      listScenarios(sec, def, false, ctx);
      return;
    }
    listScenarios(sec, def, true, ctx);
    const add = el('button', 'btn-sec btn-sm', '＋ ' + t('chn_g_add_scenario'));
    add.onclick = () => {
      const used = {}; def.scenarios.forEach((s) => { used[s.id] = true; });
      def.scenarios.push({ id: uniqueID('state', used), label: 'State ' + (def.scenarios.length + 1), when: {} });
      ctx.refresh();
    };
    sec.appendChild(add);
  }

  function listScenarios(sec, def, editable, ctx) {
    (def.scenarios || []).forEach((sc, si) => {
      const box = el('div', 'chn-scn-box');
      const head = el('div', 'chn-grow');
      const main = el('div', 'chn-grow-main');
      if (editable) {
        const lab = el('input');
        lab.value = sc.label || sc.id;
        lab.oninput = () => { sc.label = lab.value; };
        main.appendChild(lab);
      } else {
        main.appendChild(el('div', 'chn-grow-top', sc.label || sc.id));
      }
      head.appendChild(main);
      if (editable) {
        const up = el('button', 'btn-sec btn-sm', '↑');
        up.onclick = () => { if (si > 0) { const [x] = def.scenarios.splice(si, 1); def.scenarios.splice(si - 1, 0, x); ctx.refresh(); } };
        const down = el('button', 'btn-sec btn-sm', '↓');
        down.onclick = () => { if (si < def.scenarios.length - 1) { const [x] = def.scenarios.splice(si, 1); def.scenarios.splice(si + 1, 0, x); ctx.refresh(); } };
        const del = el('button', 'btn-sec btn-sm', '✕');
        del.onclick = () => { def.scenarios.splice(si, 1); ctx.refresh(); };
        head.appendChild(up); head.appendChild(down); head.appendChild(del);
      }
      box.appendChild(head);
      box.appendChild(el('div', 'chn-grow-sub', t('chn_g_when') + ':'));
      const conds = [];
      const hasAny = !!(sc.when && (sc.when.any || []).length);
      flattenConds(sc.when, conds);
      if (!conds.length) box.appendChild(el('div', 'note', t('chn_g_cond_none')));
      if (editable && hasAny) box.appendChild(el('div', 'note', t('chn_g_any_note')));
      if (editable && !hasAny) conds.forEach((c) => box.appendChild(condRow(def, sc, c, ctx)));
      else conds.forEach((c) => box.appendChild(el('div', 'chn-cond', condText(def, c))));
      if (editable && !hasAny) {
        const add = el('button', 'btn-sec btn-sm', '＋ ' + t('chn_g_cond_add'));
        add.onclick = () => {
          const first = (def.checks || [])[0];
          const c = first ? { check: first.id, field: 'ok', op: '==', value: 'true' } : { node: (def.nodes.find((n) => n.kind === 'server') || {}).id, op: '==', value: 'up' };
          sc.when.all = (sc.when.all || []).concat([c]);
          ctx.refresh();
        };
        box.appendChild(add);
      }
      sec.appendChild(box);
    });
  }

  // flattenConds is only for display and simple editing: the guided builder edits "all" rules.
  function flattenConds(when, out) {
    if (!when) return;
    if ((when.all || []).length) when.all.forEach((x) => flattenConds(x, out));
    else if ((when.any || []).length) when.any.forEach((x) => flattenConds(x, out));
    else if (when.check || when.node) out.push(when);
  }

  function condText(def, c) {
    const subject = c.check ? checkName(def, c.check) + (c.field === 'ok' ? ' · ' + t('chn_g_ok') : '') : nodeName(def, c.node);
    return subject + ' ' + (c.op || '==') + ' ' + (c.value || '');
  }

  function condRow(def, sc, c, ctx) {
    const row = el('div', 'chn-cond-row');
    const subj = el('select');
    (def.checks || []).forEach((x) => { const o = el('option', '', t('chn_g_checks') + ': ' + (x.label || x.id)); o.value = 'check:' + x.id; if (c.check === x.id) o.selected = true; subj.appendChild(o); });
    (def.nodes || []).filter((n) => n.kind === 'server').forEach((n) => { const o = el('option', '', t('chn_g_nodes') + ': ' + (n.label || n.id)); o.value = 'node:' + n.id; if (c.node === n.id) o.selected = true; subj.appendChild(o); });
    subj.onchange = () => {
      const [kind, id] = subj.value.split(':');
      delete c.check; delete c.node; delete c.field;
      if (kind === 'check') { c.check = id; c.field = 'ok'; c.op = '=='; c.value = 'true'; }
      else { c.node = id; c.op = '=='; c.value = 'up'; }
      ctx.refresh();
    };
    row.appendChild(subj);
    if (c.check) {
      const fld = el('select');
      [['ok', t('chn_g_ok')], ['value', t('chn_g_value')]].forEach(([v, l]) => { const o = el('option', '', l); o.value = v; if ((c.field || 'value') === v) o.selected = true; fld.appendChild(o); });
      fld.onchange = () => { c.field = fld.value; c.value = fld.value === 'ok' ? 'true' : ''; ctx.refresh(); };
      row.appendChild(fld);
    }
    const op = el('select');
    ['==', '!=', 'contains', 'exists'].forEach((v) => { const o = el('option', '', v); o.value = v; if (c.op === v) o.selected = true; op.appendChild(o); });
    op.onchange = () => { c.op = op.value; ctx.refresh(); };
    row.appendChild(op);
    if (c.op === 'exists') { c.value = ''; }
    else if (c.node) {
      const v = el('select');
      [['up', t('chn_g_up')], ['down', t('chn_g_down')]].forEach(([val, l]) => { const o = el('option', '', l); o.value = val; if (c.value === val) o.selected = true; v.appendChild(o); });
      v.onchange = () => { c.value = v.value; };
      row.appendChild(v);
    } else if (c.check && c.field === 'ok') {
      const v = el('select');
      [['true', t('chn_g_ok')], ['false', t('chn_g_down')]].forEach(([val, l]) => { const o = el('option', '', l); o.value = val; if (c.value === val) o.selected = true; v.appendChild(o); });
      v.onchange = () => { c.value = v.value; };
      row.appendChild(v);
    } else {
      const v = el('input');
      v.value = c.value || '';
      v.placeholder = t('chn_g_value');
      v.oninput = () => { c.value = v.value; };
      row.appendChild(v);
    }
    const del = el('button', 'btn-sec btn-sm', '✕');
    del.onclick = () => {
      const list = (sc.when.all || []);
      sc.when.all = list.filter((x) => x !== c);
      if (!sc.when.all.length) delete sc.when.all;
      ctx.refresh();
    };
    row.appendChild(del);
    return row;
  }

  function saveGuided(id, draftId, def, enabled, ed, done) {
    // Clean empty optional values so the server sees the same definition the user sees.
    (def.checks || []).forEach((x) => { ['target', 'key', 'record', 'resolver', 'prefix', 'label'].forEach((k) => { if (typeof x[k] === 'string') x[k] = x[k].trim(); }); });
    const body = { definition: def, enabled };
    const req = draftId ? api('PUT', `/api/chains/drafts/${draftId}`, { definition: def })
      : id ? api('PUT', `/api/chains/${id}`, body) : api('POST', '/api/chains', body);
    req.then((c) => {
      if (ed) ed.remove();
      showToast(t('chn_saved'), 'success');
      if (draftId) { return selectDraft(draftId); }
      state.current = c; state.sel = null; state.pos = {};
      return loadList().then(() => select(c.id));
    }).catch((e) => showErrors(ed, e.errors.length ? e.errors : [e.message]));
  }

  // ─── small dialogs over the window (add node, check, link) ─────────────────────────────

  function dialogShell(title) {
    const wrap = el('div', 'chn-dlg-wrap');
    const card = el('div', 'chn-dlg');
    const head = el('div', 'chn-dlg-head');
    head.appendChild(el('h4', '', title));
    const x = el('button', 'btn-sec btn-sm', '✕');
    head.appendChild(x);
    card.appendChild(head);
    const body = el('div', 'chn-dlg-body');
    card.appendChild(body);
    wrap.appendChild(card);
    const main = $('.chn-main') || document.body;
    main.appendChild(wrap);
    const close = () => wrap.remove();
    x.onclick = close;
    wrap.addEventListener('mousedown', (e) => { if (e.target === wrap) close(); });
    return { wrap, body, close };
  }

  function openNodeDialog(def, onAdd) {
    const dlg = dialogShell(t('chn_add_node'));
    const form = el('div', 'chn-form');
    dlg.body.appendChild(form);
    const kind = el('select');
    [['server', t('chn_kind_server')], ['storage', t('chn_kind_storage')], ['external', t('chn_kind_external')]].forEach(([v, l]) => { const o = el('option', '', l); o.value = v; kind.appendChild(o); });
    fieldRow(form, t('chn_g_kind'), kind);
    const pickWrap = el('div');
    form.appendChild(pickWrap);
    const label = el('input'); label.placeholder = t('chn_g_label');
    const role = el('input'); role.placeholder = t('chn_g_role_ph');
    fieldRow(form, t('chn_g_label'), label);
    fieldRow(form, t('chn_g_role'), role);
    const pick = el('select');
    const fill = () => {
      pickWrap.textContent = '';
      if (kind.value === 'server') {
        pick.textContent = ''; pick.appendChild(el('option', '', t('chn_g_pick_server'))).value = '';
        sshConns().forEach((c) => { const o = el('option', '', c.name + (c.host ? ' (' + c.host + ')' : '')); o.value = String(c.id); pick.appendChild(o); });
        fieldRow(pickWrap, t('chn_kind_server'), pick);
      } else if (kind.value === 'storage') {
        pick.textContent = ''; pick.appendChild(el('option', '', t('chn_g_pick_storage'))).value = '';
        (state.storages || []).forEach((s) => { const o = el('option', '', s.name); o.value = String(s.id); pick.appendChild(o); });
        fieldRow(pickWrap, t('chn_kind_storage'), pick);
      }
    };
    kind.onchange = fill; fill();
    const err = el('div', 'chn-inline-err'); form.appendChild(err);
    const add = el('button', 'btn-pri btn-sm', t('chn_g_add'));
    add.onclick = () => {
      err.textContent = '';
      const id = uniqueID(label.value || kind.value, usedNodeIDs(def));
      const n = { id, kind: kind.value, label: label.value.trim() || id, role: role.value.trim() };
      if (kind.value === 'server') { const v = Number(pick.value) || 0; if (!v) { err.textContent = t('chn_g_pick_server'); return; } n.conn_id = v; }
      if (kind.value === 'storage') { const v = Number(pick.value) || 0; if (!v) { err.textContent = t('chn_g_pick_storage'); return; } n.storage_id = v; }
      dlg.close();
      onAdd(n);
    };
    form.appendChild(add);
  }

  function openCheckDialog(def, node, existing, onSave) {
    const dlg = dialogShell(existing ? t('chn_g_edit') + ' · ' + (existing.label || existing.type) : t('chn_g_add_check') + ' · ' + (node.label || node.id));
    const form = el('div', 'chn-form');
    dlg.body.appendChild(form);
    const types = catalogTypes();
    const chk = existing ? clone(existing) : { node: node.id, type: '', label: '', target: '' };
    const allowed = Object.keys(types).filter((k) => {
      const where = types[k].where;
      return node.kind === 'server' ? where === 'server' : node.kind === 'storage' ? where === 'storage' : where === 'wrm';
    }).sort();
    const render = () => {
      form.textContent = '';
      const type = el('select');
      type.appendChild(el('option', '', '—')).value = '';
      allowed.forEach((k) => { const o = el('option', '', k + (LANGS.en['chn_type_' + k] ? ' — ' + t('chn_type_' + k) : '')); o.value = k; if (chk.type === k) o.selected = true; type.appendChild(o); });
      type.onchange = () => { chk.type = type.value; render(); };
      fieldRow(form, t('chn_g_type'), type);
      if (!chk.type) return;
      const info = types[chk.type] || {};
      const need = String(info.needs || '');
      const lab = el('input'); lab.value = chk.label || ''; lab.placeholder = t('chn_g_check_label');
      lab.oninput = () => { chk.label = lab.value; };
      fieldRow(form, t('chn_g_check_label'), lab);
      if (chk.type === 'port') {
        const p = el('input'); p.type = 'number'; p.min = 1; p.max = 65535; p.value = chk.port || '';
        p.oninput = () => { chk.port = Number(p.value) || 0; };
        fieldRow(form, t('chn_g_port'), p);
      } else if (chk.type !== 'monitor') {
        const tgt = el('input'); tgt.value = chk.target || '';
        tgt.placeholder = need.includes('target') ? (chk.type === 'dns' ? 'app.example.com' : chk.type === 'tcp' ? 'host:443' : chk.type === 'http' ? 'https://…' : chk.type === 'tls' ? 'example.com' : chk.type === 'systemd' ? 'nginx.service' : '/absolute/path') : '';
        tgt.oninput = () => { chk.target = tgt.value; };
        fieldRow(form, t('chn_g_target'), tgt);
      }
      if (chk.type === 'config') {
        const fmt = el('select');
        ['ini', 'json', 'yaml', 'env', 'kv', 'regex'].forEach((v) => { const o = el('option', '', v); o.value = v; if (chk.format === v) o.selected = true; fmt.appendChild(o); });
        fmt.onchange = () => { chk.format = fmt.value; };
        fieldRow(form, t('chn_g_format'), fmt);
      }
      if (chk.type === 'config' || chk.type === 'http') {
        const key = el('input'); key.value = chk.key || '';
        key.oninput = () => { chk.key = key.value; };
        fieldRow(form, t('chn_g_key'), key);
      }
      if (chk.type === 'dns') {
        const rec = el('select');
        ['A', 'AAAA', 'CNAME', 'TXT', 'MX', 'NS'].forEach((v) => { const o = el('option', '', v); o.value = v; if (chk.record === v) o.selected = true; rec.appendChild(o); });
        rec.onchange = () => { chk.record = rec.value; };
        fieldRow(form, t('chn_g_record'), rec);
        const res = el('input'); res.value = chk.resolver || '';
        res.oninput = () => { chk.resolver = res.value; };
        fieldRow(form, t('chn_g_resolver'), res);
      }
      if (chk.type === 'storage') {
        const pre = el('input'); pre.value = chk.prefix || '';
        pre.oninput = () => { chk.prefix = pre.value; };
        fieldRow(form, t('chn_g_prefix'), pre);
      }
      if (chk.type === 'file' || chk.type === 'storage') {
        const age = el('input'); age.type = 'number'; age.min = 0; age.value = chk.max_age_s || '';
        age.oninput = () => { chk.max_age_s = Number(age.value) || 0; };
        fieldRow(form, t('chn_g_max_age'), age);
      }
      if (chk.type === 'tls') {
        const w = el('input'); w.type = 'number'; w.min = 0; w.value = chk.warn_days || '';
        w.oninput = () => { chk.warn_days = Number(w.value) || 0; };
        fieldRow(form, t('chn_g_warn_days'), w);
      }
      const ex = el('div', 'chn-expect');
      const op = el('select');
      [['', t('chn_g_expect_none')], ['==', '=='], ['!=', '!='], ['contains', 'contains'], ['exists', 'exists']].forEach(([v, l]) => { const o = el('option', '', l); o.value = v; if ((chk.expect && chk.expect.op) === v || (!chk.expect && v === '')) o.selected = true; op.appendChild(o); });
      op.onchange = () => { chk.expect = op.value ? { op: op.value, value: (chk.expect && chk.expect.value) || '' } : undefined; render(); };
      fieldRow(form, t('chn_g_expect'), op);
      if (chk.expect && chk.expect.op !== 'exists') {
        const val = el('input'); val.value = chk.expect.value || '';
        val.oninput = () => { chk.expect.value = val.value; };
        fieldRow(form, t('chn_g_value'), val);
      }
      const err = el('div', 'chn-inline-err');
      const save = el('button', 'btn-pri btn-sm', t('chn_g_save'));
      save.onclick = () => {
        err.textContent = '';
        if (need.includes('target') && !String(chk.target || '').trim()) { err.textContent = t('chn_g_target'); return; }
        if (chk.type === 'port' && (!chk.port || chk.port < 1)) { err.textContent = t('chn_g_port'); return; }
        if ((chk.type === 'config') && (!chk.format || !String(chk.key || '').trim())) { err.textContent = t('chn_g_key'); return; }
        if (!String(chk.label || '').trim()) chk.label = chk.type + (chk.port ? ' ' + chk.port : chk.target ? ' ' + trunc(chk.target, 26) : '');
        if (!chk.id) chk.id = uniqueID(node.id + '-' + chk.type, usedCheckIDs(def));
        dlg.close();
        onSave(chk);
      };
      form.appendChild(err); form.appendChild(save);
    };
    render();
  }

  function openLinkDialog(def, preset, onSave) {
    const dlg = dialogShell(t('chn_g_add_link'));
    const form = el('div', 'chn-form');
    dlg.body.appendChild(form);
    if ((def.nodes || []).length < 2) { form.appendChild(el('div', 'note', t('chn_g_node_empty'))); return; }
    const from = el('select'), to = el('select'), kind = el('select'), label = el('input'), check = el('select');
    (def.nodes || []).forEach((n) => { [from, to].forEach((s) => { const o = el('option', '', n.label || n.id); o.value = n.id; s.appendChild(o); }); });
    if (preset && preset.from) from.value = preset.from;
    if (preset && preset.to) to.value = preset.to;
    [['data', t('chn_lk_data')], ['read', t('chn_lk_read')], ['write', t('chn_lk_write')], ['failover', t('chn_lk_failover')], ['dependency', t('chn_lk_dependency')]].forEach(([v, l]) => { const o = el('option', '', l); o.value = v; kind.appendChild(o); });
    fieldRow(form, t('chn_g_from'), from);
    fieldRow(form, t('chn_g_to'), to);
    fieldRow(form, t('chn_g_link_kind'), kind);
    fieldRow(form, t('chn_g_label'), label);
    check.appendChild(el('option', '', '—')).value = '';
    (def.checks || []).forEach((x) => { const o = el('option', '', (x.label || x.id) + ' · ' + (nodeName(def, x.node))); o.value = x.id; check.appendChild(o); });
    fieldRow(form, t('chn_g_link_check'), check);
    const err = el('div', 'chn-inline-err');
    const add = el('button', 'btn-pri btn-sm', t('chn_g_add'));
    add.onclick = () => {
      if (from.value === to.value) { err.textContent = t('chn_g_to'); return; }
      if ((def.links || []).some((l) => l.from === from.value && l.to === to.value)) { err.textContent = t('chn_link_exists'); return; }
      const l = { from: from.value, to: to.value, kind: kind.value };
      if (label.value.trim()) l.label = label.value.trim();
      if (check.value) l.check = check.value;
      dlg.close();
      onSave(l);
    };
    form.appendChild(err); form.appendChild(add);
  }

  // applyDefinition saves a structural change made on the diagram (node, link, check).
  function applyDefinition(mutator) {
    const c = state.current;
    if (!c) return;
    const def = clone(c.definition);
    mutator(def);
    const req = c.draft ? api('PUT', `/api/chains/drafts/${c.draftId}`, { definition: def })
      : api('PUT', `/api/chains/${c.id}`, { definition: def, enabled: c.enabled });
    req.then(() => {
      showToast(t('chn_link_saved'), 'success');
      if (c.draft) return selectDraft(c.draftId, true);
      return select(c.id, true);
    }).catch((e) => showToast((e.errors && e.errors[0]) || e.message, 'error'));
  }

  function editorPane(title) {
    const main = $('.chn-main');
    if (!main) return null;
    const old = $('.chn-editor');
    if (old) old.remove();
    const ed = el('div', 'chn-editor');
    const head = el('div', 'chn-editor-head');
    head.appendChild(el('h3', '', title));
    ed.appendChild(head);
    main.appendChild(ed);
    return ed;
  }

  // New chain: the guided choices come first; templates are only a shortcut.
  function openNew() {
    const ed = editorPane(t('chn_build'));
    if (!ed) return;
    const head = ed.querySelector('.chn-editor-head');
    const cancel = el('button', 'btn-sec btn-sm', t('chn_g_cancel'));
    cancel.onclick = () => ed.remove();
    head.appendChild(cancel);
    const body = el('div', 'chn-build');
    ed.appendChild(body);
    body.appendChild(el('h4', '', t('chn_how')));
    body.appendChild(el('p', 'note', t('chn_how_t')));
    const grid = el('div', 'chn-build-grid');
    body.appendChild(grid);
    const card = (icon, title, desc, fn) => {
      const b = el('button', 'chn-build-card');
      b.type = 'button';
      b.appendChild(el('span', 'ic', icon));
      b.appendChild(el('b', '', title));
      if (desc) b.appendChild(el('span', '', desc));
      b.onclick = fn;
      grid.appendChild(b);
    };
    card('🤖', t('chn_build_ai'), t('chn_build_ai_d'), () => openWizard(ed));
    card('🧩', t('chn_build_manual'), t('chn_build_manual_d'), () => { ed.remove(); openGuided(null, { definition: defaultDef(), enabled: true }); });
    card('📋', t('chn_build_tpl'), '', () => templateChooser(ed));
    card('{ }', t('chn_build_adv'), '', () => { ed.remove(); openEditor(null, { definition: defaultDef(), enabled: true }); });
  }

  // The AI wizard: pick servers/storages, describe the system, get a draft to review.
  function openWizard(ed) {
    Array.from(ed.children).forEach((c) => { if (!c.classList.contains('chn-editor-head')) c.remove(); });
    const body = el('div', 'chn-wizard');
    ed.appendChild(body);
    body.appendChild(el('h4', '', t('chn_build_ai')));
    const note = el('div', 'note');
    note.textContent = t('chn_ai_note');
    body.appendChild(note);
    const ssh = sshConns();
    body.appendChild(el('h5', '', t('chn_ai_pick_servers')));
    const servers = el('div', 'chn-picklist');
    if (!ssh.length) body.appendChild(el('div', 'note', t('chn_ai_servers_none')));
    ssh.forEach((c) => {
      const l = el('label', 'chn-pick');
      const cb = el('input'); cb.type = 'checkbox'; cb.value = String(c.id);
      l.appendChild(cb);
      l.appendChild(el('span', '', c.name + (c.host ? ' (' + c.host + ')' : '')));
      servers.appendChild(l);
    });
    body.appendChild(servers);
    const discoverBtn = el('button', 'btn-sec btn-sm', '🔍 ' + t('chn_ai_analyze'));
    body.appendChild(discoverBtn);
    const factsBox = el('div', 'chn-facts');
    body.appendChild(factsBox);
    discoverBtn.onclick = () => {
      const ids = Array.from(servers.querySelectorAll('input:checked')).map((x) => Number(x.value));
      if (!ids.length) { factsBox.textContent = t('chn_ai_servers_none'); return; }
      const old = discoverBtn.textContent;
      discoverBtn.disabled = true;
      discoverBtn.textContent = t('chn_ai_analyzing');
      factsBox.textContent = '';
      api('POST', '/api/chains/discover', { conn_ids: ids }).then((res) => {
        (res.servers || []).forEach((f) => {
          const card = el('div', 'chn-fact');
          card.appendChild(el('b', '', f.name + (f.hostname ? ' · ' + f.hostname : '')));
          if (f.os) card.appendChild(el('div', 'sub', f.os));
          if (f.error) card.appendChild(el('div', 'sub', '⚠ ' + f.error));
          const line = (label, items) => {
            if (!items || !items.length) return;
            const row = el('div', 'chn-fact-row');
            row.appendChild(el('span', 'k', label + ':'));
            items.slice(0, 12).forEach((x) => row.appendChild(el('span', 'chn-chip', String(x))));
            card.appendChild(row);
          };
          line(t('chn_ai_found_services'), f.services);
          line(t('chn_ai_found_ports'), f.ports);
          line(t('chn_ai_found_mounts'), (f.mounts || []).map((m) => m.target + (m.source ? ' ← ' + m.source : '')));
          if (!f.services?.length && !f.ports?.length && !f.mounts?.length) card.appendChild(el('div', 'note', t('chn_ai_no_facts')));
          factsBox.appendChild(card);
        });
      }).catch((e) => { factsBox.textContent = e.message; })
        .finally(() => { discoverBtn.disabled = false; discoverBtn.textContent = old; });
    };
    body.appendChild(el('h5', '', t('chn_ai_pick_storages')));
    const stores = el('div', 'chn-picklist');
    (state.storages || []).forEach((s) => {
      const l = el('label', 'chn-pick');
      const cb = el('input'); cb.type = 'checkbox'; cb.value = String(s.id);
      l.appendChild(cb); l.appendChild(el('span', '', s.name + (s.bucket ? ' · ' + s.bucket : '')));
      stores.appendChild(l);
    });
    body.appendChild(stores);
    body.appendChild(el('div', 'note', t('chn_ai_storage_optional')));
    const grid = el('div', 'chn-grid2');
    body.appendChild(grid);
    const name = el('input'); name.maxLength = 120; name.placeholder = t('chn_ai_name');
    fieldRow(grid, t('chn_g_name'), name);
    const interval = el('select');
    [15, 30, 60, 120, 300, 600].forEach((n) => { const o = el('option', '', fmt(t('chn_g_seconds'), { n })); o.value = String(n); interval.appendChild(o); });
    interval.value = '30';
    fieldRow(grid, t('chn_g_interval'), interval);
    const desc = el('textarea');
    desc.rows = 4;
    desc.maxLength = 2000;
    desc.placeholder = t('chn_ai_desc_ph');
    fieldRow(body, t('chn_ai_desc'), desc);
    const err = el('div', 'chn-errs');
    err.setAttribute('role', 'status');
    body.appendChild(err);
    const make = el('button', 'btn-pri btn-sm', t('chn_ai_make'));
    body.appendChild(make);
    make.onclick = () => {
      err.textContent = '';
      const connIDs = servers.querySelectorAll('input:checked');
      if (!connIDs.length) { err.textContent = t('chn_ai_servers_none'); return; }
      const storageIDs = Array.from(stores.querySelectorAll('input:checked')).map((x) => Number(x.value));
      const payload = {
        name: name.value.trim(), description: desc.value.trim(),
        conn_ids: Array.from(connIDs).map((x) => Number(x.value)), storage_ids: storageIDs,
        interval_s: Number(interval.value),
      };
      make.disabled = true;
      const old = make.textContent;
      make.textContent = t('chn_ai_working');
      api('POST', '/api/chains/ai-draft', payload).then((res) => {
        ed.remove();
        const reason = res.reason || '';
        const msg = res.ai ? t('chn_ai_ready')
          : (reason === 'provider_error' || reason === 'invalid_answer' || reason === 'invalid_definition') ? t('chn_ai_failed') : t('chn_ai_fallback');
        showToast(msg, res.ai ? 'success' : '');
        return loadList().then(() => selectDraft(res.draft_id));
      }).catch((e) => { err.textContent = (e.errors && e.errors[0]) || e.message; })
        .finally(() => { make.disabled = false; make.textContent = old; });
    };
  }

  // The old template grid, now one of the choices inside the new window.
  function templateChooser(ed) {
    Array.from(ed.children).forEach((c) => { if (!c.classList.contains('chn-editor-head')) c.remove(); });
    const body = el('div', 'chn-tpls');
    ed.appendChild(body);
    Promise.all([
      api('GET', '/api/chains/templates').catch(() => []),
      api('GET', '/api/connections').catch(() => []),
      api('GET', '/api/storages').catch(() => []),
    ]).then(([tpls, conns, storages]) => {
      body.appendChild(el('h4', '', t('chn_tpl_pick')));
      const grid = el('div', 'chn-tpl-grid');
      body.appendChild(grid);
      const card = (name, desc, extra, fn) => {
        const b = el('button', 'chn-tpl');
        b.type = 'button';
        b.appendChild(el('b', '', name));
        b.appendChild(el('span', '', desc));
        if (extra) { const x = el('span', '', ' · ' + extra); b.appendChild(x); }
        b.onclick = fn;
        grid.appendChild(b);
      };
      (tpls || []).forEach((tp) => {
        const k = 'chn_tpl_' + tp.id;
        card(LANGS.en[k] ? t(k) : tp.name, LANGS.en[k + '_d'] ? t(k + '_d') : tp.description,
          fmt(t('chn_nodes_n'), { n: tp.definition.nodes.length }), () => templateForm(ed, body, tp, conns || [], storages || []));
      });
      card(t('chn_tpl_empty'), t('chn_tpl_empty_d'), '', () => { ed.remove(); openEditor(null, { definition: defaultDef(), enabled: true }); });
    });
  }


  function templateForm(ed, body, tp, conns, storages) {
    body.textContent = '';
    const ssh = conns.filter((c) => (c.protocol || 'SSH').toUpperCase() === 'SSH');
    const def = JSON.parse(JSON.stringify(tp.definition));
    const form = el('div', 'chn-form');
    body.appendChild(form);
    const k = 'chn_tpl_' + tp.id;
    form.appendChild(el('h4', '', LANGS.en[k] ? t(k) : tp.name));
    const nameL = el('label', '', t('chn_tpl_name'));
    const name = el('input');
    name.value = def.name;
    nameL.appendChild(name);
    form.appendChild(nameL);
    const pick = (kind, items, labelOf) => {
      const nodes = def.nodes.filter((n) => n.kind === kind);
      if (!nodes.length) return [];
      form.appendChild(el('h5', '', t(kind === 'server' ? 'chn_tpl_servers' : 'chn_tpl_storages')));
      if (!items.length) form.appendChild(el('div', 'note', t(kind === 'server' ? 'chn_tpl_no_ssh' : 'chn_tpl_no_storage')));
      return nodes.map((n) => {
        const l = el('label', '', [n.label || n.id, n.role].filter(Boolean).join(' · '));
        const s = el('select');
        s.appendChild(el('option', '', t('chn_tpl_choose'))).value = '';
        items.forEach((it) => { const o = el('option', '', labelOf(it)); o.value = String(it.id); s.appendChild(o); });
        l.appendChild(s);
        form.appendChild(l);
        return { n, s };
      });
    };
    const picks = pick('server', ssh, (c) => c.name + (c.host ? ' (' + c.host + ')' : ''))
      .concat(pick('storage', storages, (s) => s.name + (s.bucket ? ' · ' + s.bucket : '')));
    form.appendChild(el('div', 'note', t('chn_tpl_after')));
    const errs = el('div', 'chn-errs');
    errs.setAttribute('role', 'status');
    form.appendChild(errs);
    const act = el('div', 'chn-row');
    const back = el('button', 'btn-sec btn-sm', t('chn_back'));
    back.onclick = () => { ed.remove(); openNew(); };
    const json = el('button', 'btn-sec btn-sm', '{ } ' + t('chn_edit'));
    const create = el('button', 'btn-pri btn-sm', t('chn_create'));
    act.appendChild(back); act.appendChild(json); act.appendChild(create);
    form.appendChild(act);
    const build = (strict) => {
      let missingServer = false;
      picks.forEach(({ n, s }) => {
        const v = Number(s.value) || 0;
        if (!v && n.kind === 'server') missingServer = true;
        if (n.kind === 'server') n.conn_id = v; else n.storage_id = v;
      });
      def.name = name.value.trim() || def.name;
      if (strict) {
        // A storage is optional: leave an unpicked storage node out instead of blocking.
        picks.filter(({ n, s }) => n.kind === 'storage' && !Number(s.value)).forEach(({ n }) => stripNode(def, n.id));
      }
      if (strict && missingServer) { errs.textContent = t('chn_tpl_missing'); return null; }
      return { definition: def, enabled: true };
    };
    json.onclick = () => { const v = build(false); if (v) openEditor(null, v); };
    create.onclick = () => {
      const v = build(true);
      if (!v) return;
      create.disabled = true;
      api('POST', '/api/chains', v).then((c) => {
        ed.remove();
        showToast(t('chn_saved'), 'success');
        state.current = c;
        state.sel = null;
        return loadList();
      }).catch((e) => {
        errs.textContent = '';
        (e.errors.length ? e.errors : [e.message]).forEach((x) => errs.appendChild(el('div', '', '• ' + x)));
      }).finally(() => { create.disabled = false; });
    };
  }

  function importFile() {
    const inp = document.createElement('input');
    inp.type = 'file'; inp.accept = '.json,application/json';
    inp.onchange = () => {
      const f = inp.files && inp.files[0];
      if (!f) return;
      f.text().then((txt) => {
        let data = null;
        try { data = JSON.parse(txt); } catch (_) { data = null; }
        if (!data || typeof data !== 'object') { showToast(t('chn_import_bad'), 'error'); return; }
        if (!data.definition && data.nodes) data = { definition: data, enabled: true };
        if (!data.definition) { showToast(t('chn_import_bad'), 'error'); return; }
        openEditor(null, data);
      });
    };
    inp.click();
  }

  function openEditor(id, value) {
    const ed = editorPane(t('chn_editor_t'));
    if (!ed) return;
    ed.querySelector('.chn-editor-head').insertAdjacentHTML('beforeend', `
        <label style="font-size:11.5px;display:flex;gap:6px;align-items:center;"><input type="checkbox" id="chn-ed-en"> ${esc(t('chn_enabled'))}</label>
        <button class="btn-sec btn-sm" id="chn-ed-val">${esc(t('chn_validate'))}</button>
        <button class="btn-sec btn-sm" id="chn-ed-cancel">${esc(t('chn_cancel'))}</button>
        <button class="btn-pri btn-sm" id="chn-ed-save">${esc(t('chn_save'))}</button>`);
    ed.insertAdjacentHTML('beforeend', `<div class="chn-editor-body"><textarea id="chn-ed-text" spellcheck="false" aria-label="${esc(t('chn_editor_t'))}"></textarea>
        <div class="chn-help" id="chn-ed-help"><b>${esc(t('chn_help'))}</b></div></div>
      <div class="chn-errs" id="chn-ed-errs" role="status" aria-live="polite"></div>`);
    ed.querySelector('#chn-ed-text').value = JSON.stringify(value.definition, null, 2);
    ed.querySelector('#chn-ed-en').checked = value.enabled !== false;
    ed.querySelector('#chn-ed-cancel').onclick = () => ed.remove();
    ed.querySelector('#chn-ed-val').onclick = () => submit(ed, id, true);
    ed.querySelector('#chn-ed-save').onclick = () => submit(ed, id, false);
    fillHelp(ed.querySelector('#chn-ed-help'));
  }

  function showErrors(ed, list, okText) {
    const box = ed && ed.querySelector('#chn-ed-errs');
    if (!box) { showToast((list && list[0]) || okText || t('chn_failed'), 'error'); return; }
    box.textContent = '';
    if (okText) { box.style.color = '#3fb950'; box.appendChild(el('div', '', okText)); return; }
    box.style.color = '';
    list.forEach((e) => box.appendChild(el('div', '', '• ' + e)));
  }

  function submit(ed, id, validateOnly) {
    let def = null;
    try { def = JSON.parse(ed.querySelector('#chn-ed-text').value); } catch (e) { showErrors(ed, [String(e.message || e)]); return; }
    const body = { definition: def, enabled: ed.querySelector('#chn-ed-en').checked };
    if (validateOnly) {
      api('POST', '/api/chains/validate', body).then((r) => {
        if (r.ok) showErrors(ed, [], t('chn_valid')); else showErrors(ed, r.errors || []);
      }).catch((e) => showErrors(ed, e.errors.length ? e.errors : [e.message]));
      return;
    }
    const req = id ? api('PUT', `/api/chains/${id}`, body) : api('POST', '/api/chains', body);
    req.then((c) => {
      ed.remove();
      showToast(t('chn_saved'), 'success');
      state.current = c;
      state.sel = null;
      return loadList();
    }).catch((e) => showErrors(ed, e.errors.length ? e.errors : [e.message]));
  }

  function fillHelp(box) {
    const add = (title, lines) => {
      box.appendChild(el('h5', '', title));
      lines.forEach((l) => { const d = el('div'); const c = el('code', '', l[0]); d.appendChild(c); if (l[1]) d.appendChild(document.createTextNode(' ' + l[1])); box.appendChild(d); });
    };
    Promise.all([
      api('GET', '/api/connections').catch(() => []),
      api('GET', '/api/storages').catch(() => []),
      api('GET', '/api/chains/catalog').catch(() => ({ types: {} })),
    ]).then(([conns, storages, cat]) => {
      const ssh = (conns || []).filter((c) => (c.protocol || 'SSH').toUpperCase() === 'SSH');
      add(t('chn_help_conns'), ssh.slice(0, 200).map((c) => [String(c.id), c.name]));
      if ((storages || []).length) add(t('chn_help_storages'), storages.map((s) => [String(s.id), s.name]));
      add(t('chn_help_types'), Object.keys(cat.types || {}).sort().map((k) => [k, LANGS.en['chn_type_' + k] ? t('chn_type_' + k) : cat.types[k].help]));
      box.appendChild(el('h5', '', t('chn_scenarios')));
      box.appendChild(el('div', '', t('chn_help_rules')));
      box.appendChild(el('h5', '', t('chn_links')));
      box.appendChild(el('div', '', t('chn_help_links')));
    });
  }

  // A poll finished somewhere: refresh the list, and the details when that chain is open.
  function changed(id) {
    if (!isOpen()) return;
    if ($('.chn-editor')) { loadListOnly(); return; }
    if (state.current && !state.current.draft && state.current.id === id) loadList(); else loadListOnly();
  }

  // The assistant created or replaced a draft: refresh the list and keep the open draft.
  function changedDraft(id) {
    if (!isOpen()) return;
    if ($('.chn-editor')) { loadListOnly(); return; }
    return loadList();
  }

  function loadListOnly() {
    Promise.all([api('GET', '/api/chains'), api('GET', '/api/chains/drafts').catch(() => [])])
      .then(([list, drafts]) => { state.list = list || []; state.drafts = drafts || []; renderList(); }).catch(() => {});
  }

  window.wrmChains = { open, close, changed, changedDraft };
})();
