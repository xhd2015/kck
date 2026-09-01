---
name: kck/focus
description: >-
  Focus the iTerm tab hosting a live Grok/Codex session. Lighter than open:
  never resumes or creates a window when no live host.
---

# focus

```text
kck grok focus <session-id> [--index N]
kck codex focus <session-id> [--index N]
```

## Behavior

- Thin wrapper over agent-pro `RunFocus`.
- Requires a **live** hosting iTerm tab (PID → TTY → pane).
- On success prints `focused: window W, tab T`.
- Unknown session or no live host → `Error: not found` (does **not** resume).
- Multiple hosts → error listing candidates; use `--index N`.

## vs open

| | `focus` | `open` |
|--|---------|--------|
| Live host | focus tab | focus tab |
| No live host | fail | resume in new window |

## Agent tips

1. Prefer `focus` when you only want to switch panes and know the session is live.
2. Use `open` when you need resume-if-exited.
