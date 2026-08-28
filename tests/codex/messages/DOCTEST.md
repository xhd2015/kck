# `kck codex messages`

Thin CLI over `sessions.RunMessages`: print recent coalesced Codex chat
messages from rollout JSONL (msgfmt-style).

L2 only — injectable `CodexHome` (+ optional `CodexMessagesOpts`).

# DSN (Domain Specific Notion)

## Participants

- **Caller** — `run.MainWith` with `Args: ["codex","messages",…]`.
- **`kck/run`** — help + `Error:` prefix.
- **`sessions.RunMessages`** — load rollout, page, format.

## Behaviors

- Help documents messages flags.
- Missing / unknown session → `Error:`.
- Empty rollout → `(no messages)`.
- Kinds order user → thinking → tool → assistant.
- `--limit` keeps latest; `--grep` filters; `--json` schema.

## Version

0.0.1

## Decision Tree

```text
codex/messages/
├── help/
│   ├── root-lists-messages/
│   ├── codex-usage/
│   └── messages-usage/
├── missing-session-source/
├── unknown-session/
├── empty/
├── kinds-order/
├── limit-keeps-latest/
├── json/
└── grep-filter/
```

## How to Run

```sh
doctest vet ./tests/codex/messages
doctest test ./tests/codex/messages
```

```go
import (
	"bytes"
	"testing"
	"time"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/doctest/session"

	"kck/run"
)

type Request struct {
	Args      []string
	CodexHome string
	TempDir   string
	SessionID string
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
	msgOpts := &sessions.MessagesOpts{Loc: time.UTC}
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:              append([]string(nil), req.Args...),
		Stdout:            &stdout,
		Stderr:            &stderr,
		CodexHome:         req.CodexHome,
		CodexMessagesOpts: msgOpts,
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
