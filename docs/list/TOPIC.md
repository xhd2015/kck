---
name: kck/list
description: >-
  kck list mode: live iTerm scan or --home store; filters, JSON, iTerm, limit.
---

# list

```text
kck [OPTIONS]
kck --home PATH [OPTIONS]
```

## Filters

| Flag | Effect |
|------|--------|
| `--needs-confirm` | Only needs-attention rows (live, not sendable, not exited) |
| `--sendable` | Only idle writable rows |
| `--limit N` | At most N newest rows |
| `--no-iterm` | Skip iTerm window/tab resolution (`ITERM` shows `-`) |
| `--fast` | Live only: skip `lsof`; `AGENT_SID` stays `-` |
| `--enrich` | Live only: also resolve grok title/model after sid hit |
| `--json` | Machine JSON list (buffered; no ANSI) |

There is **no** list `--send` / `--session`. To type into a pane use
`kck grok send`.

## Live vs store

- **Live** — streams rows as windows are scanned; includes agent-like cmds
  (`grok` / `codex` / `mark` / `agent-run`); omits plain bash.
- **Store** — newest-first by `updated_at`; `AGENT_RUN=yes`; sid from meta.

## `kck grok list` (Grok ids in iTerm)

```text
kck grok list [--json] [--limit N]
```

Grok-only: session ids with a **live process hard-hit** and a matching iTerm
tab. Columns: `SESSION_ID`, `ITERM`, `TITLE` (from disk summary), `WORKSPACE`.
One row per sid; multi-tab shows `w=… t=…(+N)`. Omits live PIDs with no tab.
Prefer this when you need Grok ids for `send` / `snapshot` / `open`.
Same core as `agent-pro grok session list-live`.

## Agent tips

1. Prefer `--json` when parsing programmatically.
2. Use `--needs-confirm` to find panes waiting on the user.
3. Use `--sendable` before a send workflow to confirm idle hosts.
4. Use `kck grok list` (not bare `kck`) when you only want Grok session ids.
