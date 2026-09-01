# `kck codex prompts`

Thin CLI over `codex/sessions.RunPrompts`.

## Version

0.0.1

## How to Run

```sh
doctest vet ./tests/codex/prompts
doctest test ./tests/codex/prompts
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
	CodexHome string
	TempDir   string
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
		Args:      append([]string(nil), req.Args...),
		Stdout:    &stdout,
		Stderr:    &stderr,
		CodexHome: req.CodexHome,
	})
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String()}
	if err != nil {
		resp.ErrText = err.Error()
		resp.ExitCode = 1
	}
	return resp, nil
}
```
