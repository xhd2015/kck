## Expected

- Exit 0.
- Both `agent-sess-a` and `agent-sess-b` in stdout.
- Footer exactly: `2 sessions · 1 needs attention · 1 sendable`.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "agent-sess-a", "agent a")
	assertContains(t, resp.Stdout, "agent-sess-b", "agent b")
	assertContains(t, resp.Stdout, "2 sessions · 1 needs attention · 1 sendable", "footer")
}
```
