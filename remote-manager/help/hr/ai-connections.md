---
id: ai-connections
title: AI aplikacije preko MCP-a
section: AI
order: 110
summary: Neka Claude, ChatGPT, Claude Code, Codex, Cursor i druge MCP aplikacije rade na vašim serverima kroz WRM, uz vaše odobrenje.
keywords: mcp, model context protocol, claude desktop, claude.ai, chatgpt, claude code, cursor, gemini cli, opencode, codex, konektor, token, oauth, ai veza, plugin, mcpb, most, odobri, connector, ai connection, bridge, approve
ui: #st-mcp, #mcp-pane, #btn-mcp
open: mcp_setup
---
WRM je MCP server. AI aplikacija na vašem računalu može otvoriti WRM sesiju na server koji ste dopustili i raditi na njemu po istim pravilima kao ugrađeni asistent. Model radi u aplikaciji; WRM za to nikad ne treba ključ pružatelja.

## Spajanje aplikacije
1. Administrator uključuje **AI veze** (*✦ AI → 🧭 Postavi AI* ili *Admin → AI veze*).
2. **Postavke → AI veze → Spoji AI aplikaciju** i odaberite aplikaciju:
   - **Claude Desktop:** preuzmite WRM proširenje (`.mcpb`) i otvorite ga.
   - **claude.ai / ChatGPT:** dodajte WRM-ovu MCP adresu kao konektor (za njih WRM mora biti dostupan s interneta).
   - **Claude Code, Codex, Gemini CLI, Cursor, OpenCode:** kopirajte prikazanu naredbu ili konfiguracijsku datoteku.
   - **Sve AI aplikacije odjednom:** jedna skripta za Windows ili Linux/macOS postavlja svaku pronađenu aplikaciju.
3. Odaberite **servere**, **način** (samo čitanje, pitaj, automatski) i koliko dugo pristup traje. Aplikacija se prijavljuje OAuthom ili osobnim tokenom.

## Dok radi
**🤖 AI spojen** u gornjoj traci pokazuje sesije uživo; pregled uživo pokazuje svaki korak. Promjene čekaju vaše odobrenje u WRM-u (ili u aplikaciji kad ona može pitati). Vezu opozovete jednim klikom; administratori ih sve vide u *Admin → AI veze*.

## Pomoć u vašoj AI aplikaciji
Svaka AI aplikacija spojena na WRM dobiva i alat **wrm_help**: pitajte je *kako dodati jump host u WRM-u?* i ona čita ovu pomoć.
