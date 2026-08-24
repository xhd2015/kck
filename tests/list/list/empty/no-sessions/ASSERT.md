## Expected

- Exit 0.
- Stdout contains footer `0 sessions · 0 needs attention · 0 sendable`.
- No session id rows required; may include column header.
- Trailing newline on stdout.

## Errors

- None.

## Exit Code

0

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertTrailingNewline(t, resp.Stdout, "list stdout")
	assertContains(t, resp.Stdout, "0 sessions · 0 needs attention · 0 sendable", "empty footer")
}
```
