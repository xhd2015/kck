## Expected

- Exit 0.
- Stdout contains `sess-attention`.
- Stdout does not contain `sess-sendable` or `sess-exited`.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "sess-attention", "attention session")
	assertNotContains(t, resp.Stdout, "sess-sendable", "sendable filtered out")
	assertNotContains(t, resp.Stdout, "sess-exited", "exited filtered out")
}
```
