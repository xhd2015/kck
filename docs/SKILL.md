---
name: kck
description: >-
  kck lists live iTerm agent panes (or agent-run store sessions with --home)
  and operates Grok/Codex sessions: open, focus, snapshot, send, messages,
  prompts, info, status, resolve, pickup, new. Use when the user runs kck, asks about
  live agent tabs, or /kck. Load topics: kck skill --show <topic>
---

# kck — live agent panes and Grok / Codex session ops

This skill is an **index**. Load detailed guidance with
`kck skill --show <topic>` (or `kck skill <topic> --show`).

**kck** has two list modes plus runner command groups:

- **Default (no `--home`)** — scan live iTerm windows for agent-like panes;
  stream rows as windows are scanned.
- **`--home PATH`** — list agent-run store sessions under `PATH/sessions`.
- **`kck grok …`** — list iTerm-hosted Grok ids, open, focus, snapshot, send,
  messages, prompts, info, status, resolve, pickup, new (thin wrappers over agent-pro
  session helpers).
- **`kck codex …`** — `list`, `open`, `focus`, `snapshot`, `send`, `messages`,
  `prompts`, `info`, `status`, `resolve`, `pickup`, `new` (same shapes as the Grok
  counterparts where applicable; multi-hit on one tab still refuses; Codex status
  File is always no).

Domain flags live in `kck --help`, `kck grok <cmd> --help`, and
`kck codex <cmd> --help`. This skill is for when/why and agent workflows,
not a full flag encyclopedia.

## Default agent workflow

1. Orient with bare `kck` (live multi-agent) or `kck grok list` (Grok ids in
   iTerm only). Use `kck --home <agent-run-home>` for the agent-run store.
2. Use `--needs-confirm` / `--sendable` when filtering attention vs idle.
3. For a known Grok session: `kck grok info|status <id>`, then
   `open` / `snapshot` / `messages` / `send` as needed. When the id is unknown
   in-context, `kck grok resolve` (ancestor or `--tab`) yields a bare id.
4. Prefer `kck grok send "…" --session-id <id>` over inventing list-mode send
   flags (root list has no `--send`).

## Topics

- `overview` — modes, columns (`AGENT_RUN` / `AGENT_SID`), soft-fail culture
- `list` — live + store list; `kck grok list` / `kck codex list` for iTerm-hosted ids
- `open` — focus hosting tab or resume; `--tab` / `--tab-index`; Codex: `kck codex open`
- `focus` — focus live hosting tab only (no resume); Codex: `kck codex focus`
- `snapshot` — capture visible pane text (agent-run TTY when managed, else iTerm);
  Codex: `kck codex snapshot`
- `send` — type into hosting pane; agent-run-managed `--session-id` sends via agent-run (no iTerm); `--no-agent-run` forces iTerm; `--open` resume-then-send for unmanaged; `--cron` on Grok and Codex
- `messages` — recent coalesced chat (msgfmt-style; `--limit` / `--grep` / `--offset-from-end`);
  Codex: `kck codex messages`
- `prompts` — user prompts only (`--first` / `--grep` / `--this-window` / `--this-space` / `--tab`);
  Codex: `kck codex prompts`
- `info` — session detail + Active block; Codex: `kck codex info` (File always no)
- `status` — liveness + session path; Codex: PID-only (`kck codex status`)
- `resolve` — resolve session id (ancestor walk or `--tab` / `--tab-index`);
  Codex: `kck codex resolve` with the same flags
- `pickup` — open a new empty session staged with a kck-pickup-a-session draft
  from a base session (not a native fork); Codex: `kck codex pickup`
  (agent skill: `kck skill --show kck-pickup-a-session`)
- `new` — open a new empty session via agent-run (brainstorm-prefixed prompt;
  default new terminal; `--here` silent); Codex: `kck codex new`

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

Domain usage lives in `kck --help`, `kck grok <command> --help`, and
`kck codex <command> --help`. Use this skill for when/why and nested-topic recipes.
