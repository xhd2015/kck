---
name: kck
description: >-
  kck lists live iTerm agent panes (or agent-run store sessions with --home)
  and operates Grok sessions: open/focus, snapshot, send, messages, info, status.
  Use when the user runs kck, asks about live agent tabs, or /kck.
  Load topics: kck skill --show <topic>
---

# kck — live agent panes and Grok session ops

This skill is an **index**. Load detailed guidance with
`kck skill --show <topic>` (or `kck skill <topic> --show`).

**kck** has two list modes and a `grok` command group:

- **Default (no `--home`)** — scan live iTerm windows for agent-like panes;
  stream rows as windows are scanned.
- **`--home PATH`** — list agent-run store sessions under `PATH/sessions`.
- **`kck grok …`** — list iTerm-hosted Grok ids, open/focus, snapshot, send,
  messages, info, status (thin wrappers over agent-pro session helpers).

Domain flags live in `kck --help` and `kck grok <cmd> --help`. This skill is
for when/why and agent workflows, not a full flag encyclopedia.

## Default agent workflow

1. Orient with bare `kck` (live multi-agent) or `kck grok list` (Grok ids in
   iTerm only). Use `kck --home <agent-run-home>` for the agent-run store.
2. Use `--needs-confirm` / `--sendable` when filtering attention vs idle.
3. For a known Grok session: `kck grok info|status <id>`, then
   `open` / `snapshot` / `messages` / `send` as needed.
4. Prefer `kck grok send "…" --session-id <id>` over inventing list-mode send
   flags (root list has no `--send`).

## Topics

- `overview` — modes, columns (`AGENT_RUN` / `AGENT_SID`), soft-fail culture
- `list` — live + store list; `kck grok list` for iTerm-hosted Grok ids
- `open` — focus hosting tab or resume; `--tab` / `--tab-index`
- `snapshot` — capture visible pane text (agent-run TTY when managed, else iTerm)
- `send` — type into hosting pane; agent-run-managed `--session-id` sends via agent-run (no iTerm); `--no-agent-run` forces iTerm; `--open` resume-then-send for unmanaged
- `messages` — recent coalesced chat (msgfmt-style; `--limit` / `--offset-from-end`)
- `info` — session detail + Active block
- `status` — dual-signal liveness + session path

## Retrieve topics

```bash
# list skill name + every nested topic path
kck skill --list

# root skill index (this document)
kck skill --show

# topic (both flag orders)
kck skill --show send
kck skill send --show
kck skill --show overview
kck skill --show list

# YAML frontmatter only
kck skill --show --header
kck skill --show send --header
```

## Related CLI

Domain usage lives in `kck --help` and `kck grok <command> --help`. Use this
skill for when/why and nested-topic recipes.
