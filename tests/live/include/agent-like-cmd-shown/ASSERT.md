## Expected

- Exit 0.
- Row appears: stdout contains iTerm session id `iterm-uuid-idle-mark` **or**
  runner/token `mark` with workspace `/ws/mark`.
- Footer not zero sessions.

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
	out := resp.Stdout
	hasID := strings.Contains(out, "iterm-uuid-idle-mark")
	hasMark := strings.Contains(out, "mark")
	hasWS := strings.Contains(out, "/ws/mark")
	if !(hasID || (hasMark && hasWS)) {
		t.Fatalf("want agent-like mark row (id or mark+/ws/mark); got:\n%s", out)
	}
	if strings.Contains(out, "0 sessions · 0 needs attention · 0 sendable") {
		t.Fatalf("footer must not be all zeros when mark included; got:\n%s", out)
	}
}
```
