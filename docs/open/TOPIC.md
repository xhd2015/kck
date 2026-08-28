---
name: kck/open
description: >-
  Focus the iTerm tab hosting a Grok or Codex session, or resume in a new window.
---

# open

```text
kck grok open (<session-id> | --tab SEL | --tab-index N) [OPTIONS]
kck codex open (<session-id> | --tab SEL | --tab-index N) [OPTIONS]
```

## Behavior

- If a hosting iTerm tab already exists → **focus** it (exactly one host, or
  pick with `--index` when multiple).
- Else if the id is bound in agent-run → prefer agent-run (live → focus;
  exited → `agent-run run --auto-send-or-resume --open` in a new window).
- Else → open a new iTerm window and run `grok --resume <session-id>` or
  `codex resume <session-id>`.
- `--no-agent-run` forces bare resume (skip agent-run prefer).
- `--tab` / `--tab-index` resolve **only** focuses (never resumes).

## Session source (exactly one)

- positional `<session-id>`
- `--tab SEL` — 1-based index, or `next|left|right` (`right` ≡ `next`)
- `--tab-index N` — 0-based index in this iTerm window

## Agent tips

1. Prefer open when the user wants to *see* the session UI.
2. For scripting text into a pane without focusing, use `snapshot` / `send`.
3. Unknown session → hard `Error:`.
