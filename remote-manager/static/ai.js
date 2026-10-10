// AI assistant of WRM PRO: the panel next to a terminal or file manager window, personal
// API keys and the administrators' AI tab. Loaded on demand by openAI() / the admin panel
// in index.html; uses its helpers (t, esc, showToast, apiError, uiConfirm, fmtTime,
// winObjects, winConn, winCan, currentUser, fitTerm, openReplay, settingRowHTML).
// The server decides everything (modes, approvals, limits); this file only shows it.
'use strict';

Object.assign(LANGS.en, {
  ai_title: 'AI assistant', ai_start: 'Start', ai_provider: 'Provider', ai_model: 'Model', ai_mode: 'Mode',
  ai_mode_read_only: 'Read-only', ai_mode_ask: 'Ask before every change', ai_mode_auto: 'Automatic within limits',
  ai_mode_read_only_d: 'The assistant reads logs, files, services and packages. Every change is refused.',
  ai_mode_ask_d: 'Reads run directly. Every command that changes something and every file edit waits for your approval.',
  ai_mode_auto_d: 'Changes that match your allow list run without asking until the time or action limit. Everything else still asks. Destructive commands are always refused.',
  ai_disabled: 'The AI assistant is not available', ai_no_providers: 'No AI provider is configured. An administrator adds the organisation\'s provider in Admin panel → AI assistant.',
  ai_add_own_key: 'Add my own API key', ai_my_keys: 'My API keys', ai_share_notes: 'Share my notes for this connection with the assistant',
  ai_auto_allow: 'Allowed commands (one pattern per line, * = anything)', ai_auto_deny: 'Denied commands', ai_auto_path_allow: 'Files it may change (/etc/nginx/**)',
  ai_auto_path_deny: 'Files it must not change', ai_auto_minutes: 'Time limit (minutes)', ai_auto_actions: 'Maximum actions',
  ai_auto_confirm: 'I allow the assistant to make these changes without asking, within these limits.', ai_auto_need_confirm: 'Confirm the automatic mode and its limits first.',
  ai_see: 'What can the AI see?', ai_see_tools: 'Read tools', ai_see_cmd: 'Commands: read-only ones always (in every mode); changes only as the mode allows',
  ai_see_write: 'File edits (shown as a diff)', ai_see_redacted: 'Secrets in output are replaced by [REDACTED] before they are sent to the provider.',
  ai_see_not_redacted: 'Output is sent to the provider as it is (redaction is off by policy).', ai_see_limit: 'Output and files are limited to {kb} KB per call.',
  ai_see_never: 'Never: passwords and keys stored in WRM, your WRM sign-in, other connections, your terminal screen.', ai_see_denied: '{n} kinds of sensitive files (password hashes, private keys, credentials) are never read.',
  ai_see_provider: 'Everything the assistant sees is sent to {provider}.', ai_see_notes_on: 'Your notes for this connection are shared.',
  ai_placeholder: 'Ask about this server…', ai_send: 'Send', ai_stop: 'Stop', ai_end: 'End the AI session', ai_end_q: 'End this AI session? A running command is stopped.',
  ai_close_panel: 'Hide the panel (the session keeps running)', ai_history: 'My AI sessions', ai_continue: 'Continue session', ai_new: 'New session',
  ai_approve: 'Approve', ai_deny: 'Deny', ai_edit: 'Edit', ai_note_ph: 'Note for the assistant (optional)', ai_expires: 'expires in {s} s',
  ai_wait_approval: 'Waiting for your approval', ai_d_allow: 'ran', ai_d_readonly: 'read-only', ai_d_auto: 'automatic', ai_d_approve: 'needs approval', ai_d_deny: 'denied',
  ai_o_approved: 'approved', ai_o_denied: 'denied', ai_o_timeout: 'timed out', ai_o_cancelled: 'cancelled', ai_edited: 'edited',
  ai_exit: 'exit {code}', ai_redactions: '{n} redacted', ai_output: 'Output', ai_reason: 'Why',
  ai_n_stopped: 'Stopped.', ai_n_refusal: 'The model declined this request.', ai_n_step_limit: 'The assistant stopped after the maximum number of steps.',
  ai_ended: 'The AI session has ended', ai_mode_changed: 'Mode: {mode}', ai_by: 'by {who}', ai_usage: '{in} in · {out} out tokens',
  ai_r_user: 'ended by you', ai_r_admin: 'stopped by an administrator', ai_r_killall: 'every AI session was stopped by an administrator',
  ai_r_switch: 'the kill switch was turned on', ai_r_off: 'the AI assistant was turned off', ai_r_blocked: 'an administrator turned the assistant off for you',
  ai_r_idle: 'ended after inactivity', ai_r_policy: 'the policy no longer allows the assistant on this connection',
  ai_fm_ctx: 'Directory open in the file manager', ai_active_on: 'AI active ({mode})',
  ai_auto_until: 'until {time} · {used}/{max} actions', ai_mode_auto_setup: 'Automatic mode limits',
  // personal / org providers
  ai_prov_kind: 'Type', ai_prov_name: 'Name', ai_prov_models: 'Models (comma-separated)', ai_prov_default: 'Default model',
  ai_prov_key: 'API key', ai_prov_key_keep: 'unchanged (stored encrypted)', ai_prov_base: 'Base URL (optional: gateway, private endpoint)',
  ai_prov_base_req: 'Base URL', ai_prov_endpoint: 'Endpoint (https://<resource>.openai.azure.com)', ai_prov_deployment: 'Deployment (empty: the model name)',
  ai_prov_apiver: 'API version (empty: the v1 API)', ai_prov_region: 'Region', ai_prov_project: 'Google Cloud project', ai_prov_api: 'API',
  ai_prov_ak: 'AWS access key ID', ai_prov_sk: 'AWS secret access key', ai_prov_st: 'AWS session token (optional)', ai_prov_bkey: 'Bedrock API key (instead of access keys)',
  ai_prov_sa: 'Service account key (JSON)', ai_prov_price: 'Price per million tokens (USD, input / output; empty: built-in for Anthropic)',
  ai_prov_maxtok: 'Max output tokens (0 = 16000)', ai_prov_fallback: 'Refusal fallback (Anthropic API)', ai_prov_enabled: 'Enabled', ai_prov_test: 'Test',
  ai_prov_test_ok: 'The provider answered', ai_prov_add: 'Add provider', ai_prov_edit: 'Edit provider', ai_prov_del_q: 'Delete this provider? Its sessions end.',
  ai_prov_none: 'No providers yet.', ai_models_col: 'Models', ai_prov_secrets: 'Keys', ai_kind_anthropic: 'Anthropic API', ai_kind_openai: 'OpenAI API', ai_kind_azure: 'Azure OpenAI',
  ai_kind_bedrock: 'AWS Bedrock', ai_kind_vertex: 'Google Vertex AI', ai_kind_openai_compatible: 'OpenAI-compatible (local models)',
  ai_subscriptions_note: 'Consumer Claude.ai and ChatGPT chat subscriptions cannot be used by other applications. Use an API key or an enterprise gateway here, or connect the chat app to WRM over MCP (Settings → AI connections).',
  // admin
  adm_ai_kill: 'Kill switch', ai_kill_all: 'Stop every AI session now', ai_kill_all_q: 'End every AI session of every user now?', ai_killed_n: '{n} sessions ended',
  ai_kill_switch_on: 'The kill switch is ON: nobody can use the AI assistant.', ai_kill_switch_off: 'The kill switch is off.',
  ai_policies: 'Policies', ai_providers_org: 'Organisation providers', ai_sessions: 'AI sessions', ai_usage_users: 'Usage per user', ai_destructive: 'Always blocked (destructive patterns)',
  ai_transcript: 'Transcript', ai_replay: 'Replay', ai_kill: 'Kill', ai_block: 'Turn off for this user', ai_unblock: 'Turn on again', ai_blocked: 'off',
  ai_live: 'live', ai_cost: 'Cost', ai_tokens: 'Tokens', ai_requests: 'Requests', ai_tool_calls: 'Tool calls', ai_mode_rules_help: 'Mode rules (most restrictive wins): JSON, e.g. [{"match":"tag","value":"prod","modes":["read_only"]}, {"match":"folder","value":"Production","modes":["read_only","ask"]}]. match: tag, folder, connection (id), host (wildcards).',
  // policies (Admin panel → AI assistant)
  p_ai_assistant: 'AI assistant', p_ai_assistant_d: 'Who may use the built-in AI assistant.', po_ai_assistant_off: 'Off', po_ai_assistant_admins: 'Administrators', po_ai_assistant_all: 'Everybody',
  p_ai_kill_switch: 'Kill switch', p_ai_kill_switch_d: 'Ends every AI session at once and blocks new ones until it is turned off.',
  p_ai_modes: 'Allowed modes', p_ai_modes_d: 'Global default, comma-separated: read_only, ask, auto. Mode rules can only narrow it.', p_ai_modes_ph: 'read_only,ask',
  p_ai_default_mode: 'Default mode', p_ai_default_mode_d: 'Mode of a new session when it is allowed.', po_ai_default_mode_read_only: 'Read-only', po_ai_default_mode_ask: 'Ask before every change',
  p_ai_mode_rules: 'Mode rules', p_ai_mode_rules_d: 'Per tag, folder, connection or host; the most restrictive matching rule wins.',
  p_ai_personal_keys: 'Personal API keys', p_ai_personal_keys_d: 'Who may add their own provider keys.', po_ai_personal_keys_off: 'Nobody', po_ai_personal_keys_admins: 'Administrators', po_ai_personal_keys_all: 'Everybody',
  p_ai_provider_kinds: 'Allowed provider types', p_ai_provider_kinds_d: 'anthropic, openai, azure, bedrock, vertex, openai_compatible', p_ai_provider_kinds_ph: 'anthropic,azure',
  p_ai_models: 'Allowed models', p_ai_models_d: 'Comma-separated patterns (claude-*); empty = every model of the providers.', p_ai_models_ph: 'claude-*',
  p_ai_redact_output: 'Redact output', p_ai_redact_output_d: 'Replace secrets (passwords, tokens, private keys …) before output is sent to the provider. Transcripts and the audit are always redacted.',
  p_ai_share_notes: 'Allow sharing notes', p_ai_share_notes_d: 'Users may opt in to send their connection notes to the assistant.',
  p_ai_transcript_retention_days: 'Transcript retention (days)', p_ai_transcript_retention_days_d: 'AI transcripts are deleted after this time (recordings follow the recording retention).',
  p_ai_approval_timeout_seconds: 'Approval timeout (seconds)', p_ai_approval_timeout_seconds_d: 'An approval that is not answered in time is denied.',
  p_ai_auto_max_minutes: 'Automatic mode: max minutes', p_ai_auto_max_minutes_d: 'The longest time limit a user may choose.',
  p_ai_auto_max_actions: 'Automatic mode: max actions', p_ai_auto_max_actions_d: 'The largest number of automatic actions per opt-in.',
  p_ai_allow_power: 'Shutdown and reboot', p_ai_allow_power_d: 'Blocked, or allowed with an approval of each command (never automatic).', po_ai_allow_power_off: 'Always blocked', po_ai_allow_power_ask: 'With an approval',
  p_ai_blocked_commands: 'Extra blocked commands', p_ai_blocked_commands_d: 'Comma-separated patterns refused in every mode, in addition to the built-in list.', p_ai_blocked_commands_ph: 'docker system prune*',
  p_ai_read_deny_paths: 'Extra unreadable files', p_ai_read_deny_paths_d: 'Comma-separated path patterns the assistant never reads (a pattern without / matches the file name).', p_ai_read_deny_paths_ph: '/srv/app/.env, *.p12',
  p_ai_command_timeout_seconds: 'Command time limit (seconds)', p_ai_command_timeout_seconds_d: 'A tool call is stopped after this time.',
  p_ai_output_max_kb: 'Output limit (KB)', p_ai_output_max_kb_d: 'Per tool call, also the size limit for reading files.',
  p_ai_max_steps: 'Steps per request', p_ai_max_steps_d: 'Model requests per user message.',
  p_ai_idle_minutes: 'Idle timeout (minutes)', p_ai_idle_minutes_d: 'An idle AI session ends after this time.',
});
Object.assign(LANGS.hr, {
  ai_title: 'AI asistent', ai_start: 'Pokreni', ai_provider: 'Pružatelj', ai_model: 'Model', ai_mode: 'Način rada',
  ai_mode_read_only: 'Samo čitanje', ai_mode_ask: 'Pitaj prije svake promjene', ai_mode_auto: 'Automatski u granicama',
  ai_mode_read_only_d: 'Asistent čita logove, datoteke, servise i pakete. Svaka promjena se odbija.',
  ai_mode_ask_d: 'Čitanje se izvršava odmah. Svaka naredba koja nešto mijenja i svaka izmjena datoteke čeka tvoje odobrenje.',
  ai_mode_auto_d: 'Promjene s tvoje liste dopuštenih izvršavaju se bez pitanja do isteka vremena ili broja akcija. Za sve ostalo i dalje pita. Destruktivne naredbe su uvijek odbijene.',
  ai_disabled: 'AI asistent nije dostupan', ai_no_providers: 'Nijedan AI pružatelj nije podešen. Administrator dodaje pružatelja organizacije u Admin panel → AI asistent.',
  ai_add_own_key: 'Dodaj vlastiti API ključ', ai_my_keys: 'Moji API ključevi', ai_share_notes: 'Podijeli s asistentom moje bilješke za ovu vezu',
  ai_auto_allow: 'Dopuštene naredbe (jedan uzorak po retku, * = bilo što)', ai_auto_deny: 'Zabranjene naredbe', ai_auto_path_allow: 'Datoteke koje smije mijenjati (/etc/nginx/**)',
  ai_auto_path_deny: 'Datoteke koje ne smije mijenjati', ai_auto_minutes: 'Vremensko ograničenje (minute)', ai_auto_actions: 'Najviše akcija',
  ai_auto_confirm: 'Dopuštam asistentu da ove promjene radi bez pitanja, unutar ovih granica.', ai_auto_need_confirm: 'Najprije potvrdi automatski način i njegove granice.',
  ai_see: 'Što AI vidi?', ai_see_tools: 'Alati za čitanje', ai_see_cmd: 'Naredbe: one samo za čitanje uvijek (u svakom načinu); promjene samo koliko način dopušta',
  ai_see_write: 'Izmjene datoteka (prikazane kao diff)', ai_see_redacted: 'Tajne u izlazu zamjenjuju se s [REDACTED] prije slanja pružatelju.',
  ai_see_not_redacted: 'Izlaz se pružatelju šalje kakav jest (redakcija je isključena politikom).', ai_see_limit: 'Izlaz i datoteke ograničeni su na {kb} KB po pozivu.',
  ai_see_never: 'Nikada: lozinke i ključevi spremljeni u WRM, tvoja WRM prijava, druge veze, tvoj zaslon terminala.', ai_see_denied: '{n} vrsta osjetljivih datoteka (hashevi lozinki, privatni ključevi, vjerodajnice) nikada se ne čita.',
  ai_see_provider: 'Sve što asistent vidi šalje se pružatelju {provider}.', ai_see_notes_on: 'Tvoje bilješke za ovu vezu su podijeljene.',
  ai_placeholder: 'Pitaj o ovom serveru…', ai_send: 'Pošalji', ai_stop: 'Zaustavi', ai_end: 'Završi AI sesiju', ai_end_q: 'Završiti ovu AI sesiju? Naredba koja se izvršava bit će zaustavljena.',
  ai_close_panel: 'Sakrij panel (sesija nastavlja raditi)', ai_history: 'Moje AI sesije', ai_continue: 'Nastavi sesiju', ai_new: 'Nova sesija',
  ai_approve: 'Odobri', ai_deny: 'Odbij', ai_edit: 'Uredi', ai_note_ph: 'Napomena za asistenta (neobavezno)', ai_expires: 'istječe za {s} s',
  ai_wait_approval: 'Čeka tvoje odobrenje', ai_d_allow: 'izvršeno', ai_d_readonly: 'samo čitanje', ai_d_auto: 'automatski', ai_d_approve: 'treba odobrenje', ai_d_deny: 'odbijeno',
  ai_o_approved: 'odobreno', ai_o_denied: 'odbijeno', ai_o_timeout: 'isteklo', ai_o_cancelled: 'otkazano', ai_edited: 'uređeno',
  ai_exit: 'izlaz {code}', ai_redactions: 'redigirano: {n}', ai_output: 'Izlaz', ai_reason: 'Zašto',
  ai_n_stopped: 'Zaustavljeno.', ai_n_refusal: 'Model je odbio ovaj zahtjev.', ai_n_step_limit: 'Asistent je stao nakon najvećeg broja koraka.',
  ai_ended: 'AI sesija je završena', ai_mode_changed: 'Način rada: {mode}', ai_by: 'promijenio: {who}', ai_usage: '{in} ulaznih · {out} izlaznih tokena',
  ai_r_user: 'završio si je', ai_r_admin: 'zaustavio ju je administrator', ai_r_killall: 'administrator je zaustavio sve AI sesije',
  ai_r_switch: 'uključen je prekidač za hitno zaustavljanje', ai_r_off: 'AI asistent je isključen', ai_r_blocked: 'administrator ti je isključio asistenta',
  ai_r_idle: 'završena zbog neaktivnosti', ai_r_policy: 'politika više ne dopušta asistenta na ovoj vezi',
  ai_fm_ctx: 'Direktorij otvoren u pregledniku datoteka', ai_active_on: 'AI aktivan ({mode})',
  ai_auto_until: 'do {time} · {used}/{max} akcija', ai_mode_auto_setup: 'Granice automatskog načina',
  ai_prov_kind: 'Vrsta', ai_prov_name: 'Naziv', ai_prov_models: 'Modeli (odvojeni zarezom)', ai_prov_default: 'Zadani model',
  ai_prov_key: 'API ključ', ai_prov_key_keep: 'nepromijenjen (spremljen šifrirano)', ai_prov_base: 'Osnovni URL (neobavezno: gateway, privatni endpoint)',
  ai_prov_base_req: 'Osnovni URL', ai_prov_endpoint: 'Endpoint (https://<resurs>.openai.azure.com)', ai_prov_deployment: 'Deployment (prazno: naziv modela)',
  ai_prov_apiver: 'Verzija API-ja (prazno: v1 API)', ai_prov_region: 'Regija', ai_prov_project: 'Google Cloud projekt', ai_prov_api: 'API',
  ai_prov_ak: 'AWS access key ID', ai_prov_sk: 'AWS secret access key', ai_prov_st: 'AWS session token (neobavezno)', ai_prov_bkey: 'Bedrock API ključ (umjesto access ključeva)',
  ai_prov_sa: 'Ključ servisnog računa (JSON)', ai_prov_price: 'Cijena po milijun tokena (USD, ulaz / izlaz; prazno: ugrađena za Anthropic)',
  ai_prov_maxtok: 'Najviše izlaznih tokena (0 = 16000)', ai_prov_fallback: 'Zamjenski model kod odbijanja (Anthropic API)', ai_prov_enabled: 'Uključen', ai_prov_test: 'Testiraj',
  ai_prov_test_ok: 'Pružatelj je odgovorio', ai_prov_add: 'Dodaj pružatelja', ai_prov_edit: 'Uredi pružatelja', ai_prov_del_q: 'Obrisati ovog pružatelja? Njegove sesije završavaju.',
  ai_prov_none: 'Još nema pružatelja.', ai_models_col: 'Modeli', ai_prov_secrets: 'Ključevi', ai_kind_anthropic: 'Anthropic API', ai_kind_openai: 'OpenAI API', ai_kind_azure: 'Azure OpenAI',
  ai_kind_bedrock: 'AWS Bedrock', ai_kind_vertex: 'Google Vertex AI', ai_kind_openai_compatible: 'OpenAI-kompatibilan (lokalni modeli)',
  ai_subscriptions_note: 'Potrošačke pretplate na Claude.ai i ChatGPT chat ne mogu koristiti druge aplikacije. Ovdje koristi API ključ ili enterprise gateway, ili poveži chat aplikaciju s WRM-om preko MCP-a (Postavke → AI veze).',
  adm_ai_kill: 'Prekidač za hitno zaustavljanje', ai_kill_all: 'Odmah zaustavi sve AI sesije', ai_kill_all_q: 'Odmah završiti sve AI sesije svih korisnika?', ai_killed_n: 'Završeno sesija: {n}',
  ai_kill_switch_on: 'Prekidač je UKLJUČEN: nitko ne može koristiti AI asistenta.', ai_kill_switch_off: 'Prekidač je isključen.',
  ai_policies: 'Politike', ai_providers_org: 'Pružatelji organizacije', ai_sessions: 'AI sesije', ai_usage_users: 'Potrošnja po korisniku', ai_destructive: 'Uvijek blokirano (destruktivni uzorci)',
  ai_transcript: 'Transkript', ai_replay: 'Reprodukcija', ai_kill: 'Prekini', ai_block: 'Isključi za ovog korisnika', ai_unblock: 'Ponovno uključi', ai_blocked: 'isključen',
  ai_live: 'aktivna', ai_cost: 'Trošak', ai_tokens: 'Tokeni', ai_requests: 'Zahtjevi', ai_tool_calls: 'Pozivi alata', ai_mode_rules_help: 'Pravila načina rada (najstroža pobjeđuju): JSON, npr. [{"match":"tag","value":"prod","modes":["read_only"]}, {"match":"folder","value":"Production","modes":["read_only","ask"]}]. match: tag, folder, connection (id), host (zamjenski znakovi).',
  p_ai_assistant: 'AI asistent', p_ai_assistant_d: 'Tko smije koristiti ugrađenog AI asistenta.', po_ai_assistant_off: 'Isključen', po_ai_assistant_admins: 'Administratori', po_ai_assistant_all: 'Svi',
  p_ai_kill_switch: 'Prekidač za hitno zaustavljanje', p_ai_kill_switch_d: 'Odmah završava sve AI sesije i blokira nove dok se ne isključi.',
  p_ai_modes: 'Dopušteni načini rada', p_ai_modes_d: 'Globalna zadana vrijednost, odvojeno zarezom: read_only, ask, auto. Pravila je mogu samo suziti.', p_ai_modes_ph: 'read_only,ask',
  p_ai_default_mode: 'Zadani način rada', p_ai_default_mode_d: 'Način rada nove sesije kad je dopušten.', po_ai_default_mode_read_only: 'Samo čitanje', po_ai_default_mode_ask: 'Pitaj prije svake promjene',
  p_ai_mode_rules: 'Pravila načina rada', p_ai_mode_rules_d: 'Po oznaci, mapi, vezi ili hostu; pobjeđuje najstrože pravilo koje odgovara.',
  p_ai_personal_keys: 'Osobni API ključevi', p_ai_personal_keys_d: 'Tko smije dodati vlastite ključeve pružatelja.', po_ai_personal_keys_off: 'Nitko', po_ai_personal_keys_admins: 'Administratori', po_ai_personal_keys_all: 'Svi',
  p_ai_provider_kinds: 'Dopuštene vrste pružatelja', p_ai_provider_kinds_d: 'anthropic, openai, azure, bedrock, vertex, openai_compatible', p_ai_provider_kinds_ph: 'anthropic,azure',
  p_ai_models: 'Dopušteni modeli', p_ai_models_d: 'Uzorci odvojeni zarezom (claude-*); prazno = svi modeli pružatelja.', p_ai_models_ph: 'claude-*',
  p_ai_redact_output: 'Rediguj izlaz', p_ai_redact_output_d: 'Zamijeni tajne (lozinke, tokene, privatne ključeve …) prije slanja izlaza pružatelju. Transkripti i audit uvijek su redigirani.',
  p_ai_share_notes: 'Dopusti dijeljenje bilješki', p_ai_share_notes_d: 'Korisnici mogu izabrati da asistentu pošalju bilješke veze.',
  p_ai_transcript_retention_days: 'Čuvanje transkripata (dani)', p_ai_transcript_retention_days_d: 'AI transkripti brišu se nakon tog vremena (snimke prate čuvanje snimki).',
  p_ai_approval_timeout_seconds: 'Rok za odobrenje (sekunde)', p_ai_approval_timeout_seconds_d: 'Odobrenje na koje se ne odgovori na vrijeme se odbija.',
  p_ai_auto_max_minutes: 'Automatski način: najviše minuta', p_ai_auto_max_minutes_d: 'Najdulje vremensko ograničenje koje korisnik smije izabrati.',
  p_ai_auto_max_actions: 'Automatski način: najviše akcija', p_ai_auto_max_actions_d: 'Najveći broj automatskih akcija po uključivanju.',
  p_ai_allow_power: 'Gašenje i ponovno pokretanje', p_ai_allow_power_d: 'Blokirano ili dopušteno uz odobrenje svake naredbe (nikad automatski).', po_ai_allow_power_off: 'Uvijek blokirano', po_ai_allow_power_ask: 'Uz odobrenje',
  p_ai_blocked_commands: 'Dodatno blokirane naredbe', p_ai_blocked_commands_d: 'Uzorci odvojeni zarezom, odbijeni u svakom načinu, uz ugrađeni popis.', p_ai_blocked_commands_ph: 'docker system prune*',
  p_ai_read_deny_paths: 'Dodatne nečitljive datoteke', p_ai_read_deny_paths_d: 'Uzorci putanja odvojeni zarezom koje asistent nikad ne čita (uzorak bez / odgovara nazivu datoteke).', p_ai_read_deny_paths_ph: '/srv/app/.env, *.p12',
  p_ai_command_timeout_seconds: 'Vremensko ograničenje naredbe (sekunde)', p_ai_command_timeout_seconds_d: 'Poziv alata zaustavlja se nakon tog vremena.',
  p_ai_output_max_kb: 'Ograničenje izlaza (KB)', p_ai_output_max_kb_d: 'Po pozivu alata, ujedno ograničenje veličine za čitanje datoteka.',
  p_ai_max_steps: 'Koraka po zahtjevu', p_ai_max_steps_d: 'Zahtjeva prema modelu po poruci korisnika.',
  p_ai_idle_minutes: 'Neaktivnost (minute)', p_ai_idle_minutes_d: 'Neaktivna AI sesija završava nakon tog vremena.',
});

