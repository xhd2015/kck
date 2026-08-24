## Expected

- Exit 0.
- Stdout does **not** contain `iterm-uuid-plain-zsh`.
- Footer is `0 sessions · 0 needs attention · 0 sendable`.

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertNotContains(t, resp.Stdout, "iterm-uuid-plain-zsh", "plain shell omitted")
	assertContains(t, resp.Stdout, "0 sessions · 0 needs attention · 0 sendable", "footer")
}
```
