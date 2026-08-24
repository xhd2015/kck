## Expected

- Exit 0.
- ITERM shows `w=1 t=2(+1)` (first match + one extra).

## Exit Code

0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "w=1 t=2(+1)", "ITERM multi")
}
```
