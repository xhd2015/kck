## Expected

- Exit 0.
- Stdout lists `only-here`.
- Does not invent sessions from real user `~/.agent-run` (id must be only the seed).

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "only-here", "isolated home session")
	assertContains(t, resp.Stdout, "1 sessions", "count under isolated home")
	_ = req
}
```
