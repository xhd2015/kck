## Expected

- Exit 0.
- Header includes `AGENT_SID`.
- Stdout contains `agent-sess-grok-1` (SESSION_ID prefer + AGENT_SID).
- JSON not used; human table only.

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
	assertTrailingNewline(t, resp.Stdout, "live stdout")
	assertContains(t, strings.ToUpper(resp.Stdout), "AGENT_SID", "header")
	assertContains(t, resp.Stdout, "agent-sess-grok-1", "agent session id")
	// Prefer: appears at least twice (SESSION_ID + AGENT_SID cells) when both set.
	if c := strings.Count(resp.Stdout, "agent-sess-grok-1"); c < 1 {
		t.Fatalf("want agent-sess-grok-1 at least once, got %d", c)
	}
}
```
