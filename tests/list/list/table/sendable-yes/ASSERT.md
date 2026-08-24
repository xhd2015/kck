## Expected

- Exit 0.
- Row for `ready-1` present.
- Stdout includes sendable yes token near the row (assert contains `yes` and
  footer `1 sendable`).
- Footer: `1 sessions · 0 needs attention · 1 sendable`.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "ready-1", "session")
	assertContains(t, resp.Stdout, "1 sessions · 0 needs attention · 1 sendable", "footer")
}
```
