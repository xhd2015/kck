## Expected

- Exit 0.
- Stdout contains `agent-sess-b`.
- Stdout does not contain `agent-sess-a`.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "agent-sess-b", "sendable row")
	assertNotContains(t, resp.Stdout, "agent-sess-a", "busy filtered out")
}
```
