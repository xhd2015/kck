## Expected

- Exit 0.
- STATE includes `idle`.
- Footer includes `1 sendable`.
- Footer needs attention is 0.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "idle", "STATE idle")
	assertContains(t, out, "1 sendable", "sendable count")
	assertContains(t, out, "0 needs attention", "no attention when idle sendable")
}
```
