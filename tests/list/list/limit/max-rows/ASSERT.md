## Expected

- Exit 0.
- Newest `sess-sendable` present.
- Older `sess-attention` and `sess-exited` not listed.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "sess-sendable", "newest kept")
	assertNotContains(t, resp.Stdout, "sess-attention", "limited out")
	assertNotContains(t, resp.Stdout, "sess-exited", "limited out")
}
```
