## Expected

- Exit 0.
- No `w=42` (resolution skipped).
- Session still listed.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "iterm-s", "session")
	assertNotContains(t, resp.Stdout, "w=42", "--no-iterm must skip match")
}
```
