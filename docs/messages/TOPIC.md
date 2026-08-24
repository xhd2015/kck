---
name: kck/messages
description: >-
  Print recent coalesced Grok chat messages (msgfmt-style) with per-kind
  rune caps and offset-from-end paging.
---

# messages

```text
kck grok messages (<session-id> | --tab SEL | --tab-index N) [OPTIONS]
```

## Behavior

- Reads `updates.jsonl` from disk (no live pane required for an explicit id).
- Coalesces chunks into user / thinking / tool / response messages.
- Per-kind rune caps with U+2026 ellipsis: user 4096, tool 128, thinking 512,
  response 8192.
- Each line prefixed with local `[YYYY-MM-DD HH:MM:SS]`, or `[—]` if unknown.
- Paging: skip `--offset-from-end` newest, then take last `--limit` (default 32).

## Useful flags

- `--limit N` — page size (default 32; `0` = all remaining after offset)
- `--offset-from-end N` — skip N newest before `--limit` (default 0);
  example: `--offset-from-end 32` skips the last 32 and starts the next page
- `--json` — `{session_id,total,offset_from_end,limit,messages[]}` (no ANSI)

## Agent tips

1. First page: `--limit 32`. Next older page: `--offset-from-end 32 --limit 32`.
2. Prefer `--json` when paging programmatically (use `total` to know when to stop).
3. Use `--tab` when you only know the hosting pane, not the Grok id.
