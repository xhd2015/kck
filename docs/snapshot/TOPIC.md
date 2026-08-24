---
name: kck/snapshot
description: >-
  Capture currently visible pane text for a live Grok session host
  (agent-run TTY when managed, else iTerm Contents).
---

# snapshot

```text
kck grok snapshot (<session-id> | --tab SEL | --tab-index N) [OPTIONS]
```

## Behavior

- Prints visible pane text to stdout (or `-o FILE`).
- Does **not** focus the pane.
- When the Grok id is bound to a **live agent-run** grok-tty session, prefers
  that sanitized TTY snapshot (single frame). Otherwise uses iTerm2 Contents.
- Bare grok (not under agent-run) always uses iTerm. No live host and no
  agent-run hit → hard error (no resume).

## Useful flags

- `--json` — `{"session_id","iterm_session_id","app","source","contents"}`
  (`source` is `agent-run` or `iterm`; `agent_run_session_id` when applicable)
- `--dry-run` — resolve only; do not capture
- `--iterm` — force iTerm Contents (skip agent-run prefer)
- `--index N` — disambiguate multi-host (positional id only)

## Agent tips

1. Snapshot before deciding what to `send`.
2. Prefer `--json` when the contents must be parsed.
3. Use `--iterm` to compare against the agent-run frame when debugging.
