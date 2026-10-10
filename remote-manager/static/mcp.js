// AI desktop apps (MCP) of WRM PRO: Settings → AI connections (tokens and setup snippets),
// the "AI connected" indicator with one-click revoke, approval cards for MCP sessions, the
// OAuth consent dialog and the administrators' AI connections tab. Loaded on demand by
// index.html (loadMCP); uses its helpers (t, esc, showToast, apiError, uiConfirm, fmtTime,
// copyText, currentUser, settingRowHTML, adminSettingsCache, adminTab).
// The server decides everything (scopes, modes, approvals); this file only shows it.
'use strict';

Object.assign(LANGS.en, {
  mcp_title: 'AI connections', mcp_intro: 'Let an AI app on your computer or in the browser (Claude, ChatGPT and other MCP clients) work on servers you choose. WRM checks every action against the mode and the scopes below, asks you before changes, and records everything.',
  mcp_off: 'AI connections are not available', mcp_new: 'New AI connection', mcp_name: 'Name', mcp_name_ph: 'e.g. Claude Desktop on my laptop',
  mcp_conns: 'Servers', mcp_conns_filter: 'Filter servers…', mcp_mode: 'Mode', mcp_scopes: 'Scopes', mcp_hours: 'Valid for (hours)', mcp_hours_max: 'at most {n}',
  mcp_create: 'Create token', mcp_none: 'No AI connections yet.', mcp_revoke: 'Revoke', mcp_revoke_all: 'Revoke all', mcp_revoke_q: 'Revoke this AI connection? Its sessions end at once.',
  mcp_revoke_all_q: 'Revoke every AI connection? All AI apps lose access at once.', mcp_revoked_n: '{n} revoked',
  mcp_mode_read_only: 'Read-only', mcp_mode_ask: 'Ask before every change', mcp_mode_auto: 'Automatic within limits',
  mcp_mode_read_only_d: 'The app can only read: logs, files, services and read-only commands. Every change is refused.',
  mcp_mode_ask_d: 'Reads run directly. Every change waits for your approval, in WRM or in the app.',
  mcp_mode_auto_d: 'Changes on your allow list run without asking, within a time and action limit per session. Everything else still asks.',
  mcp_sc_read_logs: 'Read logs, files and system information', mcp_sc_run_readonly: 'Run read-only commands', mcp_sc_run_with_approval: 'Run commands that change something (with approval)',
  mcp_sc_edit_file_with_approval: 'Edit files (with approval, shown as a diff)', mcp_sc_transfer: 'Copy files between the chosen servers (with approval)',
  mcp_auto_allow: 'Commands that may run without asking (one pattern per line)', mcp_auto_path_allow: 'Files that may be changed without asking (/etc/nginx/**)',
  mcp_auto_minutes: 'Time limit per session (minutes)', mcp_auto_actions: 'Maximum actions per session', mcp_auto_confirm: 'I allow these changes without asking, within these limits.',
  mcp_token_once: 'Copy the token now. WRM stores only a hash and cannot show it again.', mcp_token: 'Token', mcp_copy: 'Copy', mcp_url: 'MCP server address',
  mcp_setup: 'Set up your app', mcp_setup_claude_desktop: 'Claude Desktop (local bridge)', mcp_setup_remote: 'Claude (custom connector)', mcp_setup_chatgpt: 'ChatGPT (connector)',
  mcp_setup_generic: 'Other MCP clients', mcp_setup_cli: 'Test from a shell',
  mcp_h_claude_desktop: 'Copy the WRM program (the same file as the server, for your computer\'s system) to your computer, then add this to claude_desktop_config.json (Claude → Settings → Developer → Edit config) and restart Claude. Replace the path with the program\'s full path.',
  mcp_h_remote: 'In Claude (claude.ai or the desktop app): Settings → Connectors → Add custom connector, and enter the address below. Claude then opens WRM, where you sign in and choose the servers, the mode and the scopes. WRM must be reachable from the internet over HTTPS for claude.ai.',
  mcp_h_chatgpt: 'In ChatGPT: Settings → Apps & Connectors → Advanced settings → Developer mode, then Create (or Add custom connector). Enter the address below and choose OAuth. ChatGPT then opens WRM, where you approve the connection. WRM must be reachable from the internet over HTTPS.',
  mcp_h_generic: 'For clients that support remote MCP servers over HTTP with a header (for example Claude Code: claude mcp add --transport http …):',
  mcp_h_cli: 'Check that the token works:', mcp_oauth_box: 'Apps with OAuth (Claude, ChatGPT) need no token: add a custom connector with this address and approve it in WRM.',
  mcp_oauth_off: 'OAuth is turned off by the administrator: use a token.', mcp_kind_pat: 'token', mcp_kind_oauth: 'OAuth',
  mcp_expires: 'expires {time}', mcp_expired: 'expired', mcp_revoked: 'revoked', mcp_last_used: 'last used {time}', mcp_never_used: 'not used yet',
  mcp_live: 'connected', mcp_sessions_n: '{n} open sessions',
  mcp_ind_title: 'Connected AI apps', mcp_on: 'on {servers}',
  // approvals
  mcp_ap_title: '{client} asks to act on {server}', mcp_ap_cmd: 'Command', mcp_ap_edit_file: 'Edit {path}', mcp_ap_why: 'Why WRM asks', mcp_ap_reason: 'Reason given by the AI',
  mcp_ap_approve: 'Approve', mcp_ap_deny: 'Deny', mcp_ap_edit: 'Edit', mcp_ap_approve_edited: 'Approve edited', mcp_ap_expires: 'expires in {s} s', mcp_ap_note: 'Note for the AI (optional)',
  mcp_ap_transfer: 'File transfer',
  // consent
  mcp_cs_title: 'Connect an AI app to WRM', mcp_cs_asks: '{client} wants to work on your servers through WRM.',
  mcp_cs_redirect: 'After your answer you go back to {host}. Approve only if you started this connection in that app just now.',
  mcp_cs_approve: 'Allow', mcp_cs_deny: 'Deny', mcp_cs_expired: 'This request has expired. Start the connection in the app again.', mcp_cs_done: 'Returning to the app…',
  // admin
  mcp_adm_active: 'Active AI connections', mcp_adm_none: 'No active AI connections.', mcp_adm_clients: 'Registered OAuth apps', mcp_adm_remove: 'Remove',
  mcp_adm_remove_q: 'Remove this app? Its AI connections are revoked; it must register again.', mcp_adm_policies: 'Policies', mcp_adm_redirects: 'Redirect addresses',
  mcp_client: 'App', mcp_user: 'User', mcp_state: 'State', mcp_created: 'Created',
  mcp_adm_help: 'MCP clients reach WRM at {url}. The kill switch of the AI assistant also stops every AI connection.',
  p_ai_mcp_enabled: 'AI connections (MCP)', p_ai_mcp_enabled_d: 'Serve the MCP endpoint /mcp and the OAuth endpoints for AI apps.',
  p_ai_mcp: 'Who may use AI connections', p_ai_mcp_d: 'Who may create tokens and approve AI apps.', po_ai_mcp_off: 'Nobody', po_ai_mcp_admins: 'Administrators', po_ai_mcp_all: 'Everybody',
  p_ai_mcp_scopes: 'Allowed scopes', p_ai_mcp_scopes_d: 'read_logs, run_readonly, run_with_approval, edit_file_with_approval, transfer', p_ai_mcp_scopes_ph: 'read_logs,run_readonly',
  p_ai_mcp_modes: 'Allowed modes', p_ai_mcp_modes_d: 'read_only, ask, auto. The AI assistant\'s mode rules narrow them per connection.', p_ai_mcp_modes_ph: 'read_only,ask',
  p_ai_mcp_token_hours: 'Default lifetime (hours)', p_ai_mcp_token_hours_d: 'Lifetime of a new token or OAuth approval when the user does not choose one.',
  p_ai_mcp_max_token_hours: 'Maximum lifetime (hours)', p_ai_mcp_max_token_hours_d: 'The longest lifetime a user may choose.',
  p_ai_mcp_oauth: 'OAuth for AI apps', p_ai_mcp_oauth_d: 'Dynamic client registration and OAuth 2.1 with PKCE (Claude, ChatGPT connectors). Off: tokens only.',
  p_ai_mcp_redirect_uris: 'Allowed redirect addresses', p_ai_mcp_redirect_uris_d: 'OAuth apps may only send users back to these addresses (* = any text; https, or http on localhost).', p_ai_mcp_redirect_uris_ph: 'https://claude.ai/api/mcp/auth_callback',
  p_ai_mcp_allowed_origins: 'Allowed browser origins', p_ai_mcp_allowed_origins_d: 'Web pages (besides WRM) that may call /mcp from a browser (CORS). Empty: none.', p_ai_mcp_allowed_origins_ph: 'https://app.example.com',
  p_ai_mcp_public_url: 'Public address', p_ai_mcp_public_url_d: 'The address AI apps use for WRM, when it differs from what WRM sees (reverse proxy).', p_ai_mcp_public_url_ph: 'https://wrm.example.com',
  p_ai_mcp_elicitation: 'Approvals in the AI app', p_ai_mcp_elicitation_d: 'Ask for approvals also in the AI app (MCP elicitation) when it supports it. Approvals in WRM always work.',
  p_ai_mcp_transfer_max_mb: 'Transfer limit (MB)', p_ai_mcp_transfer_max_mb_d: 'The largest file the transfer scope copies.',
});
Object.assign(LANGS.hr, {
  mcp_title: 'AI veze', mcp_intro: 'Dopustite AI aplikaciji na računalu ili u pregledniku (Claude, ChatGPT i drugi MCP klijenti) da radi na poslužiteljima koje odaberete. WRM provjerava svaku radnju prema načinu rada i opsezima u nastavku, pita vas prije promjena i sve bilježi.',
  mcp_off: 'AI veze nisu dostupne', mcp_new: 'Nova AI veza', mcp_name: 'Naziv', mcp_name_ph: 'npr. Claude Desktop na mom laptopu',
  mcp_conns: 'Poslužitelji', mcp_conns_filter: 'Filtriraj poslužitelje…', mcp_mode: 'Način rada', mcp_scopes: 'Opsezi', mcp_hours: 'Vrijedi (sati)', mcp_hours_max: 'najviše {n}',
  mcp_create: 'Stvori token', mcp_none: 'Još nema AI veza.', mcp_revoke: 'Opozovi', mcp_revoke_all: 'Opozovi sve', mcp_revoke_q: 'Opozvati ovu AI vezu? Njezine sesije odmah završavaju.',
  mcp_revoke_all_q: 'Opozvati sve AI veze? Sve AI aplikacije odmah gube pristup.', mcp_revoked_n: 'Opozvano: {n}',
  mcp_mode_read_only: 'Samo čitanje', mcp_mode_ask: 'Pitaj prije svake promjene', mcp_mode_auto: 'Automatski u granicama',
  mcp_mode_read_only_d: 'Aplikacija može samo čitati: zapise, datoteke, servise i naredbe koje samo čitaju. Svaka promjena se odbija.',
  mcp_mode_ask_d: 'Čitanje se izvršava odmah. Svaka promjena čeka vaše odobrenje, u WRM-u ili u aplikaciji.',
  mcp_mode_auto_d: 'Promjene s vašeg popisa dopuštenih izvršavaju se bez pitanja, unutar vremenskog ograničenja i ograničenja broja radnji po sesiji. Sve ostalo i dalje pita.',
  mcp_sc_read_logs: 'Čitanje zapisa, datoteka i podataka o sustavu', mcp_sc_run_readonly: 'Izvršavanje naredbi koje samo čitaju', mcp_sc_run_with_approval: 'Izvršavanje naredbi koje nešto mijenjaju (uz odobrenje)',
  mcp_sc_edit_file_with_approval: 'Uređivanje datoteka (uz odobrenje, prikazano kao razlike)', mcp_sc_transfer: 'Kopiranje datoteka između odabranih poslužitelja (uz odobrenje)',
  mcp_auto_allow: 'Naredbe koje se smiju izvršiti bez pitanja (jedan uzorak po retku)', mcp_auto_path_allow: 'Datoteke koje se smiju mijenjati bez pitanja (/etc/nginx/**)',
  mcp_auto_minutes: 'Vremensko ograničenje po sesiji (minute)', mcp_auto_actions: 'Najviše radnji po sesiji', mcp_auto_confirm: 'Dopuštam ove promjene bez pitanja, unutar ovih ograničenja.',
  mcp_token_once: 'Kopirajte token sada. WRM sprema samo sažetak i ne može ga ponovno prikazati.', mcp_token: 'Token', mcp_copy: 'Kopiraj', mcp_url: 'Adresa MCP poslužitelja',
  mcp_setup: 'Postavljanje aplikacije', mcp_setup_claude_desktop: 'Claude Desktop (lokalni most)', mcp_setup_remote: 'Claude (prilagođeni konektor)', mcp_setup_chatgpt: 'ChatGPT (konektor)',
  mcp_setup_generic: 'Drugi MCP klijenti', mcp_setup_cli: 'Provjera iz ljuske',
  mcp_h_claude_desktop: 'Kopirajte program WRM (istu datoteku kao poslužitelj, za sustav vašeg računala) na računalo, zatim ovo dodajte u claude_desktop_config.json (Claude → Settings → Developer → Edit config) i ponovno pokrenite Claude. Zamijenite putanju punom putanjom programa.',
  mcp_h_remote: 'U Claudeu (claude.ai ili aplikacija za računalo): Settings → Connectors → Add custom connector, i upišite adresu u nastavku. Claude zatim otvara WRM, gdje se prijavite i odaberete poslužitelje, način rada i opsege. Za claude.ai WRM mora biti dostupan s interneta preko HTTPS-a.',
  mcp_h_chatgpt: 'U ChatGPT-u: Settings → Apps & Connectors → Advanced settings → Developer mode, zatim Create (ili Add custom connector). Upišite adresu u nastavku i odaberite OAuth. ChatGPT zatim otvara WRM, gdje odobravate vezu. WRM mora biti dostupan s interneta preko HTTPS-a.',
  mcp_h_generic: 'Za klijente koji podržavaju udaljene MCP poslužitelje preko HTTP-a sa zaglavljem (na primjer Claude Code: claude mcp add --transport http …):',
  mcp_h_cli: 'Provjerite radi li token:', mcp_oauth_box: 'Aplikacije s OAuth-om (Claude, ChatGPT) ne trebaju token: dodajte prilagođeni konektor s ovom adresom i odobrite ga u WRM-u.',
  mcp_oauth_off: 'Administrator je isključio OAuth: koristite token.', mcp_kind_pat: 'token', mcp_kind_oauth: 'OAuth',
  mcp_expires: 'istječe {time}', mcp_expired: 'isteklo', mcp_revoked: 'opozvano', mcp_last_used: 'zadnja uporaba {time}', mcp_never_used: 'još nije korišteno',
  mcp_live: 'povezano', mcp_sessions_n: 'Otvorenih sesija: {n}',
  mcp_ind_title: 'Povezane AI aplikacije', mcp_on: 'na {servers}',
  mcp_ap_title: '{client} traži radnju na {server}', mcp_ap_cmd: 'Naredba', mcp_ap_edit_file: 'Uređivanje {path}', mcp_ap_why: 'Zašto WRM pita', mcp_ap_reason: 'Razlog koji je naveo AI',
  mcp_ap_approve: 'Odobri', mcp_ap_deny: 'Odbij', mcp_ap_edit: 'Uredi', mcp_ap_approve_edited: 'Odobri uređeno', mcp_ap_expires: 'istječe za {s} s', mcp_ap_note: 'Napomena za AI (neobavezno)',
  mcp_ap_transfer: 'Prijenos datoteke',
  mcp_cs_title: 'Povezivanje AI aplikacije s WRM-om', mcp_cs_asks: '{client} želi raditi na vašim poslužiteljima preko WRM-a.',
  mcp_cs_redirect: 'Nakon odgovora vraćate se na {host}. Odobrite samo ako ste upravo pokrenuli ovu vezu u toj aplikaciji.',
  mcp_cs_approve: 'Dopusti', mcp_cs_deny: 'Odbij', mcp_cs_expired: 'Ovaj zahtjev je istekao. Ponovno pokrenite povezivanje u aplikaciji.', mcp_cs_done: 'Povratak u aplikaciju…',
  mcp_adm_active: 'Aktivne AI veze', mcp_adm_none: 'Nema aktivnih AI veza.', mcp_adm_clients: 'Registrirane OAuth aplikacije', mcp_adm_remove: 'Ukloni',
  mcp_adm_remove_q: 'Ukloniti ovu aplikaciju? Njezine AI veze se opozivaju; mora se ponovno registrirati.', mcp_adm_policies: 'Pravila', mcp_adm_redirects: 'Adrese za povratak',
  mcp_client: 'Aplikacija', mcp_user: 'Korisnik', mcp_state: 'Stanje', mcp_created: 'Stvoreno',
  mcp_adm_help: 'MCP klijenti pristupaju WRM-u na {url}. Prekidač za zaustavljanje AI asistenta zaustavlja i sve AI veze.',
  p_ai_mcp_enabled: 'AI veze (MCP)', p_ai_mcp_enabled_d: 'Posluži MCP krajnju točku /mcp i OAuth krajnje točke za AI aplikacije.',
  p_ai_mcp: 'Tko smije koristiti AI veze', p_ai_mcp_d: 'Tko smije stvarati tokene i odobravati AI aplikacije.', po_ai_mcp_off: 'Nitko', po_ai_mcp_admins: 'Administratori', po_ai_mcp_all: 'Svi',
  p_ai_mcp_scopes: 'Dopušteni opsezi', p_ai_mcp_scopes_d: 'read_logs, run_readonly, run_with_approval, edit_file_with_approval, transfer', p_ai_mcp_scopes_ph: 'read_logs,run_readonly',
  p_ai_mcp_modes: 'Dopušteni načini rada', p_ai_mcp_modes_d: 'read_only, ask, auto. Pravila načina rada AI asistenta ih sužavaju po vezi.', p_ai_mcp_modes_ph: 'read_only,ask',
  p_ai_mcp_token_hours: 'Zadano trajanje (sati)', p_ai_mcp_token_hours_d: 'Trajanje novog tokena ili OAuth odobrenja kad ga korisnik ne odabere.',
  p_ai_mcp_max_token_hours: 'Najdulje trajanje (sati)', p_ai_mcp_max_token_hours_d: 'Najdulje trajanje koje korisnik smije odabrati.',
  p_ai_mcp_oauth: 'OAuth za AI aplikacije', p_ai_mcp_oauth_d: 'Dinamička registracija klijenata i OAuth 2.1 s PKCE-om (konektori za Claude, ChatGPT). Isključeno: samo tokeni.',
  p_ai_mcp_redirect_uris: 'Dopuštene adrese za povratak', p_ai_mcp_redirect_uris_d: 'OAuth aplikacije smiju korisnike vratiti samo na ove adrese (* = bilo koji tekst; https, ili http na localhostu).', p_ai_mcp_redirect_uris_ph: 'https://claude.ai/api/mcp/auth_callback',
  p_ai_mcp_allowed_origins: 'Dopušteni izvori preglednika', p_ai_mcp_allowed_origins_d: 'Web stranice (osim WRM-a) koje smiju pozivati /mcp iz preglednika (CORS). Prazno: nijedna.', p_ai_mcp_allowed_origins_ph: 'https://app.example.com',
  p_ai_mcp_public_url: 'Javna adresa', p_ai_mcp_public_url_d: 'Adresa kojom AI aplikacije pristupaju WRM-u, kad se razlikuje od one koju WRM vidi (obrnuti proxy).', p_ai_mcp_public_url_ph: 'https://wrm.example.com',
  p_ai_mcp_elicitation: 'Odobrenja u AI aplikaciji', p_ai_mcp_elicitation_d: 'Traži odobrenja i u AI aplikaciji (MCP elicitation) kad je podržava. Odobrenja u WRM-u uvijek rade.',
  p_ai_mcp_transfer_max_mb: 'Ograničenje prijenosa (MB)', p_ai_mcp_transfer_max_mb_d: 'Najveća datoteka koju opseg transfer kopira.',
});

