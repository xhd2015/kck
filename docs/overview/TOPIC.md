---
name: kck/overview
description: >-
  kck modes (live vs --home store), table columns, and soft-fail culture.
---

# overview

## Modes

| Invocation | Source |
|------------|--------|
| `kck` | Live iTerm scan (agent-like panes only) |
| `kck --home PATH` | agent-run FileStore under `PATH/sessions/<id>/meta.json` |

Live mode does not require an agent-run home. Store mode skips live capture and
probes sessions under the given home.

## Columns (human table)

Common: `SESSION_ID`, `RUNNER`, `LIVE`, `SENDABLE`, `STATE`, `REASON`, `ITERM`,
`UPDATED`, `WORKSPACE`, plus always-present:

- **`AGENT_RUN`** — process sits under an agent-run parent (`yes`/`no`)
- **`AGENT_SID`** — runner-native session id (live: via `lsof` unless `--fast`;
  store: `meta.runner_session_id`)

`needs_attention` = live + not sendable + not exited. Footer counts sessions,
needs attention, and sendable.

## Soft-fail culture

Probe / iTerm list / live-capture failures print `warning:` on stderr and keep
going when partial output is still useful (exit 0). Hard usage / not-found
errors use `Error:` on stderr and non-zero exit.

## See also

- `list` — filters and JSON; also `kck grok list` for iTerm-hosted Grok ids
- `open` / `send` / `snapshot` / `info` / `status` — Grok ops
