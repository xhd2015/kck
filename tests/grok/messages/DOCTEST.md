# `kck grok messages`

Thin CLI over `sessions.RunMessages`: print recent coalesced Grok chat
messages (msgfmt-style) with per-kind rune caps and `--offset-from-end` paging.

L2 only — injectable `GrokHome` (+ optional `GrokMessagesOpts` for tab resolve).

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["grok","messages",…]`.
- **`kck/run`** — dispatches `grok messages`, prints kck-flavored help, prefixes
  errors with `Error:`.
- **`sessions.RunMessages` / `Messages`** — load `updates.jsonl`, page, format.

## Behaviors

- Root / grok / messages help document `messages`, `--limit`, `--offset-from-end`
  (including the `32` next-page example), `--grep`, `--color`, `--no-color`.
- Missing session source / unknown session → `Error:`.
- Empty / missing updates → `(no messages)`.
- Default `--limit` 32; `--offset-from-end` skips newest before limit.
- `--grep` (repeatable, AND) filters message bodies before paging; no hits →
  `(no matching messages)`.
- Text: msgfmt `Chat history (showing K of N):` + `[kind] : body` (streamed).
- `--color` / `--no-color` / auto; match spans bold-red when grepping.
- `--json`: `{session_id,total,offset_from_end,limit,messages[]}` (no ANSI;
  `total` is post-grep).

## Version

0.0.2

## Decision Tree

```text
grok/messages/
├── help/
│   ├── root-lists-messages/
│   ├── grok-usage/
│   └── messages-usage/
├── missing-session-source/
├── unknown-session/
├── empty/
├── limit-keeps-latest/
├── offset-pages/
├── offset-past-end/
├── kinds-order/
├── timestamp-prefix/
├── json/
└── grep/
    ├── filter/
    │   ├── single/
    │   ├── and/
    │   ├── no-match/
    │   └── with-limit/
    ├── errors/
    │   ├── empty-pattern/
    │   └── color-conflict/
    ├── color/
    │   ├── force-on/
    │   └── no-color/
    └── json/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-messages/` | `kck -h` mentions `grok messages`. |
| `help/grok-usage/` | `kck grok --help` lists `messages`. |
| `help/messages-usage/` | documents `--limit`, `--offset-from-end 32`, `--grep`, color flags, caps. |
| `missing-session-source/` | `Error:` usage. |
| `unknown-session/` | `Error: grok session not found`. |
| `empty/` | `(no messages)`. |
| `limit-keeps-latest/` | `--limit 2` keeps newest pair; header `showing 2 of N`. |
| `offset-pages/` | `--offset-from-end 2 --limit 2` is the prior page. |
| `offset-past-end/` | offset ≥ total → `(no messages)`. |
| `kinds-order/` | user / thinking / tool / assistant labels in order. |
| `timestamp-prefix/` | local `[YYYY-MM-DD HH:MM:SS]` from wire times. |
| `json/` | valid JSON; fields total/offset/limit/messages (+ timestamp). |
| `grep/filter/single/` | one `--grep` keeps matching bodies only. |
| `grep/filter/and/` | two `--grep`s require both substrings in the same message. |
| `grep/filter/no-match/` | `(no matching messages)`. |
| `grep/filter/with-limit/` | `--limit` applies after grep (newest matches). |
| `grep/errors/empty-pattern/` | empty `--grep` → `Error:`. |
| `grep/errors/color-conflict/` | `--color` + `--no-color` → `Error:`. |
| `grep/color/force-on/` | `--color` highlights match spans (ANSI). |
| `grep/color/no-color/` | `--no-color` emits no ANSI. |
| `grep/json/` | filtered JSON; no ANSI; `total` = match count. |

## How to Run

```sh
doctest vet ./tests/grok/messages
doctest test ./tests/grok/messages
```

```go
import (
	"bytes"
	"testing"
	"time"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/doctest/session"

	"kck/run"
)

type Request struct {
	Args      []string
	GrokHome  string
	TempDir   string
	SessionID string
	Loc       *time.Location // nil → time.Local; injectable for timestamp asserts
}

type Response struct {
	Stdout   string
	Stderr   string
	ErrText  string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	var stdout, stderr bytes.Buffer
	msgOpts := &sessions.MessagesOpts{}
	if req.Loc != nil {
		msgOpts.Loc = req.Loc
	}
	err := run.MainWith(run.Options{
		Args:             append([]string(nil), req.Args...),
		Stdout:           &stdout,
		Stderr:           &stderr,
		GrokHome:         req.GrokHome,
		GrokMessagesOpts: msgOpts,
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
