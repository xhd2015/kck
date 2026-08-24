---
name: kck/info
description: >-
  Show Grok session detail plus an Active block (live host / process signals).
---

# info

```text
kck grok info <session-id> [OPTIONS]
```

## Behavior

- Prints human-readable session metadata (paths, timestamps, workspace).
- Includes an **Active** block summarizing whether a live host/process is
  attached.
- Unknown session → hard `Error:`.

## Agent tips

1. Use `info` when the user asks “what is this session?” or for cwd / paths.
2. Prefer `status` for a compact liveness check before send/open.
