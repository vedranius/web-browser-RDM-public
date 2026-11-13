---
id: help-and-ai-help
title: Using Help and the help AI
section: AI
order: 115
summary: Search the help, pick any part of the screen to ask about it, and let the help AI explain or do it for you.
keywords: help, f1, question mark, search, pick, select, frame, explain, how to, ask ai, do it for me, documentation, manual, guide
ui: #help-panel, #btn-help
open: help:help-and-ai-help
---
**?** in the top bar or **F1** opens Help on the side; WRM keeps working next to it, and you can keep a dialog open while reading.

## Find an answer
- **Search** in English or Croatian, with or without č, ć, š, ž, đ.
- **Topics** lists every article by section; each one also has a link to the same page on GitHub.
- Help opens on the article about what is in front of you (for example the dialog you have open).

## Pick a part of the screen
1. Press **⌖ Pick** in Help.
2. Point at a button, a dialog, a window or a list: it is outlined. Click to pick it, or **drag a frame** around an area. **Esc** cancels.
3. Help shows the articles about it and puts it into the AI chat as context, so you can ask *what does this do?* or *why is this red?* Only the visible text is sent, never the values of form fields.

## Ask the help AI
Type a question in your own words, for example *how do I reach a server behind a bastion?* The answer comes from these articles, with menu paths and steps.

If you ask it to **do it for you**, it offers buttons such as *Open new connection (filled in)*, *Open Settings → AI* or *Open app-01 with the AI assistant*. Each button says exactly what it does, and **nothing happens until you press it**. A setting change asks once more before it is saved, and changes on a server still go through the AI assistant's approvals.

The help AI needs an AI service (see *Set up AI*); search and articles work without one. Administrators choose who may use it with the policy `ai_help` (default: the same people as the AI assistant). Questions are audited (`help.ask`), without their text.

Error messages have a **?** button too: it opens Help with the error filled in.
