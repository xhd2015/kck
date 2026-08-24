## Expected

- Exit 0.
- Row for `agent-sess-grok-1` appears with LIVE yes and SENDABLE no semantics
  (stdout contains session id; footer `1 needs attention`; STATE shows `busy`).
- Footer sendable count is 0 for this single busy row.

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
	out := resp.Stdout
	assertContains(t, out, "agent-sess-grok-1", "session")
	assertContains(t, out, "busy", "STATE busy")
	assertContains(t, out, "1 needs attention", "attention count")
	// SENDABLE no on the data row: prefer line containing session id also has "no"
	// after LIVE yes — flexible: footer ends with 0 sendable.
	if !strings.Contains(out, "0 sendable") {
		t.Fatalf("want 0 sendable in footer for single busy row; got:\n%s", out)
	}
}
```
