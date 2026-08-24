## Expected

- Exit 0.
- Stdout contains agent session id `agent-sess-grok-1`.
- Stdout contains runner kind `grok`.
- Footer total ≥ 1 sessions.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "agent-sess-grok-1", "agent session id")
	assertContains(t, resp.Stdout, "grok", "runner kind")
	assertContains(t, resp.Stdout, "1 sessions", "footer total")
}
```
