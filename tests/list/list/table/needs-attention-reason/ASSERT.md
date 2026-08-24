## Expected

- Exit 0.
- Stdout contains session `attn-1` and reason text `awaiting confirmation`.
- Footer includes `1 needs attention`.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "attn-1", "session id")
	assertContains(t, resp.Stdout, "awaiting confirmation", "REASON")
	assertContains(t, resp.Stdout, "1 needs attention", "footer attention count")
}
```