(function () {
  const $ = id => document.getElementById(id);
  const S = {}; // per window: {sid, es, seq, view, cfg, pol, busy}
  let cfg = null;
  const MODES = ['read_only', 'ask', 'auto'];

  const css = `
.ai-panel { position: absolute; right: 0; bottom: 0; width: 380px; max-width: 100%; display: flex; flex-direction: column; background: var(--bg2); border-left: 1px solid var(--border2); z-index: 5; border-radius: 0 0 11px 0; min-width: 0; }
.win.ai-open .win-content { margin-right: 380px; }
@container (max-width: 760px) { .win.ai-open .ai-panel { width: 100%; border-left: none; border-radius: 0 0 11px 11px; } .win.ai-open .win-content { margin-right: 0; } }
.ai-panel input[type=text], .ai-panel input[type=number], .ai-panel select, .ai-panel textarea, .ai-modal input[type=text], .ai-modal input[type=password], .ai-modal input[type=number], .ai-modal select, .ai-modal textarea {
  background: var(--bg); border: 1px solid var(--border2); color: var(--text); border-radius: 6px; padding: 6px 8px; font-size: 12px; outline: none; font-family: inherit; user-select: text; }
.ai-panel input:focus, .ai-panel select:focus, .ai-panel textarea:focus, .ai-modal input:focus, .ai-modal select:focus, .ai-modal textarea:focus { border-color: var(--accent); }
.ai-panel select, .ai-modal select { background: var(--bg3); cursor: pointer; }
.ai-head { display: flex; align-items: center; gap: 6px; padding: 6px 8px; border-bottom: 1px solid var(--border); background: var(--bg3); min-width: 0; }
.ai-head .ai-ttl { font-weight: 700; font-size: 12px; white-space: nowrap; }
.ai-head .ai-sp { flex: 1; }
.ai-ib { background: transparent; border: 1px solid transparent; color: var(--text2); border-radius: 6px; padding: 2px 6px; font-size: 12px; line-height: 18px; }
.ai-ib:hover { background: var(--bg4); color: var(--text); }
.ai-ib.danger:hover { color: var(--red); border-color: rgba(239,68,68,.35); }
.ai-body { flex: 1; overflow-y: auto; padding: 8px; display: flex; flex-direction: column; gap: 8px; min-height: 0; user-select: text; }
.ai-foot { border-top: 1px solid var(--border); padding: 6px; display: flex; gap: 6px; align-items: flex-end; }
.ai-foot textarea { flex: 1; min-height: 36px; max-height: 140px; resize: vertical; font-size: 12.5px; }
.ai-msg { border-radius: 9px; padding: 7px 9px; font-size: 12.5px; line-height: 1.5; word-wrap: break-word; overflow-wrap: anywhere; }
.ai-msg.user { background: var(--accent-d); border: 1px solid var(--accent-b, rgba(59,130,246,.3)); align-self: flex-end; max-width: 92%; white-space: pre-wrap; }
.ai-msg.assistant { background: var(--bg3); border: 1px solid var(--border2); }
.ai-msg.assistant pre, .ai-card pre { background: #0b0d12; border: 1px solid var(--border); border-radius: 6px; padding: 6px 8px; overflow-x: auto; font: 11.5px/1.45 'JetBrains Mono', monospace; white-space: pre; margin: 5px 0; max-height: 320px; }
.ai-msg code { font: 11.5px 'JetBrains Mono', monospace; background: var(--bg4); padding: 0 4px; border-radius: 4px; }
.ai-note { font-size: 11px; color: var(--text2); text-align: center; }
.ai-note.err { color: #fca5a5; }
.ai-card { border: 1px solid var(--border2); border-radius: 9px; background: var(--bg); padding: 7px 8px; font-size: 12px; display: flex; flex-direction: column; gap: 5px; min-width: 0; }
.ai-card.wait { border-color: rgba(245,158,11,.55); box-shadow: 0 0 0 1px rgba(245,158,11,.15); }
.ai-card.deny { border-color: rgba(239,68,68,.4); }
.ai-card .ai-ch { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
.ai-card .ai-tool { font-weight: 700; font-size: 11px; color: var(--text2); text-transform: uppercase; letter-spacing: .04em; }
.ai-card .ai-cmd { font: 12px/1.45 'JetBrains Mono', monospace; color: var(--text); word-break: break-all; white-space: pre-wrap; }
.ai-card .ai-why { font-size: 11px; color: var(--text2); }
.ai-card textarea { font: 12px/1.4 'JetBrains Mono', monospace; min-height: 54px; width: 100%; }
.ai-card .ai-acts { display: flex; gap: 6px; flex-wrap: wrap; align-items: center; }
.ai-card .ai-acts input { flex: 1; min-width: 120px; font-size: 11.5px; }
.ai-diff .a { color: #86efac; } .ai-diff .d { color: #fca5a5; } .ai-diff .h { color: #93c5fd; }
.ai-card details summary { cursor: pointer; font-size: 11px; color: var(--text2); }
.ai-setup { display: flex; flex-direction: column; gap: 9px; font-size: 12px; }
.ai-setup label { font-size: 11px; color: var(--text2); display: block; margin-bottom: 3px; }
.ai-setup select, .ai-setup input[type=text], .ai-setup input[type=number], .ai-setup textarea { width: 100%; }
.ai-setup textarea { min-height: 52px; font: 11.5px 'JetBrains Mono', monospace; }
.ai-row { display: flex; gap: 8px; } .ai-row > div { flex: 1; min-width: 0; }
.ai-check { display: flex; gap: 7px; align-items: flex-start; font-size: 12px; color: var(--text); }
.ai-check input { margin-top: 2px; flex: none; }
.ai-see { font-size: 11.5px; color: var(--text2); line-height: 1.5; }
.ai-see li { margin-left: 16px; }
.ai-bar { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; font-size: 11px; color: var(--text2); position: sticky; top: -8px; z-index: 2; background: var(--bg2); padding: 6px 0 4px; margin-top: -6px; }
.ai-bar select { font-size: 11.5px; padding: 2px 6px; max-width: 100%; }
#admin-body button.btn-sm.danger, .ai-modal button.btn-sm.danger, .ai-panel button.danger { background: #dc2626; border-color: #dc2626; color: #fff; }
#admin-body button.btn-sm.danger:hover, .ai-modal button.btn-sm.danger:hover { background: #b91c1c; }
.ai-modal .modal { width: 760px; max-width: calc(100vw - 24px); max-height: calc(100vh - 40px); overflow-y: auto; }
.ai-tr { display: flex; flex-direction: column; gap: 6px; }
.ai-tr .k { font-size: 10px; font-weight: 700; color: var(--text3); text-transform: uppercase; }
.ai-form .field { margin-bottom: 9px; }
.ai-form .field label { font-size: 11px; color: var(--text2); display: block; margin-bottom: 3px; }
.ai-form input[type=text], .ai-form input[type=password], .ai-form input[type=number], .ai-form select, .ai-form textarea { width: 100%; }
.ai-form textarea { min-height: 90px; font: 11px 'JetBrains Mono', monospace; }
@media (max-width: 480px) { .ai-row { flex-direction: column; gap: 6px; } .ai-modal .modal { padding: 14px; } }
`;
  if (!$('ai-style')) { const s = document.createElement('style'); s.id = 'ai-style'; s.textContent = css; document.head.appendChild(s); }

  const tf = (k, vars) => { let s = t(k); for (const [a, b] of Object.entries(vars || {})) s = s.split('{' + a + '}').join(String(b)); return s; };
  const winOf = id => winObjects.find(w => w.id === id);
  const shareQ = wo => (wo && wo.shareToken ? '?share_token=' + encodeURIComponent(wo.shareToken) : '');
  async function api(method, url, body) {
    const res = await fetch(url, {method, headers: body ? {'Content-Type': 'application/json'} : {}, body: body ? JSON.stringify(body) : undefined});
    if (!res.ok) throw new Error(await apiError(res));
    return res.json().catch(() => ({}));
  }
  function md(s) {
    // escape first, then code blocks, inline code, bold and line breaks
    const parts = String(s || '').split('```');
    return parts.map((p, i) => {
      if (i % 2) return '<pre>' + esc(p.replace(/^[a-z0-9_+-]*\n/i, '')) + '</pre>';
      return esc(p).replace(/`([^`\n]+)`/g, '<code>$1</code>').replace(/\*\*([^*\n]+)\*\*/g, '<b>$1</b>').replace(/\n/g, '<br>');
    }).join('');
  }
  function diffHTML(d) {
    return '<pre class="ai-diff">' + String(d || '').split('\n').map(l => {
      const c = l.startsWith('@@') ? 'h' : l.startsWith('+') && !l.startsWith('+++') ? 'a' : l.startsWith('-') && !l.startsWith('---') ? 'd' : '';
      return c ? `<span class="${c}">${esc(l)}</span>` : esc(l);
    }).join('\n') + '</pre>';
  }

  async function loadCfg() { cfg = await api('GET', '/api/ai/config'); return cfg; }

  // ── panel ────────────────────────────────────────────
  function panelEl(id) { return $('aip-' + id); }

  function ensurePanel(wo) {
    let p = panelEl(wo.id);
    if (p) return p;
    p = document.createElement('div');
    p.className = 'ai-panel';
    p.id = 'aip-' + wo.id;
    p.innerHTML = `<div class="ai-head"><span class="ai-ttl">✦ ${esc(t('ai_title'))}</span><span class="pill info" id="aipm-${wo.id}" style="display:none;"></span><span class="ai-sp"></span>
      <button class="ai-ib" title="${esc(t('ai_see'))}" onclick="wrmAI.see(${wo.id})">👁</button>
      <button class="ai-ib" title="${esc(t('ai_history'))}" onclick="wrmAI.history(${wo.id})">🕘</button>
      <button class="ai-ib danger" id="aipk-${wo.id}" title="${esc(t('ai_end'))}" style="display:none;" onclick="wrmAI.end(${wo.id})">⏻</button>
      <button class="ai-ib" title="${esc(t('ai_close_panel'))}" onclick="wrmAI.toggle(${wo.id})">✕</button></div>
      <div class="ai-body" id="aib-${wo.id}"></div><div class="ai-foot" id="aif-${wo.id}" style="display:none;"></div>`;
    wo.el.appendChild(p);
    return p;
  }

  function layout(wo) {
    const p = panelEl(wo.id);
    if (p) {
      const hdr = $('wh-' + wo.id);
      p.style.top = (hdr ? hdr.offsetHeight : 34) + 'px';
    }
    setTimeout(() => { try { fitTerm(wo); } catch (_) {} }, 60);
  }

  async function toggle(id) {
    const wo = winOf(id);
    if (!wo) return;
    if (wo.el.classList.contains('ai-open')) {
      wo.el.classList.remove('ai-open');
      const p = panelEl(id); if (p) p.style.display = 'none';
      layout(wo); indicators();
      return;
    }
    const p = ensurePanel(wo);
    p.style.display = '';
    wo.el.classList.add('ai-open');
    layout(wo);
    const st = S[id];
    if (st && st.sid && !st.ended) { indicators(); return; }
    await renderSetup(wo);
  }

  // setup: provider, model, mode, notes, automatic limits, what the AI sees
  async function renderSetup(wo) {
    const id = wo.id, body = $('aib-' + id);
    $('aif-' + id).style.display = 'none';
    $('aipk-' + id).style.display = 'none';
    $('aipm-' + id).style.display = 'none';
    body.innerHTML = `<div class="fm-loading">${esc(t('loading'))}</div>`;
    let pol;
    try { await loadCfg(); pol = await api('GET', `/api/ai/connections/${wo.connId}${shareQ(wo)}`); }
    catch (e) { body.innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    S[id] = Object.assign(S[id] || {}, {pol});
    if (!cfg.enabled) { body.innerHTML = `<div class="hint-box warn">${esc(t('ai_disabled'))}: ${esc(cfg.reason || '')}</div>`; return; }
    if (!pol.allowed_modes.length) { body.innerHTML = `<div class="hint-box warn">${esc(t('ai_disabled'))}</div>`; return; }
    const resume = (cfg.active || []).filter(s => s.conn_id === wo.connId && !Object.values(S).some(x => x.sid === s.id));
    const provs = cfg.providers || [];
    const own = cfg.personal_keys ? `<button class="btn-sec btn-sm" onclick="wrmAI.myKeys(${id})">🔑 ${esc(t('ai_my_keys'))}</button>` : '';
    if (!provs.length) {
      body.innerHTML = `<div class="hint-box">${esc(t('ai_no_providers'))}</div>${cfg.personal_keys ? `<button class="btn-sm" onclick="wrmAI.myKeys(${id})">🔑 ${esc(t('ai_add_own_key'))}</button>` : ''}
        <div class="hint-box">${esc(t('ai_subscriptions_note'))}</div>`;
      return;
    }
    const shared = !!wo.shareToken;
    body.innerHTML = `<div class="ai-setup">
      ${resume.map(s => `<button class="btn-sec btn-sm" onclick="wrmAI.attach(${id}, ${s.id})">↪ ${esc(t('ai_continue'))} #${s.id} · ${esc(t('ai_mode_' + s.mode))} · ${esc(s.model)}</button>`).join('')}
      <div class="ai-row"><div><label>${esc(t('ai_provider'))}</label><select id="ai-prov-${id}" onchange="wrmAI.models(${id})">${provs.map(p => `<option value="${p.id}">${esc(p.name)}${p.scope === 'user' ? ' (🔑)' : ''}</option>`).join('')}</select></div>
        <div><label>${esc(t('ai_model'))}</label><select id="ai-model-${id}"></select></div></div>
      <div><label>${esc(t('ai_mode'))}</label><select id="ai-mode-${id}" onchange="wrmAI.modeHelp(${id})">${pol.allowed_modes.map(m => `<option value="${m}" ${m === pol.default_mode ? 'selected' : ''}>${esc(t('ai_mode_' + m))}</option>`).join('')}</select>
        <div class="ai-see" id="ai-mode-d-${id}" style="margin-top:4px;"></div></div>
      <div id="ai-auto-${id}" style="display:none;">${autoForm('ai-a-' + id)}</div>
      ${pol.has_notes && pol.share_notes_allowed && !shared ? `<label class="ai-check"><input type="checkbox" id="ai-notes-${id}"> ${esc(t('ai_share_notes'))}</label>` : ''}
      <details><summary class="ai-see">👁 ${esc(t('ai_see'))}</summary><div id="ai-see-${id}">${seeHTML(pol.sees, provs[0] && provs[0].name, false)}</div></details>
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;"><button onclick="wrmAI.start(${id})">✦ ${esc(t('ai_start'))}</button>${own}</div>
      <div class="ai-see">${esc(t('ai_subscriptions_note'))}</div></div>`;
    models(id); modeHelp(id);
  }

  function autoForm(p) {
    const mx = cfg ? cfg.auto_max_minutes : 60, ma = cfg ? cfg.auto_max_actions : 50;
    return `<div class="hint-box"><b>${esc(t('ai_mode_auto_setup'))}</b>
      <div style="margin-top:6px;"><label>${esc(t('ai_auto_allow'))}</label><textarea id="${p}-allow" placeholder="systemctl restart nginx&#10;apt-get install -y *"></textarea></div>
      <div><label>${esc(t('ai_auto_deny'))}</label><textarea id="${p}-deny" placeholder="* /etc/ssh/*"></textarea></div>
      <div><label>${esc(t('ai_auto_path_allow'))}</label><textarea id="${p}-pallow" placeholder="/etc/nginx/**"></textarea></div>
      <div><label>${esc(t('ai_auto_path_deny'))}</label><textarea id="${p}-pdeny"></textarea></div>
      <div class="ai-row"><div><label>${esc(t('ai_auto_minutes'))}</label><input type="number" id="${p}-min" min="1" max="${mx}" value="${Math.min(30, mx)}"></div>
        <div><label>${esc(t('ai_auto_actions'))}</label><input type="number" id="${p}-act" min="1" max="${ma}" value="${Math.min(20, ma)}"></div></div>
      <label class="ai-check" style="margin-top:6px;"><input type="checkbox" id="${p}-ok"> ${esc(t('ai_auto_confirm'))}</label></div>`;
  }
  function readAuto(p) {
    const lines = x => ($(p + '-' + x).value || '').split('\n').map(s => s.trim()).filter(Boolean);
    return {auto: {allow: lines('allow'), deny: lines('deny'), path_allow: lines('pallow'), path_deny: lines('pdeny'),
      minutes: Number($(p + '-min').value) || 0, max_actions: Number($(p + '-act').value) || 0}, ok: $(p + '-ok').checked};
  }

  function seeHTML(sees, provider, notes) {
    if (!sees) return '';
    return `<ul class="ai-see">
      <li>${esc(t('ai_see_tools'))}: ${sees.read_tools.map(esc).join(', ')}</li>
      <li>${esc(t('ai_see_cmd'))}</li>
      <li>${esc(t('ai_see_write'))}: ${sees.write_tools.map(esc).join(', ')}</li>
      <li>${esc(sees.redacted ? t('ai_see_redacted') : t('ai_see_not_redacted'))}</li>
      <li>${esc(tf('ai_see_limit', {kb: sees.read_max_kb}))}</li>
      <li>${esc(tf('ai_see_denied', {n: sees.denied_paths}))}</li>
      <li>${esc(t('ai_see_never'))}</li>
      ${provider ? `<li>${esc(tf('ai_see_provider', {provider}))}</li>` : ''}
      ${notes ? `<li>${esc(t('ai_see_notes_on'))}</li>` : ''}</ul>`;
  }

  function models(id) {
    const sel = $('ai-prov-' + id), m = $('ai-model-' + id);
    const p = (cfg.providers || []).find(x => String(x.id) === sel.value);
    if (!p) return;
    m.innerHTML = p.models.map(x => `<option ${x === p.default_model ? 'selected' : ''}>${esc(x)}</option>`).join('');
    const see = $('ai-see-' + id);
    if (see && S[id] && S[id].pol) see.innerHTML = seeHTML(S[id].pol.sees, p.name, false);
  }
  function modeHelp(id) {
    const m = $('ai-mode-' + id).value;
    $('ai-mode-d-' + id).textContent = t('ai_mode_' + m + '_d');
    $('ai-auto-' + id).style.display = m === 'auto' ? '' : 'none';
  }

  async function start(id) {
    const wo = winOf(id);
    const mode = $('ai-mode-' + id).value;
    const body = {conn_id: wo.connId, provider_id: Number($('ai-prov-' + id).value), model: $('ai-model-' + id).value, mode,
      share_notes: !!($('ai-notes-' + id) && $('ai-notes-' + id).checked)};
    if (mode === 'auto') {
      const a = readAuto('ai-a-' + id);
      if (!a.ok) { showToast(t('ai_auto_need_confirm'), 'warning', 4000); return; }
      body.auto = a.auto; body.confirm_auto = true;
    }
    try {
      const v = await api('POST', '/api/ai/sessions' + shareQ(wo), body);
      attach(id, v.id);
    } catch (e) { showToast(e.message, 'error', 6000); }
  }

  // chat view bound to a session; the event stream replays everything from the start
  function attach(id, sid) {
    const wo = winOf(id);
    if (!wo) return;
    const old = S[id] || {};
    if (old.es) old.es.close();
    S[id] = {pol: old.pol, sid, seq: 0, ended: false, busy: false, cards: {}, cur: null};
    const body = $('aib-' + id);
    body.innerHTML = `<div class="ai-bar" id="ai-bar-${id}"></div>`;
    $('aipk-' + id).style.display = '';
    const foot = $('aif-' + id);
    foot.style.display = '';
    foot.innerHTML = `<textarea id="ai-in-${id}" rows="2" placeholder="${esc(t('ai_placeholder'))}" onkeydown="if(event.key==='Enter'&&!event.shiftKey){event.preventDefault();wrmAI.send(${id});}"></textarea>
      <button id="ai-send-${id}" onclick="wrmAI.send(${id})">${esc(t('ai_send'))}</button>
      <button id="ai-stop-${id}" class="btn-sec" style="display:none;" onclick="wrmAI.stop(${id})">■ ${esc(t('ai_stop'))}</button>`;
    const es = new EventSource(`/api/ai/sessions/${sid}/events?since=0`);
    S[id].es = es;
    ['session', 'user', 'assistant_start', 'text', 'assistant_done', 'tool_call', 'approval_required', 'approval_resolved', 'tool_result', 'mode', 'usage', 'busy', 'error', 'notice', 'ended']
      .forEach(ty => es.addEventListener(ty, ev => { try { onEvent(id, JSON.parse(ev.data)); } catch (e) { console.warn(e); } }));
    es.onerror = () => { if (S[id] && S[id].ended) es.close(); };
    indicators();
    if (!S[id].timer) S[id].timer = setInterval(() => tick(id), 1000);
  }

  function scrollEnd(id) { const b = $('aib-' + id); if (b && b.scrollHeight - b.scrollTop - b.clientHeight < 160) b.scrollTop = b.scrollHeight; }
  function add(id, html) { const b = $('aib-' + id); if (!b) return null; b.insertAdjacentHTML('beforeend', html); const el = b.lastElementChild; scrollEnd(id); return el; }

  function bar(id) {
    const st = S[id], v = st.view, el = $('ai-bar-' + id);
    if (!el || !v) return;
    const allowed = (st.pol && st.pol.allowed_modes) || [v.mode];
    const auto = v.auto ? ` · ${esc(tf('ai_auto_until', {time: new Date(v.auto.until).toLocaleTimeString(), used: v.auto.used, max: v.auto.max_actions}))}` : '';
    el.innerHTML = `<select onchange="wrmAI.setMode(${id}, this.value)" ${st.ended ? 'disabled' : ''}>${MODES.filter(m => allowed.includes(m) || m === v.mode).map(m => `<option value="${m}" ${m === v.mode ? 'selected' : ''}>${esc(t('ai_mode_' + m))}</option>`).join('')}</select>
      <span>${esc(v.provider)} · ${esc(v.model)}${auto}</span><span style="margin-left:auto;" id="ai-use-${id}"></span>`;
    usage(id);
    const pm = $('aipm-' + id);
    pm.style.display = ''; pm.textContent = t('ai_mode_' + v.mode);
  }
  function usage(id) {
    const v = S[id].view, el = $('ai-use-' + id);
    if (el && v) el.textContent = tf('ai_usage', {in: v.tokens_in || 0, out: v.tokens_out || 0}) + (v.cost_usd ? ' · $' + Number(v.cost_usd).toFixed(4) : '');
  }

  // the server's reasons for ending a session, translated where known
  function reasonText(r) {
    r = String(r || '');
    const map = {'ended by the user': 'ai_r_user', 'all sessions ended by the user': 'ai_r_user', 'killed by an administrator': 'ai_r_admin',
      'killed by an administrator (all sessions of the user)': 'ai_r_admin', 'kill switch: every AI session ended by an administrator': 'ai_r_killall',
      'kill switch turned on by an administrator': 'ai_r_switch', 'the AI assistant was turned off': 'ai_r_off',
      'the AI assistant was turned off for the user (kill switch)': 'ai_r_blocked', 'the policy no longer allows the AI assistant on this connection': 'ai_r_policy'};
    if (map[r]) return t(map[r]);
    if (r.startsWith('idle')) return t('ai_r_idle');
    return r;
  }

  function decisionPill(d) {
    const cls = {allow: 'ok', approve: 'warn', deny: 'bad'}[d.decision] || '';
    const lbl = d.decision === 'allow' ? (d.rule === 'auto' ? t('ai_d_auto') : d.read_only ? t('ai_d_readonly') : t('ai_d_allow')) : t('ai_d_' + d.decision);
    return `<span class="pill ${cls}">${esc(lbl)}</span>`;
  }

  function onEvent(id, ev) {
    const st = S[id];
    if (!st || ev.seq <= st.seq) return;
    st.seq = ev.seq;
    const d = ev.data || {};
    switch (ev.type) {
      case 'session': st.view = d; st.ended = d.ended; bar(id); (d.pending || []).forEach(a => approvalUI(id, a)); break;
      case 'user': add(id, `<div class="ai-msg user">${esc(d.text)}</div>`); break;
      case 'assistant_start': st.cur = null; break;
      case 'text':
        if (!st.cur) { st.cur = add(id, `<div class="ai-msg assistant"></div>`); st.curText = ''; }
        st.curText += d.delta; st.cur.textContent = st.curText; scrollEnd(id); break;
      case 'assistant_done':
        if (d.text) { if (!st.cur) st.cur = add(id, `<div class="ai-msg assistant"></div>`); st.cur.innerHTML = md(d.text); }
        st.cur = null; break;
      case 'tool_call': {
        const el = add(id, `<div class="ai-card ${d.decision === 'deny' ? 'deny' : ''}" id="aic-${id}-${esc(d.call_id)}">
          <div class="ai-ch"><span class="ai-tool">${esc(d.tool)}</span>${decisionPill(d)}</div>
          <div class="ai-cmd">${esc(d.display)}</div>
          ${d.reason ? `<div class="ai-why">${esc(t('ai_reason'))}: ${esc(d.reason)}</div>` : ''}
          ${d.why && d.decision !== 'allow' ? `<div class="ai-why">⚑ ${esc(d.why)}</div>` : ''}
          ${d.diff && d.decision !== 'approve' ? diffHTML(d.diff) : ''}
          <div class="ai-appr"></div><div class="ai-res"></div></div>`);
        st.cards[d.call_id] = el; break;
      }
      case 'approval_required': approvalUI(id, d.approval); break;
      case 'approval_resolved': {
        const card = st.cards[d.call_id];
        if (card) {
          card.classList.remove('wait');
          const a = card.querySelector('.ai-appr');
          const cls = d.outcome === 'approved' ? 'ok' : 'bad';
          a.innerHTML = `<div class="ai-ch"><span class="pill ${cls}">${esc(t('ai_o_' + d.outcome))}${d.edited ? ' · ' + esc(t('ai_edited')) : ''}</span>${d.by ? `<span class="ai-why">${esc(d.by)}</span>` : ''}</div>`;
        }
        indicators(); break;
      }
      case 'tool_result': {
        const card = st.cards[d.call_id];
        const head = d.error ? `<div class="ai-why" style="color:#fca5a5;">${esc(d.error)}</div>`
          : `<div class="ai-ch"><span class="pill ${d.exit_code === 0 ? '' : 'warn'}">${esc(tf('ai_exit', {code: d.exit_code}))}</span>${d.redactions ? `<span class="pill info">🔒 ${esc(tf('ai_redactions', {n: d.redactions}))}</span>` : ''}</div>`;
        const out = d.output != null ? `<details ${String(d.output).length < 600 ? 'open' : ''}><summary>${esc(t('ai_output'))} (${String(d.output).length} B)</summary><pre>${esc(d.output)}</pre></details>` : '';
        if (card) card.querySelector('.ai-res').innerHTML = head + out;
        else add(id, `<div class="ai-card ${d.error ? 'deny' : ''}"><div class="ai-ch"><span class="ai-tool">${esc(d.tool || '')}</span>${d.error ? `<span class="pill bad">${esc(t('ai_d_deny'))}</span>` : ''}</div>${head}${out}</div>`);
        scrollEnd(id); break;
      }
      case 'mode':
        if (st.view) { st.view.mode = d.mode; st.view.auto = d.auto; }
        add(id, `<div class="ai-note">${esc(tf('ai_mode_changed', {mode: t('ai_mode_' + d.mode)}))} · ${esc(tf('ai_by', {who: d.by}))}${d.why ? ' · ' + esc(d.why) : ''}</div>`);
        bar(id); indicators(); break;
      case 'usage': if (st.view) Object.assign(st.view, d); usage(id); break;
      case 'busy': {
        st.busy = d.busy;
        const s1 = $('ai-send-' + id), s2 = $('ai-stop-' + id);
        if (s1) s1.style.display = d.busy ? 'none' : ''; if (s2) s2.style.display = d.busy ? '' : 'none';
        break;
      }
      case 'error': add(id, `<div class="ai-note err">⚠ ${esc(d.error)}</div>`); break;
      case 'notice': add(id, `<div class="ai-note">${esc(t('ai_n_' + d.text))}</div>`); break;
      case 'ended':
        st.ended = true;
        if (st.es) st.es.close();
        add(id, `<div class="ai-note err">⏻ ${esc(t('ai_ended'))}: ${esc(reasonText(d.reason))}</div><button class="btn-sm" style="align-self:center;" onclick="wrmAI.newSession(${id})">✦ ${esc(t('ai_new'))}</button>`);
        const f = $('aif-' + id); if (f) f.style.display = 'none';
        $('aipk-' + id).style.display = 'none';
        if (st.view) st.view.ended = true;
        bar(id); indicators(); break;
    }
  }

  function approvalUI(id, a) {
    const st = S[id], card = st.cards[a.call_id];
    if (!card) return;
    card.classList.add('wait');
    const box = card.querySelector('.ai-appr');
    const isWrite = !!a.path;
    box.innerHTML = `<div class="ai-why"><b>⏳ ${esc(t('ai_wait_approval'))}</b> · <span data-exp="${new Date(a.expires).getTime()}"></span></div>
      ${isWrite ? diffHTML(a.diff) + `<details><summary>✏️ ${esc(t('ai_edit'))}</summary><textarea id="ai-ed-${id}-${esc(a.id)}">${esc(a.content)}</textarea></details>`
        : `<textarea id="ai-ed-${id}-${esc(a.id)}" spellcheck="false">${esc(a.command)}</textarea>`}
      <div class="ai-acts"><input type="text" id="ai-nt-${id}-${esc(a.id)}" placeholder="${esc(t('ai_note_ph'))}">
        <button class="btn-sm" onclick="wrmAI.decide(${id}, '${esc(a.id)}', true, ${isWrite})">✓ ${esc(t('ai_approve'))}</button>
        <button class="btn-sec btn-sm" onclick="wrmAI.decide(${id}, '${esc(a.id)}', false, ${isWrite})">✕ ${esc(t('ai_deny'))}</button></div>`;
    tick(id); scrollEnd(id); indicators();
  }

  function tick(id) {
    const b = $('aib-' + id);
    if (!b) { const st = S[id]; if (st && st.timer) { clearInterval(st.timer); st.timer = null; } return; }
    b.querySelectorAll('[data-exp]').forEach(el => {
      const s = Math.max(0, Math.round((Number(el.dataset.exp) - Date.now()) / 1000));
      el.textContent = tf('ai_expires', {s});
    });
  }

  async function decide(id, aid, ok, isWrite) {
    const wo = winOf(id), st = S[id];
    const ed = $(`ai-ed-${id}-${aid}`), note = $(`ai-nt-${id}-${aid}`);
    const body = {decision: ok ? 'approve' : 'deny', note: note ? note.value : ''};
    if (ok && ed) { if (isWrite) body.content = ed.value; else body.command = ed.value; }
    try { await api('POST', `/api/ai/sessions/${st.sid}/approvals/${encodeURIComponent(aid)}${shareQ(wo)}`, body); }
    catch (e) { showToast(e.message, 'error', 5000); }
  }

  async function send(id) {
    const wo = winOf(id), st = S[id], inp = $('ai-in-' + id);
    let text = (inp.value || '').trim();
    if (!text || !st || !st.sid) return;
    if (wo.mode === 'sftp' && wo.path) text = `[${t('ai_fm_ctx')}: ${wo.path}]\n` + text;
    try { await api('POST', `/api/ai/sessions/${st.sid}/prompt${shareQ(wo)}`, {text}); inp.value = ''; }
    catch (e) { showToast(e.message, 'error', 5000); }
  }
  async function stop(id) { const st = S[id]; if (st && st.sid) api('POST', `/api/ai/sessions/${st.sid}/stop${shareQ(winOf(id))}`).catch(e => showToast(e.message, 'error')); }
  async function end(id) {
    const st = S[id];
    if (!st || !st.sid || !(await uiConfirm(t('ai_end_q'), {okText: t('ai_end'), danger: true}))) return;
    api('POST', `/api/ai/sessions/${st.sid}/kill${shareQ(winOf(id))}`).catch(e => showToast(e.message, 'error'));
  }
  function newSession(id) { const st = S[id]; if (st && st.es) st.es.close(); S[id] = {pol: st && st.pol}; renderSetup(winOf(id)); indicators(); }

  async function setMode(id, mode) {
    const wo = winOf(id), st = S[id];
    const body = {mode};
    if (mode === 'auto') {
      const r = await autoDialog();
      if (!r) { bar(id); return; }
      body.auto = r; body.confirm = true;
    }
    try { st.view = await api('POST', `/api/ai/sessions/${st.sid}/mode${shareQ(wo)}`, body); bar(id); }
    catch (e) { showToast(e.message, 'error', 6000); bar(id); }
  }

  // a small modal for the automatic mode opt-in during a session
  function modal(html) {
    let m = $('ai-modal');
    if (!m) { m = document.createElement('div'); m.id = 'ai-modal'; m.className = 'modal-overlay ai-modal'; m.style.zIndex = 9600; document.body.appendChild(m); }
    m.innerHTML = `<div class="modal">${html}</div>`;
    m.classList.add('open');
    return m;
  }
  function closeAIModal() { const m = $('ai-modal'); if (m) m.classList.remove('open'); }
  function autoDialog() {
    return new Promise(resolve => {
      modal(`<div style="display:flex;align-items:center;gap:10px;"><h3 style="margin:0;flex:1;">✦ ${esc(t('ai_mode_auto'))}</h3><button class="icon-btn" id="ai-ad-x">✕</button></div>
        <div class="ai-see" style="margin:8px 0;">${esc(t('ai_mode_auto_d'))}</div><div class="ai-setup">${autoForm('ai-ad')}</div>
        <div class="modal-footer"><button class="btn-sec" id="ai-ad-c">${esc(t('cancel'))}</button><button id="ai-ad-ok">${esc(t('ai_start'))}</button></div>`);
      const done = v => { closeAIModal(); resolve(v); };
      $('ai-ad-x').onclick = $('ai-ad-c').onclick = () => done(null);
      $('ai-ad-ok').onclick = () => { const a = readAuto('ai-ad'); if (!a.ok) { showToast(t('ai_auto_need_confirm'), 'warning'); return; } done(a.auto); };
    });
  }

  function see(id) {
    const st = S[id] || {};
    const prov = st.view ? st.view.provider : '';
    modal(`<div style="display:flex;align-items:center;gap:10px;"><h3 style="margin:0;flex:1;">👁 ${esc(t('ai_see'))}</h3><button class="icon-btn" onclick="wrmAI.closeModal()">✕</button></div>
      ${st.view ? `<div class="hint-box"><b>${esc(t('ai_mode_' + st.view.mode))}</b>: ${esc(t('ai_mode_' + st.view.mode + '_d'))}</div>` : ''}
      ${seeHTML(st.pol && st.pol.sees, prov, st.view && st.view.share_notes)}<div class="hint-box">${esc(t('ai_subscriptions_note'))}</div>`);
  }

  async function history(id) {
    let list = [];
    try { list = await api('GET', '/api/ai/sessions'); } catch (e) { showToast(e.message, 'error'); return; }
    modal(`<div style="display:flex;align-items:center;gap:10px;"><h3 style="margin:0;flex:1;">🕘 ${esc(t('ai_history'))}</h3><button class="icon-btn" onclick="wrmAI.closeModal()">✕</button></div>
      ${sessionsTable(list, false)}`);
  }

  function sessionsTable(list, admin) {
    if (!list.length) return `<div class="chat-empty">${esc(t('none'))}</div>`;
    return `<div class="tbl-scroll"><table class="admin-table"><thead><tr><th>${esc(t('started'))}</th>${admin ? `<th>${esc(t('user'))}</th>` : ''}<th>${esc(t('server'))}</th><th>${esc(t('ai_model'))}</th><th>${esc(t('status'))}</th><th>${esc(t('ai_tokens'))}</th><th></th></tr></thead><tbody>` +
      list.map(s => `<tr><td style="white-space:nowrap;">${esc(fmtTime(s.started_at))}</td>${admin ? `<td>${esc(s.user)}</td>` : ''}
        <td>${esc(s.connection)}<div class="sub mono">${esc(s.host)}</div></td><td>${esc(s.model)}<div class="sub">${esc(s.provider)} · ${esc(t('ai_mode_' + s.mode))}</div></td>
        <td><span class="pill ${s.live ? 'ok' : s.status === 'killed' ? 'warn' : ''}">${esc(s.live ? t('ai_live') : s.status)}</span>${s.end_reason ? `<div class="sub">${esc(reasonText(s.end_reason))}</div>` : ''}</td>
        <td class="mono">${s.tokens_in}/${s.tokens_out}${s.cost_usd ? `<div class="sub">$${Number(s.cost_usd).toFixed(4)}</div>` : ''}<div class="sub">${s.tool_calls} ${esc(t('ai_tool_calls'))}</div></td>
        <td><div style="display:flex;gap:4px;flex-wrap:wrap;"><button class="btn-sec btn-sm" onclick="wrmAI.transcript(${s.id}, ${admin})">${esc(t('ai_transcript'))}</button>
          ${s.recording_session ? `<button class="btn-sec btn-sm" onclick="wrmAI.closeModal(); openReplay(${s.recording_session})">▶</button>` : ''}
          ${admin && s.live ? `<button class="btn-sm danger" onclick="wrmAI.adminKill(${s.id})">${esc(t('ai_kill'))}</button>` : ''}</div></td></tr>`).join('') + '</tbody></table></div>';
  }

  async function transcript(sid, admin) {
    let v;
    try { v = await api('GET', (admin ? '/api/admin/ai/sessions/' : '/api/ai/sessions/') + sid); } catch (e) { showToast(e.message, 'error'); return; }
    const s = v.session || {};
    const icon = {user: '🧑', assistant: '✦', tool_call: '🛠', tool_result: '📄', approval: '⏳', decision: '⚖', mode: '⚙', error: '⚠', system: 'ℹ', end: '⏻'};
    modal(`<div style="display:flex;align-items:center;gap:10px;"><h3 style="margin:0;flex:1;">${esc(t('ai_transcript'))} #${sid}</h3><button class="icon-btn" onclick="wrmAI.closeModal()">✕</button></div>
      <div class="ai-see" style="margin:6px 0 10px;">${esc(s.user || '')} · ${esc(s.connection || '')} (${esc(s.host || '')}) · ${esc(s.provider || '')} / ${esc(s.model || '')} · ${esc(fmtTime(s.started_at))} · ${s.tokens_in || 0}/${s.tokens_out || 0} · $${Number(s.cost_usd || 0).toFixed(4)}</div>
      <div class="ai-tr">${(v.messages || []).map(m => `<div class="ai-card"><div class="ai-ch"><span>${icon[m.kind] || '·'}</span><span class="k">${esc(m.kind)}</span><span class="ai-why">${esc(m.actor)} · ${esc(fmtTime(m.ts))}</span></div>
        ${m.kind === 'assistant' ? `<div class="ai-msg assistant">${md(m.content)}</div>` : `<pre>${esc(m.content)}</pre>`}</div>`).join('')}</div>`);
  }

  // ── indicators on the windows ─────────────────────────
  function indicators() {
    winObjects.forEach(wo => {
      const st = S[wo.id], b = $('wai-' + wo.id), btn = $('wb-ai-' + wo.id);
      const on = !!(st && st.sid && !st.ended);
      if (b) {
        b.style.display = on ? '' : 'none';
        if (on && st.view) b.title = tf('ai_active_on', {mode: t('ai_mode_' + st.view.mode)});
      }
      if (btn) btn.classList.toggle('on', on || wo.el.classList.contains('ai-open'));
    });
  }

  function onWinClose(wo) {
    const st = S[wo.id];
    if (st && st.sid && !st.ended) api('POST', `/api/ai/sessions/${st.sid}/kill${shareQ(wo)}`).catch(() => {});
    if (st && st.es) st.es.close();
    if (st && st.timer) clearInterval(st.timer);
    delete S[wo.id];
  }

  async function refresh() {
    // sessions ended elsewhere (admin kill, policy): the event stream reports it; nothing else to poll
    indicators();
  }

  // ── providers (personal and organisation) ─────────────
  const KINDS = ['anthropic', 'openai', 'azure', 'bedrock', 'vertex', 'openai_compatible'];
  const MODEL_HINT = {anthropic: 'claude-opus-5-5, claude-sonnet-5-5, claude-haiku-5-5', bedrock: 'anthropic.claude-opus-5-5', vertex: 'claude-opus-5-5'};

  function provForm(p, scope) {
    p = p || {kind: 'anthropic', enabled: true, refusal_fallback: true};
    const kinds = scope === 'user' && cfg ? KINDS.filter(k => (cfg.provider_kinds || KINDS).includes(k)) : KINDS;
    const f = (k, lbl, val, type, ph) => `<div class="field" data-k="${k}"><label>${esc(lbl)}</label><input type="${type || 'text'}" id="aipf-${k}" value="${esc(val == null ? '' : val)}" placeholder="${esc(ph || '')}" autocomplete="off"></div>`;
    const sec = (k, lbl, has) => f(k, lbl, '', 'password', has ? t('ai_prov_key_keep') : '');
    return `<div class="ai-form">
      <div class="ai-row"><div class="field"><label>${esc(t('ai_prov_name'))}</label><input type="text" id="aipf-name" value="${esc(p.name || '')}"></div>
        <div class="field"><label>${esc(t('ai_prov_kind'))}</label><select id="aipf-kind" onchange="wrmAI.formKind()" ${p.id ? 'disabled' : ''}>${kinds.map(k => `<option value="${k}" ${k === p.kind ? 'selected' : ''}>${esc(t('ai_kind_' + k))}</option>`).join('')}</select></div></div>
      ${f('models', t('ai_prov_models'), p.models, 'text', MODEL_HINT[p.kind] || '')}
      ${f('default_model', t('ai_prov_default'), p.default_model)}
      <div class="kind-anthropic kind-openai kind-azure kind-bedrock kind-openai_compatible">${sec('api_key', t('ai_prov_key') + ' / ' + t('ai_prov_bkey'), p.has_api_key)}</div>
      <div class="kind-anthropic kind-openai kind-bedrock kind-vertex">${f('base_url', t('ai_prov_base'), p.base_url)}</div>
      <div class="kind-azure kind-openai_compatible">${f('base_url2', p.kind === 'azure' ? t('ai_prov_endpoint') : t('ai_prov_base_req'), p.base_url)}</div>
      <div class="kind-azure ai-row">${f('deployment', t('ai_prov_deployment'), p.deployment)}${f('api_version', t('ai_prov_apiver'), p.api_version, 'text', '2024-10-21')}</div>
      <div class="kind-bedrock kind-vertex ai-row">${f('region', t('ai_prov_region'), p.region, 'text', 'eu-central-1 / global')}<div class="field"><label>${esc(t('ai_prov_api'))}</label><select id="aipf-api_mode"></select></div></div>
      <div class="kind-vertex">${f('project', t('ai_prov_project'), p.project)}<div class="field"><label>${esc(t('ai_prov_sa'))}</label><textarea id="aipf-service_account" placeholder="${esc(p.has_service_account ? t('ai_prov_key_keep') : '{ "type": "service_account", … }')}"></textarea></div></div>
      <div class="kind-bedrock">${sec('access_key_id', t('ai_prov_ak'), p.has_access_key)}${sec('secret_access_key', t('ai_prov_sk'), p.has_access_key)}${sec('session_token', t('ai_prov_st'), false)}</div>
      <div class="ai-row">${f('price_in', t('ai_prov_price'), p.price_in || '', 'number')}${f('price_out', '​', p.price_out || '', 'number')}</div>
      ${f('max_tokens', t('ai_prov_maxtok'), p.max_tokens || 0, 'number')}
      <label class="ai-check kind-anthropic"><input type="checkbox" id="aipf-fallback" ${p.refusal_fallback !== false ? 'checked' : ''}> ${esc(t('ai_prov_fallback'))}</label>
      <label class="ai-check"><input type="checkbox" id="aipf-enabled" ${p.enabled !== false ? 'checked' : ''}> ${esc(t('ai_prov_enabled'))}</label></div>`;
  }
  function formKind(p) {
    const k = $('aipf-kind').value;
    document.querySelectorAll('.ai-form [class*="kind-"]').forEach(el => { el.style.display = el.classList.contains('kind-' + k) ? '' : 'none'; });
    const am = $('aipf-api_mode'), cur = (p && p.api_mode) || am.value;
    const opts = k === 'bedrock' ? ['converse', 'messages'] : k === 'vertex' ? ['anthropic', 'openai'] : [];
    am.innerHTML = opts.map(o => `<option ${o === cur ? 'selected' : ''}>${o}</option>`).join('');
    const ml = $('aipf-models');
    if (ml && !ml.value) ml.placeholder = MODEL_HINT[k] || '';
  }
  function readProvForm(p) {
    const v = id => { const el = $('aipf-' + id); return el ? el.value.trim() : ''; };
    const k = $('aipf-kind').value;
    const out = {name: v('name'), kind: k, models: v('models'), default_model: v('default_model'), price_in: Number(v('price_in')) || 0, price_out: Number(v('price_out')) || 0,
      max_tokens: Number(v('max_tokens')) || 0, enabled: $('aipf-enabled').checked, refusal_fallback: $('aipf-fallback').checked,
      base_url: (k === 'azure' || k === 'openai_compatible') ? v('base_url2') : v('base_url'), deployment: v('deployment'), api_version: v('api_version'),
      region: v('region'), project: v('project'), api_mode: v('api_mode')};
    for (const s of ['api_key', 'access_key_id', 'secret_access_key', 'session_token', 'service_account']) { const x = v(s); if (x) out[s] = x; }
    if (!out.models && !p) out.models = (MODEL_HINT[k] || '');
    return out;
  }

  async function providersDialog(scope, after) {
    const base = scope === 'org' ? '/api/admin/ai/providers' : '/api/ai/providers';
    let list = [];
    try { list = await api('GET', base); } catch (e) { showToast(e.message, 'error'); return; }
    const html = `<div style="display:flex;align-items:center;gap:10px;"><h3 style="margin:0;flex:1;">🔑 ${esc(t(scope === 'org' ? 'ai_providers_org' : 'ai_my_keys'))}</h3><button class="icon-btn" onclick="wrmAI.closeModal()">✕</button></div>
      <div class="hint-box">${esc(t('ai_subscriptions_note'))}</div>${provList(list, scope)}<div style="margin-top:10px;"><button class="btn-sm" onclick="wrmAI.editProv('${scope}', 0)">＋ ${esc(t('ai_prov_add'))}</button></div>`;
    modal(html);
    provAfter = after;
  }
  let provAfter = null, provCache = [];
  function provList(list, scope) {
    provCache = list;
    if (!list.length) return `<div class="chat-empty">${esc(t('ai_prov_none'))}</div>`;
    return `<div class="tbl-scroll"><table class="admin-table"><thead><tr><th>${esc(t('ai_prov_name'))}</th><th>${esc(t('ai_prov_kind'))}</th><th>${esc(t('ai_models_col'))}</th><th>${esc(t('ai_prov_secrets'))}</th><th></th></tr></thead><tbody>` +
      list.map(p => `<tr><td>${esc(p.name)}${p.enabled ? '' : ` <span class="pill">${esc(t('off'))}</span>`}</td><td>${esc(t('ai_kind_' + p.kind))}${p.region ? `<div class="sub">${esc(p.region)}</div>` : ''}</td>
        <td class="mono" style="font-size:11px;">${esc(p.models)}</td><td>${p.has_api_key || p.has_access_key || p.has_service_account ? '🔒' : '—'}</td>
        <td><div style="display:flex;gap:4px;flex-wrap:wrap;"><button class="btn-sec btn-sm" onclick="wrmAI.testProv('${scope}', ${p.id})">${esc(t('ai_prov_test'))}</button>
          <button class="btn-sec btn-sm" onclick="wrmAI.editProv('${scope}', ${p.id})">✏️</button><button class="btn-sec btn-sm" onclick="wrmAI.delProv('${scope}', ${p.id})">🗑</button></div></td></tr>`).join('') + '</tbody></table></div>';
  }
  function editProv(scope, pid) {
    const p = provCache.find(x => x.id === pid);
    modal(`<div style="display:flex;align-items:center;gap:10px;"><h3 style="margin:0;flex:1;">${esc(t(p ? 'ai_prov_edit' : 'ai_prov_add'))}</h3><button class="icon-btn" onclick="wrmAI.closeModal()">✕</button></div>
      ${provForm(p, scope)}<div class="modal-footer"><button class="btn-sec" onclick="wrmAI.providers('${scope}')">${esc(t('cancel'))}</button><button onclick="wrmAI.saveProv('${scope}', ${pid})">${esc(t('save'))}</button></div>`);
    formKind(p);
  }
  async function saveProv(scope, pid) {
    const base = scope === 'org' ? '/api/admin/ai/providers' : '/api/ai/providers';
    const p = provCache.find(x => x.id === pid);
    const body = readProvForm(p);
    try { await api(pid ? 'PUT' : 'POST', pid ? `${base}/${pid}` : base, body); showToast(t('saved'), 'success', 1200); }
    catch (e) { showToast(e.message, 'error', 6000); return; }
    if (scope === 'org' && adminTab === 'ai') { closeAIModal(); adminAI(); } else providersDialog(scope, provAfter);
    if (provAfter) provAfter();
  }
  async function testProv(scope, pid) {
    const base = scope === 'org' ? '/api/admin/ai/providers' : '/api/ai/providers';
    showToast(t('loading'), 'info', 1500);
    try { const r = await api('POST', `${base}/${pid}/test`, {}); r.ok ? showToast(t('ai_prov_test_ok') + (r.reply ? ': ' + r.reply : ''), 'success', 5000) : showToast(r.error, 'error', 8000); }
    catch (e) { showToast(e.message, 'error', 8000); }
  }
  async function delProv(scope, pid) {
    if (!(await uiConfirm(t('ai_prov_del_q'), {okText: t('delete'), danger: true}))) return;
    const base = scope === 'org' ? '/api/admin/ai/providers' : '/api/ai/providers';
    try { await api('DELETE', `${base}/${pid}`); } catch (e) { showToast(e.message, 'error'); return; }
    if (scope === 'org' && adminTab === 'ai') adminAI(); else providersDialog(scope, provAfter);
  }

  // ── admin tab ─────────────────────────────────────────
  const AI_POLICY_KEYS = ['ai_assistant', 'ai_kill_switch', 'ai_modes', 'ai_default_mode', 'ai_mode_rules', 'ai_personal_keys', 'ai_provider_kinds', 'ai_models',
    'ai_redact_output', 'ai_share_notes', 'ai_allow_power', 'ai_blocked_commands', 'ai_read_deny_paths', 'ai_approval_timeout_seconds', 'ai_auto_max_minutes',
    'ai_auto_max_actions', 'ai_command_timeout_seconds', 'ai_output_max_kb', 'ai_max_steps', 'ai_idle_minutes', 'ai_transcript_retention_days'];

  async function adminAI() {
    const body = $('admin-body');
    let sum, provs, sessions, settings;
    try {
      [sum, provs, sessions, settings] = await Promise.all([api('GET', '/api/admin/ai/summary'), api('GET', '/api/admin/ai/providers'),
        api('GET', '/api/admin/ai/sessions?limit=100'), api('GET', '/api/admin/settings')]);
    } catch (e) { body.innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    adminSettingsCache = settings;
    const ks = sum.kill_switch;
    body.innerHTML = `
      <div class="adm-h">⏻ ${esc(t('adm_ai_kill'))}</div>
      <div class="hint-box ${ks ? 'warn' : ''}" style="display:flex;gap:10px;align-items:center;flex-wrap:wrap;"><span style="flex:1;min-width:200px;">${esc(ks ? t('ai_kill_switch_on') : t('ai_kill_switch_off'))} · ${sum.active.length} ${esc(t('ai_live'))}</span>
        <button class="btn-sm danger" onclick="wrmAI.killAll()">⏻ ${esc(t('ai_kill_all'))}</button></div>
      <div class="adm-h">🔑 ${esc(t('ai_providers_org'))}</div>
      <div class="hint-box">${esc(t('ai_subscriptions_note'))}</div>
      ${provList(provs, 'org')}<div style="margin:8px 0 4px;"><button class="btn-sm" onclick="wrmAI.editProv('org', 0)">＋ ${esc(t('ai_prov_add'))}</button></div>
      <div class="adm-h">✦ ${esc(t('ai_sessions'))}</div>${sessionsTable(sessions, true)}
      <div class="adm-h">📊 ${esc(t('ai_usage_users'))}</div>
      <div class="tbl-scroll"><table class="admin-table"><thead><tr><th>${esc(t('user'))}</th><th>${esc(t('ai_sessions'))}</th><th>${esc(t('ai_tokens'))}</th><th>${esc(t('ai_cost'))}</th><th></th></tr></thead><tbody>` +
      sum.usage.map(u => `<tr><td>${esc(u.user)}${u.blocked ? ` <span class="pill bad">${esc(t('ai_blocked'))}</span>` : ''}</td><td>${u.sessions}${u.active ? ` · <span class="pill ok">${u.active} ${esc(t('ai_live'))}</span>` : ''}</td>
        <td class="mono">${u.tokens_in}/${u.tokens_out}</td><td class="mono">$${Number(u.cost_usd).toFixed(4)}</td>
        <td><div style="display:flex;gap:4px;flex-wrap:wrap;">${u.active ? `<button class="btn-sm danger" onclick="wrmAI.killUser(${u.user_id})">${esc(t('ai_kill'))}</button>` : ''}
          <button class="btn-sec btn-sm" onclick="wrmAI.blockUser(${u.user_id}, ${!u.blocked})">${esc(t(u.blocked ? 'ai_unblock' : 'ai_block'))}</button></div></td></tr>`).join('') + `</tbody></table></div>
      <div class="adm-h">⚙ ${esc(t('ai_policies'))}</div><div class="hint-box">${esc(t('ai_mode_rules_help'))}</div>
      <div class="policy-group">${AI_POLICY_KEYS.map(k => settingRowHTML(adminSettingsCache.find(x => x.key === k))).join('')}</div>
      <div class="adm-h">⛔ ${esc(t('ai_destructive'))}</div><ul class="ai-see">${sum.destructive_rules.map(r => `<li><span class="mono">${esc(r.id)}</span> — ${esc(r.what)}</li>`).join('')}</ul>`;
  }
  async function killAll() {
    if (!(await uiConfirm(t('ai_kill_all_q'), {okText: t('ai_kill_all'), danger: true}))) return;
    try { const r = await api('POST', '/api/admin/ai/kill-all'); showToast(tf('ai_killed_n', {n: r.ended}), 'success'); } catch (e) { showToast(e.message, 'error'); }
    adminAI();
  }
  async function killUser(uid) {
    try { const r = await api('POST', `/api/admin/ai/users/${uid}/kill`); showToast(tf('ai_killed_n', {n: r.ended}), 'success'); } catch (e) { showToast(e.message, 'error'); }
    adminAI();
  }
  async function blockUser(uid, blocked) {
    try { await api('PUT', `/api/admin/ai/users/${uid}`, {blocked}); } catch (e) { showToast(e.message, 'error'); }
    adminAI();
  }
  async function adminKill(sid) {
    try { await api('POST', `/api/admin/ai/sessions/${sid}/kill`); } catch (e) { showToast(e.message, 'error'); }
    closeAIModal(); if (adminTab === 'ai') adminAI();
  }

  window.wrmAI = {
    toggle, start, attach, models, modeHelp, send, stop, end, decide, setMode, see, history, transcript, newSession,
    myKeys: id => providersDialog('user', () => { const wo = winOf(id); if (wo && !(S[id] && S[id].sid)) renderSetup(wo); }),
    providers: scope => providersDialog(scope, provAfter), editProv, saveProv, testProv, delProv, formKind: () => formKind(),
    closeModal: closeAIModal, adminTab: adminAI, killAll, killUser, blockUser, adminKill, onWinClose, refresh, indicators,
  };
})();
