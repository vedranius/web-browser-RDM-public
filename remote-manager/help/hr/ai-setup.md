---
id: ai-setup
title: Set up AI (the short way)
section: AI
order: 100
summary: One page that turns AI on, adds an AI service (Claude, ChatGPT, Gemini, local models…) and connects AI apps.
keywords: ai, setup, set up, turn on, enable, provider, api key, claude, anthropic, openai, chatgpt, gemini, ollama, local model, openrouter, mistral, groq, deepseek, opencode, free, mcp, connect, wizard
ui: #ai-setup, #st-ai
open: ai_setup
---
**✦ AI → Set up AI** (also in Help) is one page with every step and a light for each: green when done, with a button for what is missing.

## 1. Turn the AI on (administrator)
Choose who may use it: **Administrators** or **Everybody**. Nothing is sent to an AI service until a service is added and somebody asks.

## 2. Add an AI service
**Add an AI service** opens a short wizard: pick the service, paste its key, choose the model, **Test**.
- **Claude (Anthropic), ChatGPT (OpenAI), Gemini, Mistral, xAI, DeepSeek, Groq, OpenRouter, Together, Fireworks, Cerebras, Cohere, Perplexity:** paste an API key from the service's console.
- **OpenCode Zen:** free models **without a key**, a good way to try it.
- **Ollama (local models, nothing leaves your network):** the wizard finds Ollama, and an administrator can install it, start it and download a model with one click each.
- Azure OpenAI, AWS Bedrock and Google Vertex AI are in the full form (*Admin → AI assistant → Providers*).
Administrators add services for everybody; users add their own key under *My API keys* when the policy `ai_personal_keys` allows it.

## 3. Connect your AI apps (optional)
To work with WRM from **Claude Desktop, claude.ai, ChatGPT, Claude Code, Codex, Gemini CLI, Cursor, OpenCode** or another MCP app:
1. An administrator turns on **AI connections (MCP)** on the same page.
2. **Connect an AI app** opens the wizard in *Settings → AI connections*: pick the app, choose servers and mode, and copy the one line or download the file it shows. **All AI apps at once** downloads one setup script.

## What happens next
- **✦ AI** next to a terminal, in the file editor, the log viewer, Chains and Help.
- The AI **never runs anything by itself**: it says what it wants to run, and changes wait for your approval.
