# `kck grok prompts`

Thin CLI over `sessions.RunPrompts`: list user prompts with `--first`,
`--head` / `--tail`, `--main`, repeatable `--grep` (AND), and live host scopes.

L2 — injectable `GrokHome` (+ optional `GrokPromptsOpts`).

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["grok","prompts",…]`.
- **`kck/run`** — dispatches `grok prompts`, prints kck-flavored help, prefixes
  errors with `Error:`.
- **`sessions.RunPrompts`** — load updates.jsonl user prompts, filter, format.

## Behaviors

- Help documents `--first`, `--main`, `--head`, `--tail`, scopes, `--grep`.
- Per-session clip: `--head N` first N; `--tail N` last N; `--first` ≡ `--head 1`.
- `--head` and `--tail` mutually exclusive; `--first` exclusive with both.
- `--main` (role) skips subagent-class sessions in multi mode.
- Text: `[YYYY-MM-DD HH:MM:SS] body`; omission markers `(...K omitted...)`.

## Version

0.0.2

## Decision Tree

```text
grok/prompts/
├── help/
│   ├── root-lists-prompts/
│   ├── grok-usage/
│   └── prompts-usage/
├── clip/                         # which prompts appear (greatest effect)
│   ├── head/keeps-first-2/
│   ├── tail/keeps-last-2/
│   └── first/equals-head-1/
├── errors/
│   ├── head-and-tail/
│   ├── first-and-head/
│   └── tail-zero/
└── role/
    └── multi-skips-subagent/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-prompts/` | root `-h` mentions `grok prompts` |
| `help/grok-usage/` | `kck grok --help` lists `prompts` |
| `help/prompts-usage/` | documents `--first`, `--main`, `--head`/`--tail`, scopes |
| `clip/head/keeps-first-2/` | `--head 2` → p1,p2 + trailing omission; no p3+ |
| `clip/tail/keeps-last-2/` | `--tail 2` → leading omission + p4,p5 |
| `clip/first/equals-head-1/` | `--first` → only p1 + trailing omission |
| `errors/head-and-tail/` | `--head`+`--tail` → `Error:` |
| `errors/first-and-head/` | `--first`+`--head` → `Error:` |
| `errors/tail-zero/` | `--tail 0` → `Error:` |
| `role/multi-skips-subagent/` | `--main` multi keeps main session only |

## How to Run

```sh
doctest vet ./tests/grok/prompts
doctest test ./tests/grok/prompts
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/doctest/session"

	"kck/run"
)

type Request struct {
	Args      []string
	GrokHome  string
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
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:     append([]string(nil), req.Args...),
		Stdout:   &stdout,
		Stderr:   &stderr,
		GrokHome: req.GrokHome,
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
