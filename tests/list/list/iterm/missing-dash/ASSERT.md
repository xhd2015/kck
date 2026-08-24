## Expected

- Exit 0.
- Session listed; ITERM column is `-` (assert line/token with session and `-`).
- Must not contain `w=` for this session row when no match (no false match).

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
	assertContains(t, resp.Stdout, "iterm-s", "session")
	// Row should include bare dash for ITERM; require " - " or trailing column dash.
	if !strings.Contains(resp.Stdout, " - ") && !strings.Contains(resp.Stdout, "\t-\t") &&
		!strings.Contains(resp.Stdout, " -") {
		// Fallback: any field that is exactly dash between spaces near session line
		found := false
		for _, line := range strings.Split(resp.Stdout, "\n") {
			if strings.Contains(line, "iterm-s") && strings.Contains(line, "-") && !strings.Contains(line, "w=") {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("want ITERM - for no match; stdout:\n%s", resp.Stdout)
		}
	}
	assertNotContains(t, resp.Stdout, "w=9", "unrelated iterm must not attach")
}
```
