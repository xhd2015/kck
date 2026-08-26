---
name: kck/messages
description: >-
  Print recent coalesced Grok chat messages (msgfmt-style) with per-kind
  rune caps, optional --grep AND filter, and offset-from-end paging.
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
- Optional `--grep`: keep messages whose body contains every pattern (AND;
  case-insensitive literal), then page.
- Paging: skip `--offset-from-end` newest matches, then take last `--limit`
  (default 32).
- Text header: offset 0 full → `showing all K of N`; offset 0 partial →
  `showing last K of N`; offset > 0 → `showing lo-hi(K) of N` (1-based
  oldest→newest indices).
- Streams message lines as they are formatted; `--json` is one document.

## Useful flags

- `--limit N` — page size (default 32; `0` = all remaining after offset)
- `--offset-from-end N` — skip N newest before `--limit` (default 0);
  example: `--offset-from-end 32` skips the last 32 and starts the next page
- `--grep P` — repeatable; AND across patterns (applied before paging)
- `--color` / `--no-color` — force ANSI on/off (auto: TTY + `NO_COLOR`)
- `--json` — `{session_id,total,offset_from_end,limit,messages[]}` (no ANSI;
  `total` is post-grep count)

## Agent tips

1. First page: `--limit 32`. Next older page: `--offset-from-end 32 --limit 32`.
2. Prefer `--json` when paging programmatically (use `total` to know when to stop).
3. Use `--tab` when you only know the hosting pane, not the Grok id.
4. Narrow with `--grep A --grep B` when both tokens must appear in one message.
