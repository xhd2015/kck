## Expected

- Exit 0.
- Exactly one of `agent-sess-a` or `agent-sess-b` appears (not both).
- Footer total is 1 session.

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
	hasA := strings.Contains(out, "agent-sess-a")
	hasB := strings.Contains(out, "agent-sess-b")
	if hasA == hasB {
		// both true or both false
		t.Fatalf("want exactly one of agent-sess-a/b; a=%v b=%v\n%s", hasA, hasB, out)
	}
	assertContains(t, out, "1 sessions", "limit footer")
}
```
