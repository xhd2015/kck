# kck

List live iTerm agent panes and operate Grok sessions (open, snapshot, send, messages, info, status). Thin CLI over [agent-pro](https://github.com/xhd2015/agent-pro) session helpers; agent-run-managed sessions prefer agent-run for send/resume.

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

# embedded skill (when/why + topics)
kck skill --list
kck skill --show
kck skill --show send
```

Domain flags: `kck --help` and `kck grok <command> --help`.

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
```

Source topics: [`docs/SKILL.md`](docs/SKILL.md) and [`docs/*/TOPIC.md`](docs/).

Probe / list failures that still allow partial output print `warning:` on stderr and exit 0. Hard usage / not-found errors print `Error:` on stderr and exit non-zero.
