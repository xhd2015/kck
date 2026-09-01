---
name: kck-pickup-a-session
description: >-
  Read an existing Grok or Codex session's messages with kck grok|codex
  messages <session-id> and fulfill the user's stated intent. Requires an
  explicit session-id and intent; pause and ask if either is missing. Triggers:
  /kck-pickup-a-session, /kck-from-a-session, read session messages, from that
  session, what did that session say, pickup from session.
---

# kck-pickup-a-session

Given an explicit **session-id** and a stated **intent**, fetch that session's
messages with `kck grok messages` or `kck codex messages` and fulfill the intent.

This skill is owned by **kck** and staged into new empty sessions by
`kck grok pickup` / `kck codex pickup` (not a native fork).

Patrol / kick live agents → `track-user-working-on`. Flag encyclopedia →
`kck skill --show messages`.

## When / not

| Use | Not |
|-----|-----|
| User named a session-id and what to do with its messages | Inventing a sid from `kck` / `kck grok list` |
| Summarize, extract decisions/APIs, continue from prior turns | Latency analysis → `analyse-agent-session-tool-latency` |
| Narrow with `--grep` then answer the ask | Dumping the full transcript when a page + grep suffice |
| Staged via `kck grok|codex pickup` (new empty session + draft) | Native fork (`--fork-session`) — pickup is reference-only |

## Hard gates

1. **Session-id** — must be explicit in the user message. If missing or ambiguous, **stop and ask**. Do not pick from a live list.
2. **Intent** — what to do with those messages. If unclear, ask once.

## Runner choice

| Context | Command |
|---------|---------|
| Grok base / `kck grok pickup` draft / user said Grok | `kck grok messages <session-id> …` |
| Codex base / `kck codex pickup` draft / user said Codex | `kck codex messages <session-id> …` |
| Unclear which runner | Ask once; do not guess from `kck` list |

## Procedure

1. Confirm both gates (or ask). Pick grok vs codex (table above).
2. Fetch (prefer latest page — from the bottom):

```bash
kck grok messages <session-id> --limit 32
kck grok messages <session-id> --grep A --grep B --limit 32
kck grok messages <session-id> --offset-from-end 32 --limit 32
kck grok messages <session-id> --grep A --json

kck codex messages <session-id> --limit 32
kck codex messages <session-id> --grep A --grep B --limit 32
kck codex messages <session-id> --offset-from-end 32 --limit 32
kck codex messages <session-id> --grep A --json
```

- `--grep` is repeatable (AND, case-insensitive literal); applied before paging.
- Prefer `--json` when paging or extracting (`total` is post-grep).
- Page older with `--offset-from-end` only as needed for the intent.

3. Fulfill the intent from the fetched messages. Shape the reply to the ask —
   no fixed report template. Do not paste secrets from the transcript.
4. More flags / paging detail: `kck skill --show messages`.

## Wrong → correct

| Wrong | Correct |
|-------|---------|
| Guess sid from `kck grok list` | Ask for the session-id |
| Run messages with no stated goal | Ask what to do with the messages |
| Dump every turn by default | `--limit 32`, `--grep` when the ask names terms, page only if needed |
| Use `kck grok messages` for a Codex session | `kck codex messages` (and the reverse) |

## Pitfalls

- `--grep` filters **message body**; patterns are AND across repeats.
- Cap truncation is per kind (see `kck skill --show messages`); re-fetch with
  `--json` / another page if a critical span was clipped.
- This skill does not send, open, or focus; use other kck flows if the user
  asks for those actions separately. To *start* a new empty session with this
  skill pre-staged, use `kck grok pickup` / `kck codex pickup`.
