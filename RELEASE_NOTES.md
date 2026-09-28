## Web Remote Manager PRO v10.0.0-mimo

[![Support me on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/vedranius)

Browser-based remote server management: SSH terminal, SFTP/FTP/FTPS file manager, server-to-server transfer, sharing, real-time collaboration and **voice calls**. One self-contained binary with an embedded web UI.

v10 is a big release. It brings back **voice calls in collaboration rooms**, adds **user management, members and roles** for sharing, and **hardens security** throughout so WRM can be used in companies. Please read **Upgrading from v9** below before you replace the binary.

---

### 🎙 Voice calls in collaboration rooms (back, rebuilt)

Everyone in a share room can talk. It works like Discord or Jitsi:

- **🎙 Join voice** in the collaboration bar; **mute** (Ctrl+Shift+M), **deafen** (Ctrl+Shift+D), **leave**.
- **Push-to-talk** with a key of your choice (e.g. F8), or voice activity.
- **Microphone and speaker selection**, live input level meter, test sound, **noise suppression**, **echo cancellation**, **automatic gain control**.
- **Speaking indicators** (green ring), **per-person volume**, join/leave sounds, connection quality per person (P2P or relay, round-trip time).
- **Listen-only**: without a microphone, or without permission to use it, you can still join and listen.
- **Stable on real networks**: WebRTC with "perfect negotiation", automatic ICE restarts when the network changes, and automatic rejoin after reconnects.
- **Built-in TURN relay** (UDP and TCP, port 3478): calls work behind NAT and strict firewalls with no extra software. It accepts only short-lived credentials issued to room participants. By default it refuses to relay to private or loopback networks.
- **End-to-end encrypted** audio (DTLS-SRTP). The relay cannot decrypt it.
- Moderators can **mute someone for everyone**. Admins can turn voice on or off, limit participants per call (default 12), and add external STUN/TURN servers.

> Browsers allow the microphone only on **HTTPS** pages (or `localhost`). The quickest way to get HTTPS is `HTTPS_SELF_SIGNED=1`. For calls across the internet, open **UDP+TCP 3478** and the **UDP relay range** (default `49152-65535`) in the firewall. Behind NAT, set `turn_public_ip`.

### 👥 Users, members & roles

- **Admin panel → Users**: create users (a temporary password is generated if you leave it empty), set the display name, make or remove admin, **reset password**, **reset 2FA**, **sign out everywhere**, unlock, **disable/enable**, delete. The last active administrator cannot be removed.
- **Shares have three access modes**: 👥 *Members only* (default), 🏢 *Everyone signed in*, 🌐 *Anyone with the link* (guests; admins can disable this mode).
- **Members** with an individual role each, plus a role for everyone else. There are **five roles**:

  | | Observer | Viewer | Operator | Moderator | Owner |
  |---|:-:|:-:|:-:|:-:|:-:|
  | Chat & voice | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Watch shared terminals | ✔ | ✔ | ✔ | ✔ | ✔ |
  | Browse & download files | — | ✔ | ✔ | ✔ | ✔ |
  | Open terminals, change files, transfer | — | — | ✔ | ✔ | ✔ |
  | Share own terminal, receive keyboard control | — | — | ✔ | ✔ | ✔ |
  | Change roles, mute, remove, ban people | — | — | — | ✔ | ✔ |
  | Share settings & members | — | — | — | — | ✔ |

  The **server enforces roles** on every file operation, terminal and room action.
- **People** view per share: everyone who ever opened it, with last seen and IP. From there you can change a person's role or **ban** them.
- **Share expiry** (1 hour to 30 days, or a custom date), **pause/resume**, **new link** (the old link stops working immediately).
- **Revocation takes effect at once.** Deleting, pausing or rotating a share, removing a member, lowering a role, banning, or disabling a user closes their open terminals and room connections.

### 💬 Collaboration

- A **new collaboration bar** shows the share, your role, avatars of the people online, and the voice controls.
- **👥 People panel**: who is in the call and who is online, with roles, mute and deafen state, shared terminals, raised hands, connection quality and volume. Moderators have a **⋯** menu for *change role*, *mute*, *stop sharing*, *lower hand*, *remove* and *ban*.
- **Chat history** is kept, so people who join later see it. The chat has timestamps, clickable links, an unread badge and notifications. **File exchange** up to the policy limit: images get a preview, and other files are always downloaded, never opened in the page.
- **✋ Raise hand.**
- **Terminal sharing**: you can share several terminals at once. Viewers get a **snapshot of the current screen** (with colors) and then the live output. Terminal data is sent only to people who watch it.
- **Remote keyboard control**: *Request control* → *Allow / Deny*. Keystrokes go only to the granted terminal, and only while it is shared.
- **Identities come from the server.** Signed-in users appear under their account name. Guests pick a name and are marked *guest*.
- The room **reconnects automatically** and restores watching, sharing and the voice call.
- **Per-participant send queues**: one slow client can no longer stall a room.

### 🔐 Security hardening

- **Two-factor authentication (TOTP)** with QR code and single-use **recovery codes**. The policy can **require 2FA** for admins or for everyone.
- **Self-registration is closed by default.** Only the first account (the administrator) registers itself.
- **Accounts lock** after repeated wrong passwords (policy). Per-IP rate limiting. Unknown users get constant-time responses. After an admin reset, the user **must change the password**.
- **Sessions**: 256-bit tokens, **stored hashed**. Idle and maximum lifetime (policy). A list of **signed-in devices** with remote sign-out.
- **Stored secrets** (connection passwords, private keys, 2FA secrets) are encrypted with **AES-256-GCM** and a **random key per installation** (`remote_manager.db.key`, or your own `ENCRYPTION_KEY` / `ENCRYPTION_KEY_FILE`).
- **Secrets are never sent to the browser**, not even to their owner. The *Edit connection* dialog shows "saved and encrypted — leave empty to keep it". *Export* leaves secrets out. *Export with passwords & keys* asks for your password again.
- **SSH host key verification** (trust on first use, or strict). A changed key is refused and shown with its fingerprint. **FTPS certificates** are validated, or pinned on first use.
- **Server-side key files** (*Key file* / *Auto (~/.ssh)*) are limited to administrators by default.
- **CSRF protection** on every state-changing API call. **Content-Security-Policy**, `X-Frame-Options`, `nosniff`, `Referrer-Policy: no-referrer`, `Permissions-Policy`, and HSTS on HTTPS. Request-size limits. No directory listings.
- **All assets are served locally** (xterm.js and fonts are embedded). There is no CDN dependency at runtime.
- **Share passwords** are held in signed cookies, which become invalid when the password changes. Guests get their own identity cookie. Share connection lists stay hidden until the share password is entered.
- **HTTPS made easy**: `HTTPS_CERT_FILE` + `HTTPS_KEY_FILE`, or `HTTPS_SELF_SIGNED=1`. Behind a reverse proxy, `WRM_TRUST_PROXY=1` gives correct client IPs and secure cookies.
- **Audit log** of sign-ins (failed ones too), account and policy changes, connections, shares, room joins and moderation, terminals, file changes, host keys, exports and imports. It is searchable in the admin panel, exportable as **CSV**, and written as `AUDIT …` lines to the server log (journald/syslog → SIEM).
- **Recovery from the command line**: `wrm -reset-password USER [-reset-2fa]`.

### 🛡 Admin panel

*Settings → 🛡 Admin panel* has these tabs:
- **Overview**: version, uptime, users, HTTPS, relay status, key source, **security warnings**, **open terminals** (with *End*), live rooms.
- **Users**
- **Shares**: every share on the server.
- **Security policies**
- **Voice & network**
- **Host keys**
- **Audit log**

Any policy can also be **forced by an environment variable** (`WRM_<KEY>`, e.g. `WRM_REQUIRE_2FA=all`). A forced policy is shown as locked in the panel.

### ✨ Other improvements

- **Settings** are reorganised into tabs: General, Terminal, Security, Voice & audio, Data.
- **Connections**: SSH keyboard-interactive authentication (servers that ask for the password that way). *Duplicate* copies the stored secrets on the server, so they never pass through the browser. Deleting a connection closes its open terminals.
- **Uploads** respect the server-side size limit (`max_upload_mb`), and partial files are removed.
- Translations (English / Hrvatski) cover the whole new UI.
- A **unit test suite** (TOTP, recovery codes, roles, share access, CSRF, encryption and migration, signed values, settings, input validation) runs in CI with the race detector.

---

### ⬆️ Upgrading from v9

Replace the binary and start it with the same database, and the same `ENCRYPTION_KEY` if you set one. Everything is migrated automatically.

- **Stored secrets are re-encrypted.** If `ENCRYPTION_KEY` was not set, v9 used a public default key. v10 generates a random key file, `remote_manager.db.key`, next to the database and re-encrypts all secrets with it. **Back that file up separately from the database.** If you did set `ENCRYPTION_KEY`, nothing changes.
- **Self-registration is now closed.** Existing accounts keep working. To reopen it: *Admin panel → Security policies*.
- **SSH host keys are now verified.** The first connection to each server remembers its key.
- **Key file / Auto (~/.ssh)** authentication is now admin-only, unless you turn on the policy *Allow server key files for all users*.
- **Existing shares keep their behaviour**: anyone with the link can open them, with the role *Operator*. Edit a share to switch it to members only or to change roles.
- **Behind a reverse proxy**, state-changing API calls now require the `Origin` to match the `Host`. Pass the `Host` header through (`proxy_set_header Host $host;`) or set `WRM_ALLOWED_ORIGINS`. Also set `WRM_TRUST_PROXY=1`.
- For voice calls, open **UDP+TCP 3478** and the relay port range in the firewall (see above).

---

### 📦 Downloads

| Platform | Architecture | Binary |
|---|---|---|
| Linux | x86-64 | `wrm-pro-v10.0.0-mimo-linux-amd64` |
| Linux | arm64 | `wrm-pro-v10.0.0-mimo-linux-arm64` |
| Linux | ARMv7 (Raspberry Pi) | `wrm-pro-v10.0.0-mimo-linux-armv7` |
| Linux | ARMv6 | `wrm-pro-v10.0.0-mimo-linux-armv6` |
| Linux | 32-bit | `wrm-pro-v10.0.0-mimo-linux-386` |
| Windows | x86-64 | `wrm-pro-v10.0.0-mimo-windows-amd64.exe` |
| Windows | arm64 | `wrm-pro-v10.0.0-mimo-windows-arm64.exe` |
| macOS | Intel | `wrm-pro-v10.0.0-mimo-darwin-amd64` |
| macOS | Apple Silicon | `wrm-pro-v10.0.0-mimo-darwin-arm64` |
| macOS | Universal | `wrm-pro-v10.0.0-mimo-darwin-universal` |
| Android | arm64 (Termux) | `wrm-pro-v10.0.0-mimo-android-arm64` |
| FreeBSD | x86-64 | `wrm-pro-v10.0.0-mimo-freebsd-amd64` |
| FreeBSD | arm64 | `wrm-pro-v10.0.0-mimo-freebsd-arm64` |
| OpenBSD | x86-64 | `wrm-pro-v10.0.0-mimo-openbsd-amd64` |

Verify integrity with `SHA256SUMS.txt`. The Android build has no built-in TURN relay; configure an external TURN server there if you need one.

### 🚀 Quick start

**Linux / macOS**
```bash
chmod +x wrm-pro-v10.0.0-mimo-linux-amd64
HTTPS_SELF_SIGNED=1 ./wrm-pro-v10.0.0-mimo-linux-amd64
# open https://<server>:8080 — create the administrator account (the first account)
```
On macOS, if Gatekeeper blocks the file: `xattr -d com.apple.quarantine wrm-pro-*-darwin-*`.

**Windows**: double-click `wrm-pro-v10.0.0-mimo-windows-amd64.exe`, or in PowerShell:
```powershell
$env:HTTPS_SELF_SIGNED=1; .\wrm-pro-v10.0.0-mimo-windows-amd64.exe
```

**Android (Termux)**
```bash
pkg install wget
wget https://github.com/vedranius/web-browser-RDM-public/releases/download/v10.0.0-mimo/wrm-pro-v10.0.0-mimo-android-arm64
chmod +x wrm-pro-v10.0.0-mimo-android-arm64 && ./wrm-pro-v10.0.0-mimo-android-arm64
```

On the first start WRM creates `remote_manager.db` (the database) and `remote_manager.db.key` (the encryption key) next to each other, both with mode `0600`. **Back up the key file.**

### ⚙️ Most important settings

```
PORT=8080 / LISTEN_ADDR=:8080     listening address
DB_PATH=./remote_manager.db       database
ENCRYPTION_KEY / ENCRYPTION_KEY_FILE   your own key for stored secrets (else <DB_PATH>.key)
HTTPS_CERT_FILE + HTTPS_KEY_FILE  HTTPS with your certificate
HTTPS_SELF_SIGNED=1               HTTPS with a generated self-signed certificate
WRM_TRUST_PROXY=1                 behind a reverse proxy (real client IPs, secure cookies)
WRM_ALLOWED_ORIGINS=host          extra allowed Origin hosts
WRM_<POLICY>=value                force a policy, e.g. WRM_REQUIRE_2FA=all, WRM_TURN_PUBLIC_IP=203.0.113.10
```

For the full documentation (reverse proxy, systemd, firewall, API), see the [README](https://github.com/vedranius/web-browser-RDM-public/blob/v10.0.0-mimo/README.md). To report a vulnerability, see [SECURITY.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.0.0-mimo/SECURITY.md).

---

### 📜 License

Web Remote Manager PRO is **source-available** under the [PolyForm Noncommercial License 1.0.0](https://github.com/vedranius/web-browser-RDM-public/blob/v10.0.0-mimo/LICENSE). It is free for personal, educational, non-profit and other noncommercial use, including forks and modifications. **Commercial use requires a separate license**; see [COMMERCIAL-LICENSE.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.0.0-mimo/COMMERCIAL-LICENSE.md). Contributions are welcome; see [CONTRIBUTING.md](https://github.com/vedranius/web-browser-RDM-public/blob/v10.0.0-mimo/CONTRIBUTING.md).

☕ **Like WRM?** Support its development on **[Ko-fi](https://ko-fi.com/vedranius)**. Thank you!
