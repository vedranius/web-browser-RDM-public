---
id: ai-setup
title: Postavi AI (kratki put)
section: AI
order: 100
summary: Jedna stranica koja uključuje AI, dodaje AI servis (Claude, ChatGPT, Gemini, lokalni modeli…) i spaja AI aplikacije.
keywords: ai, postavi, postavljanje, uključi, omogući, pružatelj, servis, api ključ, claude, anthropic, openai, chatgpt, gemini, ollama, lokalni model, openrouter, mistral, groq, deepseek, opencode, besplatno, mcp, spoji, čarobnjak, modeli, setup, enable, provider, api key, local model, free, connect, wizard
ui: #ai-setup, #st-ai
open: ai_setup
---
**✦ AI → 🧭 Postavi AI** (i u *Postavke → AI asistent* te u Pomoći) je jedna stranica sa svim koracima i svjetlom za svaki: zeleno kad je gotov, s gumbom za ono što nedostaje.

## 1. Uključite AI (administrator)
Odaberite tko ga smije koristiti: **Administratori** ili **Svi**. AI servisu se ništa ne šalje dok se servis ne doda i netko ne pita.

## 2. Dodajte AI servis
**Dodaj AI servis** otvara kratki čarobnjak: odaberite servis, zalijepite ključ, odaberite model, **Testiraj**. Kad test uspije, WRM se spaja na servis i **automatski dodaje sve modele koje nudi** (za Ollamu sve instalirane modele). Gumb **⟳ Modeli** na popisu servisa to ponavlja kad servis doda nove modele.
- **Claude (Anthropic), ChatGPT (OpenAI), Gemini, Mistral, xAI, DeepSeek, Groq, OpenRouter, Together, Fireworks, Cerebras, Cohere, Perplexity:** zalijepite API ključ iz konzole servisa.
- **OpenCode Zen:** besplatni modeli **bez ključa**, dobar način da isprobate.
- **Ollama (lokalni modeli, ništa ne napušta vašu mrežu):** čarobnjak pronalazi Ollamu, a administrator je može instalirati, pokrenuti i preuzeti model s po jednim klikom.
- Azure OpenAI, AWS Bedrock i Google Vertex AI su u punom obrascu (*Admin → AI asistent → Pružatelji*).
Administratori dodaju servise za sve; korisnici dodaju svoj ključ pod *Moji API ključevi* kad pravilo `ai_personal_keys` to dopušta. Koje modele korisnici smiju birati ograničava *Admin → AI → Dopušteni modeli*.

## 3. Spojite AI aplikacije (neobavezno)
Za rad s WRM-om iz **Claude Desktopa, claude.ai, ChatGPT-a, Claude Codea, Codexa, Gemini CLI-ja, Cursora, OpenCodea** ili druge MCP aplikacije:
1. Administrator uključuje **AI veze (MCP)** na istoj stranici.
2. **Spoji AI aplikaciju** otvara čarobnjak u *Postavke → AI veze*: odaberite aplikaciju, servere i način, pa kopirajte jedan redak ili preuzmite datoteku koju pokaže. **Sve AI aplikacije odjednom** preuzima jednu skriptu za postavljanje.

## Što dalje
- **✦ AI** pokraj terminala, u uređivaču datoteka, pregledniku logova, Lancima i Pomoći.
- AI **nikad ništa ne pokreće sam**: kaže što želi pokrenuti, a promjene čekaju vaše odobrenje.
