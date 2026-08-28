---
name: kck/status
description: >-
  Dual-signal Grok session liveness plus session path (compact status line).
---

# status

```text
kck grok status <session-id> [OPTIONS]
kck codex status <session-id> [OPTIONS]
```

## Behavior

- Grok: dual signals — `active_sessions.json` (File) + live PIDs.
- Codex: PID-only (File always no; no Codex active-session registry). Includes
  the rollout path for quick navigation.
- Unknown session → hard `Error:`.

## Agent tips

1. Cheap preflight before `send` / `open` when you only need alive/dead.
2. Use `info` when you need richer metadata and the Active block.
