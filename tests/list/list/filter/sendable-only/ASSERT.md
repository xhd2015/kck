## Expected

- Exit 0.
- `sess-sendable` present; `sess-attention` and `sess-exited` absent.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "sess-sendable", "sendable session")
	assertNotContains(t, resp.Stdout, "sess-attention", "attention filtered out")
	assertNotContains(t, resp.Stdout, "sess-exited", "exited filtered out")
}
```
