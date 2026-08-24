## Expected

- Success (exit 0 / nil product error).
- Stdout documents `--home`, `--json`, `--needs-confirm`, `--sendable`,
  `--no-iterm`, `--fast`.
- Stdout does not advertise removed list `--send` / `--session` stub.
- Stdout ends with trailing newline.

## Errors

- None.

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertTrailingNewline(t, resp.Stdout, "help stdout")
	lower := strings.ToLower(resp.Stdout)
	for _, want := range []string{"--home", "--json", "--needs-confirm", "--sendable", "--no-iterm", "--fast"} {
		if !strings.Contains(lower, strings.ToLower(want)) {
			t.Fatalf("help must document %q; got:\n%s", want, resp.Stdout)
		}
	}
	if strings.Contains(lower, "not implemented yet") {
		t.Fatalf("help must not mention send stub; got:\n%s", resp.Stdout)
	}
	// List mode must not document removed --send MESSAGE / --session ID.
	if strings.Contains(resp.Stdout, "--send MESSAGE") || strings.Contains(resp.Stdout, "--session ID") {
		t.Fatalf("help must not document removed --send/--session; got:\n%s", resp.Stdout)
	}
}
```
