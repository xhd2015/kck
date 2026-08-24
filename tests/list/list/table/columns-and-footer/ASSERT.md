## Expected

- Header line includes: SESSION_ID, RUNNER, LIVE, SENDABLE, STATE, REASON,
  AGENT_RUN, AGENT_SID, ITERM, UPDATED, WORKSPACE (order as specified;
  whitespace-flexible).
- Footer: `2 sessions · 1 needs attention · 1 sendable` for twoSessionFixture
  (s-new attention, s-old sendable).
- Exit 0; trailing newline.

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
	assertTrailingNewline(t, resp.Stdout, "list stdout")
	upper := strings.ToUpper(resp.Stdout)
	for _, col := range []string{
		"SESSION_ID", "RUNNER", "LIVE", "SENDABLE", "STATE", "REASON",
		"AGENT_RUN", "AGENT_SID", "ITERM", "UPDATED", "WORKSPACE",
	} {
		if !strings.Contains(upper, col) {
			t.Fatalf("header/columns missing %q; stdout:\n%s", col, resp.Stdout)
		}
	}
	assertContains(t, resp.Stdout, "2 sessions · 1 needs attention · 1 sendable", "footer")
}
```
