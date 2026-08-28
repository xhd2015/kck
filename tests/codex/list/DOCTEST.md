# `kck codex list`

Thin CLI over `agent-pro` `codex/sessions.RunListLive`: list Codex session ids
that currently have a hosting iTerm2 tab.

L2 only — injectable `CodexHome` + `CodexListLiveOpts`. No live iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["codex","list",…]`.
- **`kck/run`** — dispatches `codex list`, prints kck-flavored help.
- **`sessions.RunListLive`** — core discovery (agent-pro).

## Behaviors

- Root / codex help mention `list`.
- One host → table row with sid + iterm.
- Empty → `0 sessions`.

## Version

0.0.1

## Decision Tree

```text
codex/list/
├── help/
│   ├── root-lists-list/
│   ├── codex-usage/
│   └── list-usage/
├── one-host/
└── empty/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-list/` | `kck -h` mentions `codex list`. |
| `help/codex-usage/` | `kck codex --help` lists `list`. |
| `help/list-usage/` | `kck codex list --help` documents `--json` / `--limit`. |
| `one-host/` | Table contains sid + `w=` iterm. |
| `empty/` | `0 sessions`. |

## How to Run

```sh
doctest vet ./tests/codex/list
doctest test ./tests/codex/list
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"

	"kck/run"
)

type Request struct {
	Args       []string
	CodexHome  string
	TempDir    string
	ProjectDir string
	SessionID  string
	Procs      []sessions.FocusProc
	OpenFiles  map[int][]string
	ITerm      []iterm2.SessionRef
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
	fake := &sessions.ListLiveFake{
		FocusFake: sessions.FocusFake{
			Procs:     req.Procs,
			OpenFiles: req.OpenFiles,
			ITerm:     req.ITerm,
		},
	}
	var stdout, stderr bytes.Buffer
	err := run.MainWith(run.Options{
		Args:              append([]string(nil), req.Args...),
		Stdout:            &stdout,
		Stderr:            &stderr,
		CodexHome:         req.CodexHome,
		CodexListLiveOpts: fake.ListLiveOpts(),
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
