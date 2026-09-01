# kck

List live iTerm agent panes and operate Grok/Codex sessions (open, snapshot, send, messages, info, status, resolve, pickup). Thin CLI over [agent-pro](https://github.com/xhd2015/agent-pro) session helpers; agent-run-managed sessions prefer agent-run for send/resume.

## Requirements

- **macOS** with **iTerm2** for live pane scan and iTerm-hosted Grok ops
- Optional **agent-run home** for `kck --home PATH` (store list) and managed send/resume

## Install

Build from a checkout that already has dependency trees under `external/` (see `go.mod` `replace` directives; `/external` is local-only):

```sh
go install .
# or: go build -o ~/go/bin/kck .
```

Then ensure `~/go/bin` (or your chosen output dir) is on `PATH`.

## Quick start

```sh
# live iTerm agent panes (streams as windows are scanned)
kck

# agent-run store sessions
kck --home ~/.agent-run

# Grok ids hosted in iTerm
kck grok list

# inspect / operate a session
kck grok info <session-id>
kck grok status <session-id>
kck grok send "hello" --session-id <session-id>
kck grok pickup "summarize decisions" --session-id <session-id>

# Codex: list / open / send / messages / info / status / resolve / snapshot / pickup
kck codex list
kck codex open --tab 1
kck codex send "hello" --session-id <id>
kck codex messages <id>
kck codex info <id>
kck codex status <id>
kck codex resolve --tab 1
kck codex snapshot --tab 1
kck codex pickup "extract TODOs" --tab 1

# embedded skill (when/why + topics)
kck skill --list
kck skill --show
kck skill --show send
kck skill --show pickup
```

Domain flags: `kck --help`, `kck grok <command> --help`, and `kck codex <command> --help`.

## Commands

| Command | Role |
|---------|------|
| `kck` / `kck --home PATH` | Live iTerm panes, or agent-run store under `PATH/sessions` |
| `kck grok list` | Grok ids in iTerm tabs |
| `kck grok open` | Focus hosting tab or resume |
| `kck grok snapshot` | Capture visible pane text |
| `kck grok send` | Type into hosting pane (agent-run when managed) |
| `kck grok messages` | Recent coalesced chat |
| `kck grok info` / `status` | Detail + Active; dual-signal liveness |
| `kck grok resolve` | Resolve id (ancestor walk or `--tab`) |
| `kck grok pickup` | New empty session staged from a base session (kck-pickup-a-session) |
| `kck codex list` | Codex ids in iTerm tabs |
| `kck codex open` | Focus hosting tab or resume |
| `kck codex snapshot` | Capture visible Codex pane text |
| `kck codex send` | Type into hosting pane (agent-run when managed) |
| `kck codex messages` | Recent coalesced chat |
| `kck codex info` / `status` | Detail + Active; PID liveness (File always no) |
| `kck codex resolve` | Resolve Codex id (ancestor walk or `--tab`) |
| `kck codex pickup` | New empty session staged from a base session (kck-pickup-a-session) |
| `kck skill` | Show/install embedded skill docs |

## Docs

Deep when/why lives in the embedded skill (not a second flag encyclopedia):

```sh
kck skill --show overview   # modes, columns, soft-fail culture
kck skill --show list
kck skill --show open
kck skill --show send
kck skill --show snapshot
kck skill --show messages
kck skill --show info
kck skill --show status
kck skill --show resolve
kck skill --show pickup
```

Source topics: [`docs/SKILL.md`](docs/SKILL.md) and [`docs/*/TOPIC.md`](docs/).

Probe / list failures that still allow partial output print `warning:` on stderr and exit 0. Hard usage / not-found errors print `Error:` on stderr and exit non-zero.
