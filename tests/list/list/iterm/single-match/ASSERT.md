## Expected

- Exit 0.
- Stdout contains ITERM cell `w=42 t=3` (exact token).
- No `(+` multi suffix for this row.

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "w=42 t=3", "ITERM single")
	if strings.Contains(resp.Stdout, "w=42 t=3(+") {
		t.Fatalf("single match must not use (+N); stdout:\n%s", resp.Stdout)
	}
}
```