(function () {
  const $ = id => document.getElementById(id);
  const tf = (k, vars) => { let s = t(k); for (const [a, b] of Object.entries(vars || {})) s = s.split('{' + a + '}').join(String(b)); return s; };
  async function api(method, url, body) {
    const res = await fetch(url, {method, headers: body ? {'Content-Type': 'application/json'} : {}, body: body ? JSON.stringify(body) : undefined});
    if (!res.ok) throw new Error(await apiError(res));
    return res.json().catch(() => ({}));
  }
  const css = `
.mcp-form .field { margin: 8px 0; }
.mcp-form .field > label { font-size: 11px; color: var(--text2); display: block; margin-bottom: 3px; }
.mcp-form input[type=text], .mcp-form input[type=number], .mcp-form select, .mcp-form textarea { width: 100%; box-sizing: border-box; }
.mcp-form textarea { min-height: 60px; font: 11px 'JetBrains Mono', monospace; }
.mcp-row { display: flex; gap: 10px; } .mcp-row > * { flex: 1; min-width: 0; }
.mcp-conns { max-height: 170px; overflow: auto; border: 1px solid var(--border2); border-radius: 8px; padding: 4px 8px; background: var(--bg); }
.mcp-check { display: flex; align-items: flex-start; gap: 8px; font-size: 12px; padding: 3px 0; color: var(--text2); cursor: pointer; }
.mcp-check input { margin-top: 2px; } .mcp-check .sub { font-size: 10.5px; color: var(--text3); }
.mcp-form .field label.mcp-check, .mcp-form label.mcp-check { display: flex; text-transform: none; letter-spacing: normal; font-size: 12px; font-weight: 400; margin: 0; color: var(--text2); }
.mcp-form .field .mcp-check input, .mcp-form .mcp-check input { width: auto; padding: 0; flex: none; }
.mcp-card input[type=text], .mcp-card textarea { width: 100%; padding: 6px 9px; background: var(--bg3); border: 1px solid var(--border2); color: var(--text); border-radius: 6px; }
.mcp-grant { border: 1px solid var(--border2); border-radius: 8px; padding: 9px 11px; margin: 6px 0; background: var(--bg); display: flex; gap: 10px; align-items: flex-start; }
.mcp-grant .g-main { flex: 1; min-width: 0; font-size: 12px; } .mcp-grant .g-main b { color: var(--text); }
.mcp-grant .sub { font-size: 10.5px; color: var(--text3); margin-top: 3px; word-break: break-word; }
.mcp-snip { position: relative; } .mcp-snip pre { background: var(--bg); border: 1px solid var(--border2); border-radius: 8px; padding: 9px 11px; font: 11px 'JetBrains Mono', monospace; white-space: pre-wrap; word-break: break-all; margin: 6px 0; max-height: 220px; overflow: auto; }
.mcp-snip .btn-sm { position: absolute; top: 6px; right: 6px; }
.mcp-tabs { display: flex; flex-wrap: wrap; gap: 4px; margin: 8px 0 4px; } .mcp-tabs button { font-size: 11px; padding: 4px 9px; }
.mcp-tabs button.active { background: var(--accent); color: #fff; border-color: var(--accent); }
#mcp-cards { position: fixed; right: 16px; bottom: 16px; z-index: 9700; display: flex; flex-direction: column; gap: 10px; max-width: min(460px, calc(100vw - 32px)); max-height: calc(100vh - 32px); overflow: auto; }
.mcp-card { background: var(--bg2); border: 1px solid var(--yellow, #f59e0b); border-radius: 10px; padding: 12px 14px; box-shadow: 0 10px 30px rgba(0,0,0,.45); font-size: 12px; color: var(--text2); }
.mcp-card h4 { margin: 0 0 8px; font-size: 13px; color: var(--text); } .mcp-card pre { background: var(--bg); border: 1px solid var(--border2); border-radius: 6px; padding: 7px 9px; font: 11px 'JetBrains Mono', monospace; white-space: pre-wrap; word-break: break-all; max-height: 200px; overflow: auto; margin: 4px 0 8px; }
.mcp-card pre .a { color: #86efac; } .mcp-card pre .d { color: #fca5a5; } .mcp-card pre .h { color: var(--text3); }
.mcp-card > .sub { display: block; margin-top: 3px; font-size: 11px; color: var(--text3); }
.mcp-card .acts { display: flex; gap: 6px; flex-wrap: wrap; margin-top: 8px; } .mcp-card textarea, .mcp-card input[type=text] { width: 100%; box-sizing: border-box; font: 11px 'JetBrains Mono', monospace; }
.mcp-card textarea { min-height: 70px; }
#mcp-pop { position: fixed; z-index: 9650; top: 52px; right: 12px; width: min(420px, calc(100vw - 24px)); background: var(--bg2); border: 1px solid var(--border2); border-radius: 10px; padding: 12px 14px; box-shadow: 0 12px 34px rgba(0,0,0,.5); display: none; }
#mcp-pop.open { display: block; }
@media (max-width: 480px) { .mcp-row { flex-direction: column; gap: 0; } #mcp-cards { right: 8px; left: 8px; bottom: 8px; max-width: none; } }
`;
  if (!$('mcp-style')) { const s = document.createElement('style'); s.id = 'mcp-style'; s.textContent = css; document.head.appendChild(s); }

  let cfg = null;
  function rel(iso) {
    const d = new Date(iso);
    if (isNaN(d)) return '';
    return fmtTime(iso);
  }
  function diffHTML(d) {
    return String(d || '').split('\n').map(l => {
      const c = l.startsWith('@@') ? 'h' : l.startsWith('+') && !l.startsWith('+++') ? 'a' : l.startsWith('-') && !l.startsWith('---') ? 'd' : '';
      return c ? `<span class="${c}">${esc(l)}</span>` : esc(l);
    }).join('\n');
  }

  // ── the grant form (settings and consent) ─────────────
  function formHTML(prefix, c, pre) {
    pre = pre || {};
    const scopes = pre.scopes && pre.scopes.length ? pre.scopes : c.scopes.filter(s => s === 'read_logs' || s === 'run_readonly');
    const conns = c.connections.map(x => `<label class="mcp-check" data-name="${esc((x.name + ' ' + x.host).toLowerCase())}"><input type="checkbox" class="${prefix}-conn" value="${x.id}">
      <span>${esc(x.name)} <span class="sub">${esc(x.host)}${x.modes && x.modes.length ? '' : ' · ' + esc(t('mcp_off'))}</span></span></label>`).join('');
    return `<div class="mcp-form">
      ${pre.noName ? '' : `<div class="field"><label>${esc(t('mcp_name'))}</label><input type="text" id="${prefix}-name" maxlength="100" placeholder="${esc(t('mcp_name_ph'))}"></div>`}
      <div class="field"><label>${esc(t('mcp_conns'))}</label>${c.connections.length > 8 ? `<input type="text" id="${prefix}-filter" placeholder="${esc(t('mcp_conns_filter'))}" style="margin-bottom:4px;" oninput="wrmMCP.filter('${prefix}')">` : ''}
        <div class="mcp-conns" id="${prefix}-conns">${conns || `<div class="sub">${esc(t('none'))}</div>`}</div></div>
      <div class="mcp-row"><div class="field"><label>${esc(t('mcp_mode'))}</label><select id="${prefix}-mode" onchange="wrmMCP.modeChanged('${prefix}')">${c.modes.map(m => `<option value="${m}" ${m === 'read_only' ? 'selected' : ''}>${esc(t('mcp_mode_' + m))}</option>`).join('')}</select></div>
        <div class="field"><label>${esc(t('mcp_hours'))} <span class="sub">(${esc(tf('mcp_hours_max', {n: c.max_hours}))})</span></label><input type="number" id="${prefix}-hours" min="1" max="${c.max_hours}" value="${c.default_hours}"></div></div>
      <div class="sub" id="${prefix}-mode-d" style="font-size:11px;color:var(--text3);">${esc(t('mcp_mode_read_only_d'))}</div>
      <div id="${prefix}-auto" style="display:none;">
        <div class="field"><label>${esc(t('mcp_auto_allow'))}</label><textarea id="${prefix}-auto-allow" placeholder="systemctl restart nginx"></textarea></div>
        <div class="field"><label>${esc(t('mcp_auto_path_allow'))}</label><textarea id="${prefix}-auto-paths"></textarea></div>
        <div class="mcp-row"><div class="field"><label>${esc(t('mcp_auto_minutes'))}</label><input type="number" id="${prefix}-auto-min" min="1" max="${c.auto_max_minutes}" value="${Math.min(30, c.auto_max_minutes)}"></div>
          <div class="field"><label>${esc(t('mcp_auto_actions'))}</label><input type="number" id="${prefix}-auto-act" min="1" max="${c.auto_max_actions}" value="${Math.min(20, c.auto_max_actions)}"></div></div>
        <label class="mcp-check"><input type="checkbox" id="${prefix}-auto-ok"><span>${esc(t('mcp_auto_confirm'))}</span></label>
      </div>
      <div class="field"><label>${esc(t('mcp_scopes'))}</label>${c.scopes.map(s => `<label class="mcp-check"><input type="checkbox" class="${prefix}-scope" value="${s}" ${scopes.includes(s) ? 'checked' : ''}>
        <span>${esc(t('mcp_sc_' + s))} <span class="sub mono">${esc(s)}</span></span></label>`).join('')}</div>
    </div>`;
  }
  function filter(prefix) {
    const q = ($(prefix + '-filter').value || '').toLowerCase();
    document.querySelectorAll(`#${prefix}-conns .mcp-check`).forEach(l => { l.style.display = !q || l.dataset.name.includes(q) ? '' : 'none'; });
  }
  function modeChanged(prefix) {
    const m = $(prefix + '-mode').value;
    $(prefix + '-mode-d').textContent = t('mcp_mode_' + m + '_d');
    $(prefix + '-auto').style.display = m === 'auto' ? '' : 'none';
  }
  function formValue(prefix) {
    const lines = id => ($(id).value || '').split('\n').map(x => x.trim()).filter(Boolean);
    const body = {
      conn_ids: [...document.querySelectorAll(`.${prefix}-conn:checked`)].map(x => +x.value),
      mode: $(prefix + '-mode').value, hours: +$(prefix + '-hours').value || 0,
      scopes: [...document.querySelectorAll(`.${prefix}-scope:checked`)].map(x => x.value),
    };
    if ($(prefix + '-name')) body.name = $(prefix + '-name').value.trim();
    if (body.mode === 'auto') {
      body.auto = {allow: lines(prefix + '-auto-allow'), path_allow: lines(prefix + '-auto-paths'), minutes: +$(prefix + '-auto-min').value, max_actions: +$(prefix + '-auto-act').value};
      body.confirm_auto = $(prefix + '-auto-ok').checked;
    }
    return body;
  }

  // ── setup snippets ────────────────────────────────────
  function snippets(url, tok) {
    const base = url.replace(/\/mcp$/, '');
    const out = {
      claude_desktop: {help: t('mcp_h_claude_desktop'), code: JSON.stringify({mcpServers: {wrm: {command: '/path/to/wrm-pro', args: ['-mcp-stdio', '-url', url], env: {WRM_MCP_TOKEN: tok}}}}, null, 2)},
      remote: {help: t('mcp_h_remote'), code: url},
      chatgpt: {help: t('mcp_h_chatgpt'), code: url},
      generic: {help: t('mcp_h_generic'), code: JSON.stringify({mcpServers: {wrm: {type: 'http', url, headers: {Authorization: 'Bearer ' + tok}}}}, null, 2) +
        `\n\nclaude mcp add --transport http wrm ${url} --header "Authorization: Bearer ${tok}"`},
      cli: {help: t('mcp_h_cli'), code: `curl -sS ${url} -H "Authorization: Bearer ${tok}" -H "Content-Type: application/json" -H "Accept: application/json, text/event-stream" \\\n  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"curl"}}}'`},
    };
    void base;
    return out;
  }
  let snipState = {};
  function snippetsHTML(url, tok, which) {
    snipState = {url, tok};
    const all = snippets(url, tok);
    const keys = ['claude_desktop', 'remote', 'chatgpt', 'generic', 'cli'];
    which = which || keys[0];
    return `<div class="adm-h">🧩 ${esc(t('mcp_setup'))}</div><div class="mcp-tabs">${keys.map(k => `<button class="btn-sec btn-sm ${k === which ? 'active' : ''}" onclick="wrmMCP.snip('${k}')">${esc(t('mcp_setup_' + k))}</button>`).join('')}</div>
      <div class="hint-box">${esc(all[which].help)}</div><div class="mcp-snip"><pre id="mcp-snip-code">${esc(all[which].code)}</pre><button class="btn-sec btn-sm" onclick="wrmMCP.copy('mcp-snip-code')">${esc(t('mcp_copy'))}</button></div>`;
  }
  function snip(k) { const box = $('mcp-snips'); if (box) box.innerHTML = snippetsHTML(snipState.url, snipState.tok, k); }
  function copy(id) { const el = $(id); if (el) copyText(el.textContent); }

  // ── Settings → AI connections ─────────────────────────
  async function settingsPane() {
    const pane = $('mcp-pane');
    if (!pane) return;
    let grants;
    try { [cfg, grants] = await Promise.all([api('GET', '/api/mcp/config'), api('GET', '/api/mcp/grants')]); }
    catch (e) { pane.innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    if (!cfg.enabled) { pane.innerHTML = `<div class="hint-box warn">${esc(t('mcp_off'))}: ${esc(cfg.reason)}</div>`; return; }
    pane.innerHTML = `<div class="hint-box">${esc(t('mcp_intro'))}</div>
      <div class="field" style="margin:8px 0;"><label style="font-size:11px;color:var(--text2);">${esc(t('mcp_url'))}</label><div class="mcp-snip"><pre id="mcp-url">${esc(cfg.url)}</pre><button class="btn-sec btn-sm" onclick="wrmMCP.copy('mcp-url')">${esc(t('mcp_copy'))}</button></div>
        <div class="sub" style="font-size:11px;color:var(--text3);">${esc(cfg.oauth ? t('mcp_oauth_box') : t('mcp_oauth_off'))}</div></div>
      <div style="display:flex;align-items:center;gap:8px;"><div class="adm-h" style="flex:1;">🔌 ${esc(t('mcp_title'))}</div>${grants.some(g => g.state === 'active') ? `<button class="btn-danger btn-sm" onclick="wrmMCP.revokeAll()">${esc(t('mcp_revoke_all'))}</button>` : ''}</div>
      <div id="mcp-grants">${grantsHTML(grants)}</div>
      <div class="adm-h">＋ ${esc(t('mcp_new'))}</div>${formHTML('mcpn', cfg)}
      <div style="margin-top:8px;"><button onclick="wrmMCP.create()">${esc(t('mcp_create'))}</button></div>
      <div id="mcp-created"></div>`;
  }
  function grantsHTML(list) {
    if (!list.length) return `<div class="sub" style="font-size:11.5px;color:var(--text3);">${esc(t('mcp_none'))}</div>`;
    return list.map(g => {
      const st = g.state === 'active' ? (g.live ? `<span class="pill ok">● ${esc(t('mcp_live'))}</span>` : '') : `<span class="pill bad">${esc(t('mcp_' + g.state))}</span>`;
      return `<div class="mcp-grant"><div class="g-main"><b>${esc(g.name)}</b> <span class="pill">${esc(t('mcp_kind_' + g.kind))}</span> ${st}
        <div class="sub">${esc(g.connections.map(c => c.name).join(', ') || '—')} · ${esc(t('mcp_mode_' + g.mode))} · <span class="mono">${esc(g.scopes.join(', '))}</span></div>
        <div class="sub">${g.state === 'active' ? esc(tf('mcp_expires', {time: rel(g.expires_at)})) : ''}${g.last_used_at ? ' · ' + esc(tf('mcp_last_used', {time: rel(g.last_used_at)})) : ' · ' + esc(t('mcp_never_used'))}${g.hint ? ` · …${esc(g.hint)}` : ''}${g.sessions.length ? ' · ' + esc(tf('mcp_sessions_n', {n: g.sessions.length})) : ''}</div></div>
        ${g.state === 'active' ? `<button class="btn-danger btn-sm" onclick="wrmMCP.revoke(${g.id})">${esc(t('mcp_revoke'))}</button>` : ''}</div>`;
    }).join('');
  }
  async function create() {
    const body = formValue('mcpn');
    let r;
    try { r = await api('POST', '/api/mcp/grants', body); } catch (e) { showToast(e.message, 'error', 6000); return; }
    await settingsPane();
    $('mcp-created').innerHTML = `<div class="hint-box warn" style="margin-top:12px;">${esc(t('mcp_token_once'))}</div>
      <div class="mcp-snip"><pre id="mcp-token">${esc(r.token)}</pre><button class="btn-sec btn-sm" onclick="wrmMCP.copy('mcp-token')">${esc(t('mcp_copy'))}</button></div>
      <div id="mcp-snips">${snippetsHTML(r.url, r.token)}</div>`;
    $('mcp-created').scrollIntoView({block: 'nearest'});
  }
  async function revoke(id) {
    if (!(await uiConfirm(t('mcp_revoke_q'), {okText: t('mcp_revoke'), danger: true}))) return;
    try { await api('DELETE', `/api/mcp/grants/${id}`); } catch (e) { showToast(e.message, 'error'); }
    after();
  }
  async function revokeAll() {
    if (!(await uiConfirm(t('mcp_revoke_all_q'), {okText: t('mcp_revoke_all'), danger: true}))) return;
    try { const r = await api('POST', '/api/mcp/revoke-all'); showToast(tf('mcp_revoked_n', {n: r.revoked}), 'success'); } catch (e) { showToast(e.message, 'error'); }
    after();
  }
  function after() {
    if ($('mcp-pane') && $('st-mcp') && $('st-mcp').classList.contains('active')) settingsPane();
    if (typeof adminTab !== 'undefined' && adminTab === 'mcp' && $('admin-body')) adminMCP();
    refresh();
  }

  // ── the "AI connected" indicator and approval cards ───
  let lastActive = {live: [], pending: []};
  async function refresh() {
    let d;
    try { d = await api('GET', '/api/mcp/active'); } catch (_) { return; }
    indicator(d);
  }
  function indicator(d) {
    lastActive = d;
    const btn = $('btn-mcp');
    if (btn) {
      btn.style.display = d.live.length || d.pending.length ? '' : 'none';
      const bar = $('top-bar'); if (bar) bar.classList.toggle('mcp-on', !!(d.live.length || d.pending.length));
      const b = $('mcp-badge');
      if (b) { b.textContent = d.pending.length || d.live.length; b.style.background = d.pending.length ? 'var(--yellow)' : ''; }
    }
    if ($('mcp-pop') && $('mcp-pop').classList.contains('open')) renderPop();
    cards(d.pending);
  }
  function renderPop() {
    const d = lastActive;
    $('mcp-pop').innerHTML = `<div style="display:flex;align-items:center;gap:8px;margin-bottom:6px;"><b style="flex:1;font-size:13px;">🤖 ${esc(t('mcp_ind_title'))}</b>
      ${d.live.length ? `<button class="btn-danger btn-sm" onclick="wrmMCP.revokeAll()">⏻ ${esc(t('mcp_revoke_all'))}</button>` : ''}<button class="icon-btn" onclick="wrmMCP.togglePop(false)">✕</button></div>` +
      (d.live.length ? d.live.map(g => `<div class="mcp-grant"><div class="g-main"><b>${esc(g.client_name || g.name)}</b> <span class="pill ok">● ${esc(t('mcp_live'))}</span>
        <div class="sub">${esc(tf('mcp_on', {servers: (g.sessions.length ? g.sessions.map(s => s.connection) : g.connections.map(c => c.name)).join(', ')}))} · ${esc(t('mcp_mode_' + g.mode))}</div>
        <div class="sub">${esc(tf('mcp_expires', {time: rel(g.expires_at)}))}</div></div><button class="btn-danger btn-sm" onclick="wrmMCP.revoke(${g.id})">${esc(t('mcp_revoke'))}</button></div>`).join('')
        : `<div class="sub" style="font-size:11.5px;color:var(--text3);">${esc(t('mcp_none'))}</div>`) +
      `<div style="margin-top:8px;"><button class="btn-sec btn-sm" onclick="wrmMCP.togglePop(false); openSettings('st-mcp')">⚙ ${esc(t('mcp_title'))}</button></div>`;
  }
  function togglePop(open) {
    let p = $('mcp-pop');
    if (!p) { p = document.createElement('div'); p.id = 'mcp-pop'; document.body.appendChild(p); }
    if (open === undefined) open = !p.classList.contains('open');
    p.classList.toggle('open', open);
    if (open) { renderPop(); refresh(); }
  }
  const shown = {};
  let cardTimer = null;
  function cards(pending) {
    let box = $('mcp-cards');
    if (!box) { box = document.createElement('div'); box.id = 'mcp-cards'; document.body.appendChild(box); }
    const ids = new Set(pending.map(p => p.approval.id));
    for (const id of Object.keys(shown)) if (!ids.has(id)) { const el = $('mcp-card-' + id); if (el) el.remove(); delete shown[id]; }
    for (const p of pending) {
      const a = p.approval;
      if (shown[a.id]) continue;
      shown[a.id] = p;
      const el = document.createElement('div');
      el.className = 'mcp-card'; el.id = 'mcp-card-' + a.id;
      const what = a.transfer ? `<div>${esc(t('mcp_ap_transfer'))}</div><pre>${esc(a.transfer)}</pre>`
        : a.command ? `<div>${esc(t('mcp_ap_cmd'))}</div><pre id="mcp-cmd-${a.id}">${esc(a.command)}</pre>`
        : `<div>${esc(tf('mcp_ap_edit_file', {path: a.path}))}</div><pre>${diffHTML(a.diff)}</pre>`;
      el.innerHTML = `<h4>🤖 ${esc(tf('mcp_ap_title', {client: p.client, server: p.connection}))}</h4>${what}
        ${a.why ? `<div class="sub">${esc(t('mcp_ap_why'))}: ${esc(a.why)}</div>` : ''}${a.reason ? `<div class="sub">${esc(t('mcp_ap_reason'))}: ${esc(a.reason)}</div>` : ''}
        <div id="mcp-edit-${a.id}"></div>
        <input type="text" id="mcp-note-${a.id}" maxlength="500" placeholder="${esc(t('mcp_ap_note'))}" style="margin-top:6px;">
        <div class="acts"><button class="btn-sm" onclick="wrmMCP.decide(${p.session}, '${a.id}', true)">✔ ${esc(t('mcp_ap_approve'))}</button>
          <button class="btn-danger btn-sm" onclick="wrmMCP.decide(${p.session}, '${a.id}', false)">✕ ${esc(t('mcp_ap_deny'))}</button>
          ${a.transfer ? '' : `<button class="btn-sec btn-sm" onclick="wrmMCP.edit('${a.id}')">✎ ${esc(t('mcp_ap_edit'))}</button>`}
          <span class="sub" style="margin-left:auto;align-self:center;" data-exp="${esc(a.expires)}"></span></div>`;
      box.appendChild(el);
    }
    if (!cardTimer && Object.keys(shown).length) cardTimer = setInterval(tickCards, 1000);
    tickCards();
  }
  function tickCards() {
    const els = document.querySelectorAll('#mcp-cards [data-exp]');
    if (!els.length && cardTimer) { clearInterval(cardTimer); cardTimer = null; return; }
    els.forEach(el => { const s = Math.max(0, Math.round((new Date(el.dataset.exp) - Date.now()) / 1000)); el.textContent = tf('mcp_ap_expires', {s}); });
  }
  function edit(aid) {
    const p = shown[aid];
    if (!p) return;
    const a = p.approval, box = $('mcp-edit-' + aid);
    if (box.innerHTML) return;
    box.innerHTML = a.command ? `<input type="text" id="mcp-edv-${aid}" value="${esc(a.command)}">` : `<textarea id="mcp-edv-${aid}">${esc(a.content || '')}</textarea>`;
    const acts = box.parentElement.querySelector('.acts button');
    acts.textContent = '✔ ' + t('mcp_ap_approve_edited');
  }
  async function decide(sid, aid, approve) {
    const p = shown[aid];
    const body = {decision: approve ? 'approve' : 'deny', note: ($('mcp-note-' + aid) || {}).value || ''};
    const ed = $('mcp-edv-' + aid);
    if (approve && ed && p) { if (p.approval.command) body.command = ed.value; else body.content = ed.value; }
    try { await api('POST', `/api/ai/sessions/${sid}/approvals/${aid}`, body); } catch (e) { showToast(e.message, 'error', 6000); }
    const el = $('mcp-card-' + aid); if (el) el.remove();
    delete shown[aid];
    refresh();
  }

  // ── OAuth consent ─────────────────────────────────────
  async function consent(id) {
    let req, c;
    try { [req, c] = await Promise.all([api('GET', `/api/mcp/authorize/${encodeURIComponent(id)}`), api('GET', '/api/mcp/config')]); }
    catch (e) { showModal(`<h3>🤖 ${esc(t('mcp_cs_title'))}</h3><div class="hint-box warn">${esc(e.message || t('mcp_cs_expired'))}</div><div class="modal-footer"><button onclick="wrmMCP.closeModal()">${esc(t('close'))}</button></div>`); clearParam(); return; }
    cfg = c;
    if (!req.enabled) { showModal(`<h3>🤖 ${esc(t('mcp_cs_title'))}</h3><div class="hint-box warn">${esc(t('mcp_off'))}: ${esc(req.reason)}</div><div class="modal-footer"><button onclick="wrmMCP.answer('${esc(id)}', false)">${esc(t('mcp_cs_deny'))}</button></div>`); return; }
    showModal(`<h3>🤖 ${esc(t('mcp_cs_title'))}</h3>
      <div class="hint-box"><b>${esc(tf('mcp_cs_asks', {client: req.client_name}))}</b><br>${esc(tf('mcp_cs_redirect', {host: req.redirect_host}))}</div>
      ${formHTML('mcpc', c, {scopes: (req.scopes || []).filter(s => c.scopes.includes(s)), noName: true})}
      <div class="modal-footer"><button class="btn-sec" onclick="wrmMCP.answer('${esc(id)}', false)">${esc(t('mcp_cs_deny'))}</button><button onclick="wrmMCP.answer('${esc(id)}', true)">${esc(t('mcp_cs_approve'))}</button></div>`);
  }
  async function answer(id, approve) {
    const body = approve ? Object.assign({decision: 'approve'}, formValue('mcpc')) : {decision: 'deny'};
    let r;
    try { r = await api('POST', `/api/mcp/authorize/${encodeURIComponent(id)}`, body); } catch (e) { showToast(e.message, 'error', 6000); return; }
    clearParam();
    $('mcp-modal').querySelector('.modal').innerHTML = `<h3>🤖 ${esc(t('mcp_cs_title'))}</h3><div class="hint-box">${esc(t('mcp_cs_done'))}</div>`;
    location.href = r.redirect;
  }
  function clearParam() {
    try { const u = new URL(location.href); u.searchParams.delete('mcp_authorize'); history.replaceState(null, '', u.pathname + u.search + u.hash); } catch (_) {}
  }
  function showModal(html) {
    let m = $('mcp-modal');
    if (!m) { m = document.createElement('div'); m.id = 'mcp-modal'; m.className = 'modal-overlay'; m.style.zIndex = 9650; document.body.appendChild(m); }
    m.innerHTML = `<div class="modal" style="width:560px;max-width:calc(100vw - 32px);max-height:calc(100vh - 32px);overflow:auto;">${html}</div>`;
    m.classList.add('open');
  }
  function closeModal() { const m = $('mcp-modal'); if (m) m.classList.remove('open'); }

  // ── Admin panel → AI connections ──────────────────────
  const MCP_POLICY_KEYS = ['ai_mcp_enabled', 'ai_mcp', 'ai_mcp_modes', 'ai_mcp_scopes', 'ai_mcp_token_hours', 'ai_mcp_max_token_hours', 'ai_mcp_oauth',
    'ai_mcp_elicitation', 'ai_mcp_redirect_uris', 'ai_mcp_allowed_origins', 'ai_mcp_public_url', 'ai_mcp_transfer_max_mb'];
  async function adminMCP() {
    const body = $('admin-body');
    let grants, clients, settings, c;
    try {
      [grants, clients, settings, c] = await Promise.all([api('GET', '/api/admin/mcp/grants'), api('GET', '/api/admin/mcp/clients'), api('GET', '/api/admin/settings'), api('GET', '/api/mcp/config')]);
    } catch (e) { body.innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    adminSettingsCache = settings;
    const active = grants.filter(g => g.state === 'active');
    body.innerHTML = `<div class="hint-box">${esc(tf('mcp_adm_help', {url: c.url}))}</div>
      <div style="display:flex;align-items:center;gap:8px;"><div class="adm-h" style="flex:1;">🔌 ${esc(t('mcp_adm_active'))}</div>${active.length ? `<button class="btn-danger btn-sm" onclick="wrmMCP.adminRevokeAll()">⏻ ${esc(t('mcp_revoke_all'))}</button>` : ''}</div>
      ${active.length ? `<div class="tbl-scroll"><table class="admin-table"><thead><tr><th>${esc(t('mcp_user'))}</th><th>${esc(t('mcp_client'))}</th><th>${esc(t('mcp_conns'))}</th><th>${esc(t('mcp_mode'))}</th><th>${esc(t('mcp_state'))}</th><th></th></tr></thead><tbody>` +
        active.map(g => `<tr><td>${esc(g.user)}</td><td>${esc(g.client_name || g.name)} <span class="pill">${esc(t('mcp_kind_' + g.kind))}</span><div class="sub">${esc(tf('mcp_expires', {time: rel(g.expires_at)}))}</div></td>
          <td>${esc(g.connections.map(x => x.name).join(', '))}<div class="sub mono">${esc(g.scopes.join(', '))}</div></td><td>${esc(t('mcp_mode_' + g.mode))}</td>
          <td>${g.live ? `<span class="pill ok">● ${esc(t('mcp_live'))}</span>` : ''}<div class="sub">${g.last_used_at ? esc(tf('mcp_last_used', {time: rel(g.last_used_at)})) : esc(t('mcp_never_used'))}${g.last_ip ? ' · ' + esc(g.last_ip) : ''}</div></td>
          <td><button class="btn-danger btn-sm" onclick="wrmMCP.adminRevoke(${g.id})">${esc(t('mcp_revoke'))}</button></td></tr>`).join('') + '</tbody></table></div>'
        : `<div class="chat-empty">${esc(t('mcp_adm_none'))}</div>`}
      <div class="adm-h">⚙ ${esc(t('mcp_adm_policies'))}</div>
      <div class="policy-group">${MCP_POLICY_KEYS.map(k => settingRowHTML(adminSettingsCache.find(x => x.key === k))).join('')}</div>
      <div class="adm-h">🧩 ${esc(t('mcp_adm_clients'))}</div>
      ${clients.length ? `<div class="tbl-scroll"><table class="admin-table"><thead><tr><th>${esc(t('mcp_client'))}</th><th>${esc(t('mcp_adm_redirects'))}</th><th>${esc(t('mcp_created'))}</th><th></th></tr></thead><tbody>` +
        clients.map(x => `<tr><td>${esc(x.client_name)}<div class="sub mono">${esc(x.client_id)}</div>${x.active ? `<span class="pill ok">${x.active}</span>` : ''}</td>
          <td class="mono" style="font-size:10.5px;word-break:break-all;">${esc((x.redirect_uris || []).join(' '))}</td><td>${esc(rel(x.created_at))}<div class="sub">${esc(x.created_ip)}</div></td>
          <td><button class="btn-sec btn-sm" onclick="wrmMCP.adminRemoveClient(${x.id})">${esc(t('mcp_adm_remove'))}</button></td></tr>`).join('') + '</tbody></table></div>'
        : `<div class="chat-empty">${esc(t('none'))}</div>`}`;
  }
  async function adminRevoke(id) {
    if (!(await uiConfirm(t('mcp_revoke_q'), {okText: t('mcp_revoke'), danger: true}))) return;
    try { await api('DELETE', `/api/admin/mcp/grants/${id}`); } catch (e) { showToast(e.message, 'error'); }
    adminMCP();
  }
  async function adminRevokeAll() {
    if (!(await uiConfirm(t('mcp_revoke_all_q'), {okText: t('mcp_revoke_all'), danger: true}))) return;
    try { const r = await api('POST', '/api/admin/mcp/revoke-all'); showToast(tf('mcp_revoked_n', {n: r.revoked}), 'success'); } catch (e) { showToast(e.message, 'error'); }
    adminMCP();
  }
  async function adminRemoveClient(id) {
    if (!(await uiConfirm(t('mcp_adm_remove_q'), {okText: t('mcp_adm_remove'), danger: true}))) return;
    try { await api('DELETE', `/api/admin/mcp/clients/${id}`); } catch (e) { showToast(e.message, 'error'); }
    adminMCP();
  }

  window.wrmMCP = {
    settingsPane, create, revoke, revokeAll, filter, modeChanged, snip, copy, refresh, indicator, togglePop, decide, edit,
    consent, answer, closeModal, adminTab: adminMCP, adminRevoke, adminRevokeAll, adminRemoveClient,
  };
})();
