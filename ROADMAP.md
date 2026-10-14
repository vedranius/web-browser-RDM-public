# Roadmap

Released work is in [CHANGELOG.md](CHANGELOG.md). The owner prioritized the Windows Medium hotfix **v12.4.3b** first. The complete remaining scope below is accepted for **v12.4.4**; it is a plan, not a list of completed features. Original v12.4.1/v12.4.2 request labels are retained only for provenance.

## v12.4.4 — Monitoring, logs, MCP workflows and remaining requests

Status: accepted TODO; implement after publishing v12.4.3b. Medium itself is delivered in v12.4.3 and corrected in v12.4.3b. Do not represent this future scope as delivered by the hotfix.

### Import verification and issue #5

- Review the final comment on [issue #5](https://github.com/vedranius/web-browser-RDM-public/issues/5): PuTTY import and a 4 MB mRemoteNG XML rejected as JSON. The current code already has dedicated mRemoteNG XML and PuTTY registry importers, proxy mapping and large-import regressions. Verify the exact visible UI workflow before calling the report resolved.
- Reproduce a 4 MB `confCons.xml`, including `.xml`/`.config` extensions, master-password handling, encrypted/unencrypted supported formats, nesting, duplicates and proxy/jump links. Verify PuTTY `.reg` UTF-16/UTF-8, protocol/port/username/key/proxy settings and useful unsupported-field warnings.
- Correct format routing/discovery/error guidance if the generic JSON importer still receives XML. Verify size limits through the real browser and reverse-proxy/API paths; preserve encrypted secrets and owner boundaries.
- Link acknowledgement, actual verification/fixes, commits, changelog, PR and v12.4.4 release to issue #5. Close only after the final user's scenario is reproduced and verified.

### Easier AI/MCP onboarding and commands

- Simplify and automate WRM MCP setup for supported AI applications, including local command-line clients and cloud OAuth clients. Provide application-specific detection, actionable setup/install/config generation, connection tests and status. Check current primary vendor documentation before generating commands/configuration.
- Resolve the reported read-only `.ini`/`.py` inventory search denial: reproduce the permitted scope and command/tool parser, offer structured read-only search/inventory tools where possible, and show the exact policy/scope blocker. A rejection must not be bypassed by rewriting a command or granting unrestricted execution.
- Add **Automatic within limits** to MCP permissions, consistent with the internal assistant's supported command/file limits and approval/audit behavior. Keep read-only, approval and bounded automatic modes distinct; do not equate automatic with unrestricted access.
- Add `/wrm` onboarding/help and `/wrm server` selection. Expose supported MCP prompts/actions and WRM assistant commands; where a client can render a picker/dropdown, show only authorized servers with hostname/folder context. For clients without custom slash/dropdown support, provide a clear supported prompt/text/picker-link fallback. Verify per-client behavior instead of claiming a global command registration.
- Cloud connectors need a reachable HTTPS/OAuth endpoint; local `localhost` instructions are appropriate only to clients on that host. Generate the correct URL/transport for the chosen client and deployment, with explicit setup/authorization for publishing access.
- Allow existing MCP connections/tokens to edit authorized servers and permission mode/limits without revoke/recreate. Owner/admin validation and checks on every call apply immediately, including active sessions; removed/revoked scopes cannot keep cached access. Changes are audited without exposing bearer tokens.
- Validate the screenshot's wizard at mobile widths, including checkbox/radio alignment, server search, selected scope, an automatic mode and URL guidance. Never publish its real server names/addresses.

### Monitoring workspace and server metrics

- Add **Monitoring** as a primary workspace and move service Functionality inspection into it with preserved links/configuration. Keep Git/catalog service associations and diagnostics reachable from Git. Group by hostname/connection with duplicate-name context and the existing planned searches/filters.
- Provide a Proxmox-style overview: CPU, memory, load/uptime, filesystem capacity, network rates and server/agent availability, per-server detail/history/charts, freshness and healthy/warning/failed/unknown states. Missing metrics are unknown, not zero. Metrics describe this server; they do not imply hypervisor features on arbitrary hosts.
- Use installed Zabbix agent/agent 2 where feasible. Evaluate passive polling over an authorized network/SSH route with supported protocol/keys and TLS/PSK/certificate settings; an active-only agent needs the existing Zabbix server/API path or an explicitly configured receiver. Explain which transport is supported and test actual agent versions/platforms.
- Offer agent detection/setup guidance and optional administrator-approved installation/configuration after reviewing the target OS, permissions, allowed hosts and TLS. Do not silently open agent ports, change existing Zabbix configuration or enable arbitrary `system.run` commands.
- Bound polling/concurrency/retention, encrypt connection secrets, enforce ownership/share scope and audit configuration. Problem/recovery notifications should deduplicate states and report unavailable/stale checks.

Feasibility references: [Zabbix passive/active agent checks](https://www.zabbix.com/documentation/current/en/manual/appendix/items/activepassive), [supported agent items](https://www.zabbix.com/documentation/current/en/manual/config/items/itemtypes/zabbix_agent). Implementation must distinguish passive request/response from active push; direct polling is not sufficient for an active-only installation.

### Logs discovery, search and live windows

- Add a **Logs** tab under Monitoring. Configure multiple owned log roots, including `/var/log` and service-specific directories such as `/srv/demo/service/log`. Discover readable logs and map them to server/service/operator. Keep permission errors, missing roots, rotations, caps and truncated results explicit.
- Combine server/hostname, service, operator, log selection, optional date/time range and keyword filters. Search/filter controls stay at the top of the actual scroll container and preserve selection/search across refreshes.
- Highlight keyword matches, show next/previous match and counts, support copy/download and open the selected log in a terminal. Treat search text and filenames as data; bounded remote reads, safe roots/symlinks and existing access/share permissions apply.
- Parse timestamp formats/timezones explicitly; do not discard lines silently when timestamps cannot be parsed. Handle large files, multiline events, binary/compressed/unsupported files and offsets without returning unlimited content.
- Select multiple logs and follow them concurrently in independently movable/resizable windows. Provide vertical/horizontal tile and cascade layouts, optionally showing the latest 5–10 lines per panel. Each title includes operator, hostname/server and log path/name.
- Follow real-time append/rotation/truncation with pause/resume, autoscroll control, reconnect/gap indicators, bounded buffers and cancellation when a panel closes. Browser refresh/reconnect must not spawn duplicate tails. Respect scope changes and server disconnection immediately.

### Terminal selection, AI and context menu

- Add **Send selection to AI** to terminal/context actions for selected text or a confirmed file/path. Preserve server/share identity and current directory, preview/identify the selected context, apply provider/account policy and bounded read limits, and treat selected text as untrusted data rather than an executable instruction.
- Preserve normal Copy/Paste and related applicable browser/terminal actions in custom right-click menus. Keep keyboard shortcuts and a way to open the native browser menu; do not intercept editor/input context menus unnecessarily. Clipboard restrictions need a clear fallback.
- Add **Open terminal here** from the current terminal directory or a verified selected directory, and applicable file-manager/log contexts. Use the same authorized server/route/credentials/share and start the new terminal in the safely quoted directory. Do not execute text from a selection as shell code; handle remote shell/platform differences.

### Combined acceptance and release

- Test the latest issue #5 scenario in a real browser, MCP search denials/allowed inventory, all permission modes, mutable token scopes, supported client installation/onboarding and `/wrm` selection behavior. Scope changes must invalidate unauthorized pending/active tool actions.
- Exercise real Zabbix agent/agent 2 fixtures, passive/active distinctions, authorization/TLS errors, unavailable and stale metrics, bounded retention and duplicate hostnames. Do not claim production agents were tested without access.
- Exercise logs with multiple roots/operators, rotations, truncation, huge files, timestamp variants/timezones, keyword navigation and several live tiled/cascade windows. Verify cancellation, reconnect, buffer limits and permissions.
- Verify selection-to-AI and Open terminal here with spaces, quotes, malicious-looking filenames, shares and different shells. Verify native copy/paste and fallback behavior.
- Apply en/hr translations, additive migrations, full race/vet/format/UI tests and real-browser checks at 390/1100/1440/1920 px. Link the complete accepted scope, implementation commits, PR, changelog and the published v12.4.4 release. Remove this roadmap section only when delivered.

### Client GitLab access/cache and replay (carried from the original v12.4.1 request)

Status: accepted backlog; implementation scheduled for v12.4.4.

#### GitLab access through the user's computer

- Add an optional client access mode for a configured GitLab source. The user's computer authenticates to the GitLab API with an access token, fetches repository data and forwards it to WRM over its authenticated WRM session, allowing WRM to operate when its server cannot reach GitLab.
- Provide a client-side GitLab URL/access-token entry and repository access check. Use the token for API authorization on that computer; the computer's domain/VPN connection supplies network reachability. Keep the client token in session memory by default, separate from the repository cache, and never export existing server-stored tokens to the browser.
- Show the access mode, client connection status, synchronization progress and actionable connection/authentication errors in Git settings and the workspace. Preserve the existing direct server mode.
- On connection of an authorized user with client synchronization enabled, check access and synchronize the configured catalog/services. A subsequent user's connection or manual **Check / Refresh** starts a new synchronization.
- Cover projects required by the catalog, branches, tags, resolved commits, repository trees and file contents needed for comparison, diffs and deployment. Declare size/pagination limits, missing content and unsupported operations explicitly; never label a partial snapshot complete.
- Distinguish unreachable GitLab, rejected/expired tokens and insufficient repository/API permissions. Check token permissions against the API operations used by synchronization; use read access for repository synchronization and separate permission for writes.
- Check browser constraints for token-authenticated requests before choosing the transport: CORS/preflight for token headers, permitted WRM origin, HTTPS, WRM CSP and local-network permissions. Access tokens do not remove these transport constraints. If plain browser requests are blocked, provide a documented client bridge through a narrowly scoped extension or companion on that computer. Report which transport is required; do not require browser security to be disabled.
- Bind relay requests to the authorized GitLab source and synchronization job. Validate endpoints, responses, sizes and repository/blob identities on WRM; apply timeouts, cancellation and concurrency limits. Keep authentication material out of cache payloads and logs.
- Identify operations that require a live client, including uncached repository reads and GitLab writes. Disable unavailable operations with an explanation; retain the existing review/approval checks for deployments and writes.

#### Persistent cache and refresh ownership

- Keep the latest complete successful synchronization on the WRM server, including after the contributing browser disconnects or WRM restarts. Use it until a newer successful check replaces it.
- Display the source, snapshot generation, last successful refresh, contributing user and cache completeness. Distinguish **using cached data**, **refreshing**, **refresh failed** and **live** states.
- Prepare each refresh in isolation and publish its snapshot atomically. A failed, canceled or interrupted refresh keeps the previous usable snapshot. An older synchronization finishing late must not overwrite a newer successful snapshot.
- Reuse immutable trees/blobs and transfer only missing or changed data. Apply storage limits and retention to obsolete snapshots without removing content still referenced by the active snapshot or a reviewed deployment.
- Add an explicit shared Git source/cache scope for teams that need the newest authorized user's check to refresh the common view. Current Git sources are per user: connecting another user must not silently expose that user's private repositories or change other users' credentials. Enforce scope membership and source permissions on every cache read and refresh.
- Keep deployment plans tied to their reviewed commit and content; a background refresh cannot silently change an approved target. Surface stale cache and missing content before comparison/deployment, and audit successful/failed refreshes without secrets.

#### Replay window scaling

- Audit the existing terminal `replayFit`, `ResizeObserver`, fullscreen layout and graphical replay scaling. Automatic scaling already exists; reproduce any remaining horizontal overflow before changing it.
- Make the complete recorded frame fit the available window width without left/right scrolling, preserving its aspect ratio and recorded terminal dimensions.
- Recalculate after opening, seeking, recorded resize events, font loading, modal/browser resize, mobile rotation and fullscreen changes. Review minimum scale, flex sizing, padding/borders and transformed-element overflow.
- Keep replay controls accessible on narrow screens and verify both terminal recordings and graphical desktop recordings.

#### Acceptance and release

- Validate with a token-authenticated GitLab fixture reachable from the client while blocked from WRM: connect, verify repository access, synchronize, disconnect, restart WRM, use cached views/comparisons and refresh from another authorized client with its own token.
- Cover successful/failed/interrupted refreshes, concurrent clients finishing out of order, changed/deleted refs, pagination, large/missing blobs, invalid/expired tokens, insufficient token permissions, token-header preflight/CORS/local-network rejection and isolated/shared permissions. Verify that tokens never appear in cached data, API responses or logs.
- Confirm a cache refresh does not alter a reviewed deployment and that unavailable live operations explain the required client connection.
- Replay recordings with very wide terminals/desktops and changing dimensions at 390, 1100, 1440 and 1920 px, with modal resizing and fullscreen. Assert no horizontal overflow and inspect screenshots.
- During implementation, run the repository's Go race, vet/format and UI checks, plus browser validation. Update en/hr text, affected architecture/security/API documentation, changelog and release notes for v12.4.4; link implementation commits, PR and published release to the accepted scope.

Browser references: [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS), [Fetch credentials](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch), [local network access](https://developer.mozilla.org/en-US/docs/Web/Security/Defenses/Local_network_access).

### Long AI requests, navigation and visual labels (carried from the original v12.4.2 request)

Status: accepted for v12.4.4; implementation pending. The original Medium integration was moved forward and delivered in v12.4.3 using the supplied reference archives.

#### Ollama timeout on longer AI requests

- Investigate the reported `net/http: timeout awaiting response headers` from the OpenAI-compatible Ollama `/v1/chat/completions` endpoint during longer requests. The current AI HTTP transport has a fixed three-minute response-header timeout and a separate ten-minute request limit.
- Reproduce slow model loading, long prompt evaluation, queued requests and long streamed responses. Inspect header wait, overall request deadline, stream handling and any intervening proxy limits separately.
- Add appropriate configurable provider/request timeouts with documented defaults and bounds so valid long requests can finish. Keep cancellation responsive and distinguish provider unreachability, response-header timeout, total request timeout and interrupted streaming in the UI.
- Show that the assistant is still waiting/processing and preserve its diagnostic context without exposing credentials. Review retry behavior so a timeout cannot silently repeat a tool action or execute a generated command twice.
- Verify local Ollama/OpenAI-compatible access as well as other providers; longer timeout settings must not break stop/cancel, approvals or tool execution limits.

#### Functionality organized by server hostname

- Make server hostname the primary grouping of Functionality, with connection name and service/install details under each server. Distinguish identically named hosts using connection identity/context; use connection details as an explicit fallback when hostname is unknown.
- Add searchable hostname/connection, service/repository, environment, installation path, unit/worker name and warning/error text.
- Add combinable filters for server, service, environment, runtime kind, health state, installation problems and never-checked/stale results. Include quick views for failures/warnings and a clear reset action.
- Show per-server health totals and last-check time, collapsible groups, expand/collapse controls, matching-result counts and a useful empty state. Keep selected filters when results refresh.
- Keep Supervisor/systemd/resource findings and their installations together under the appropriate server; configuration and error-log actions must remain easy to reach.

#### Search and filters stay visible while scrolling

- Keep the search/filter controls at the top of the scrolling view in Git Overview, Installations, Functionality and Activity.
- Apply the same behavior wherever WRM has a search/filter area, including Simple/Medium/Advanced and settings/admin lists as applicable. Inventory the affected screens so this is applied consistently.
- Use the correct scroll container and sticky offsets for existing headers/tab rows; give the toolbar an opaque background and appropriate stacking so results cannot obscure it.
- Wrap controls at mobile widths and preserve keyboard focus, entered queries and filter values across refreshes. Avoid horizontal page overflow or a toolbar that hides most of the available results on a short screen.

#### Folder and connection colors/icons

- Let users choose a color for folders and individual connections, reset it to default and preview the choice before saving.
- Offer a curated icon picker for folders and connections: systems/IT operations, servers/racks, network/router/switch, database, storage/backup, security, cloud, monitoring, OTT/video, live broadcast, encoder/transcoder and audio. A folder's selected icon replaces its folder symbol; a connection icon appears next to its protocol indicator.
- Show visual labels consistently in the connection tree, connection pickers and applicable window/tab/folder contexts. Keep the connection protocol identifiable and preserve text labels.
- Store preferences with additive schema changes and retain them across restarts, config export/import and authorized shares. Existing connections/folders retain their defaults; owner permissions control edits.
- Use supported icon identifiers and validated color values. Keep icons/text legible in light/dark themes and selected/hover states, add accessible names and do not rely on color alone for status or identification.

#### Acceptance and release

- Cover delayed response headers, long streaming, provider cancellation, timeout diagnostics and refusal to duplicate tool actions. Validate the original long-request scenario with Ollama.
- Verify Functionality search/filter combinations on a multi-server catalog with duplicate host/service names, errors, stale and unknown checks; verify that result refreshes preserve navigation state.
- Verify sticky controls in each inventoried view at 390, 1100, 1440 and 1920 px, browser zoom and short viewport heights. Inspect scroll behavior and screenshots.
- Cover icon/color validation, ownership/share visibility, old-database compatibility, defaults and export/import persistence; verify accessibility and both themes.
- During implementation, run the repository's Go race, vet/format, UI and real-browser checks. Add en/hr text and update affected feature/API/security documentation, changelog and release notes; link commits, PR and the v12.4.4 release to this scope. Deliver this remaining scope in v12.4.4.
