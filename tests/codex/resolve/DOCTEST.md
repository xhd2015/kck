# `kck codex resolve`

Thin CLI over `agent-pro` `codex/sessions.RunResolve`: resolve a Codex session
id by walking ancestors to the nearest codex runner (open-file paths), or from
a sibling iTerm2 tab (`--tab` / `--tab-index`).

L2 only — injectable `CodexHome` + `CodexResolveOpts`. No live ps / lsof / iTerm.

# DSN (Domain Specific Notion)

## Participants

- **Caller** — harness via `run.MainWith` with `Args: ["codex","resolve",…]`.
- **`kck/run`** — dispatches `codex resolve`, prints kck-flavored help, prefixes
  errors with `Error:`.
- **`sessions.RunResolve`** — core ancestor / tab resolve (agent-pro).

## Behaviors

- Root help mentions `codex resolve`.
- `kck codex --help` / `kck codex resolve --help` document `--pid` / tab flags.
- Default success: bare session id on stdout.
- `--tab N` resolves sibling tab codex id.
- Ancestor miss → `Error: no ancestor codex`.
- `--pid` + `--tab` → mutual exclusion `Error:`.
- `--dry-run` prints `[dry-run]` plan (no bare-id-only shape).

## Version

0.0.1

## Decision Tree

```text
codex/resolve/
├── help/
│   ├── root-lists-resolve/
│   ├── codex-usage/
│   └── resolve-usage/
├── hit/
│   ├── bare/
│   ├── tab/
│   └── dry-run/
├── miss/
│   └── no-ancestor/
└── flags/
    └── pid-and-tab/
```

## Test Index

| Leaf | Contract |
|---|---|
| `help/root-lists-resolve/` | `kck -h` mentions `codex resolve`. |
| `help/codex-usage/` | `kck codex --help` lists `resolve`. |
| `help/resolve-usage/` | `kck codex resolve --help` documents flags; kck-branded. |
| `hit/bare/` | Default stdout is bare session id. |
| `hit/tab/` | `--tab 2` resolves sibling tab codex id. |
| `hit/dry-run/` | `[dry-run]` plan; would resolve fixture id. |
| `miss/no-ancestor/` | `Error: no ancestor codex`. |
| `flags/pid-and-tab/` | `Error:` mutual exclusion. |

## How to Run

```sh
doctest vet ./tests/codex/resolve
doctest test ./tests/codex/resolve
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/agent-pro/pkgs/procresolve"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"

	"kck/run"
)

type Request struct {
	Args             []string
	PID              int
	Procs            []FixtureProc
	FocusProcs       []sessions.FocusProc
	OpenFiles        map[int][]string
	ITerm            []iterm2.SessionRef
	CurrentSessionID string
	ControllingTTY   string
	CodexHome        string
	TempDir          string
}

type FixtureProc struct {
	PID  int
	PPID int
	Cmd  string
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
	procs := make([]procresolve.Proc, 0, len(req.Procs))
	for _, p := range req.Procs {
		procs = append(procs, procresolve.Proc{PID: p.PID, PPID: p.PPID, Cmd: p.Cmd})
	}
	snap := procs
	files := req.OpenFiles
	if files == nil {
		files = map[int][]string{}
	}
	focusSnap := append([]sessions.FocusProc(nil), req.FocusProcs...)
	itermSnap := append([]iterm2.SessionRef(nil), req.ITerm...)
	resolveOpts := &sessions.ResolveOpts{
		PID: req.PID,
		ListProcs: func() []procresolve.Proc {
			return append([]procresolve.Proc(nil), snap...)
		},
		Lsof: func(pid int) []string {
			return files[pid]
		},
		CodexHome: req.CodexHome,
		ListFocusProcs: func() []sessions.FocusProc {
			return append([]sessions.FocusProc(nil), focusSnap...)
		},
		ListITerm: func() ([]iterm2.SessionRef, error) {
			return append([]iterm2.SessionRef(nil), itermSnap...), nil
		},
		CurrentSessionID: func() string { return req.CurrentSessionID },
		ControllingTTY:   func() string { return req.ControllingTTY },
		AncestorTTYs:     func() []string { return nil },
	}
	err := run.MainWith(run.Options{
		Args:             append([]string(nil), req.Args...),
		Stdout:           &stdout,
		Stderr:           &stderr,
		CodexHome:        req.CodexHome,
		CodexResolveOpts: resolveOpts,
	})
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
