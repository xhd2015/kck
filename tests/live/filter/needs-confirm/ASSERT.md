## Expected

- Exit 0.
- Stdout contains `agent-sess-a`.
- Stdout does not contain `agent-sess-b`.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "agent-sess-a", "attention row")
	assertNotContains(t, resp.Stdout, "agent-sess-b", "sendable filtered out")
}
```
