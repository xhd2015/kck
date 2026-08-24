## Expected

- Exit 0.
- Footer includes `0 sessions · 0 needs attention · 0 sendable`.
- Trailing newline on stdout.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertTrailingNewline(t, resp.Stdout, "stdout")
	assertContains(t, resp.Stdout, "0 sessions · 0 needs attention · 0 sendable", "footer")
}
```
