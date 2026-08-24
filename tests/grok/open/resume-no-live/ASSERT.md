## Expected

- Resumes in new window; no focus.

```go
import (
	"strings"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.Focused) != 0 {
		t.Fatalf("Focused=%v", resp.Focused)
	}
	if len(resp.Opened) != 1 {
		t.Fatalf("Opened=%v", resp.Opened)
	}
	entry := resp.Opened[0]
	if !strings.Contains(entry, "--resume") || !strings.Contains(entry, req.SessionID) {
		t.Fatalf("opened follow-up: %q", entry)
	}
	if strings.Contains(entry, "--fork-session") {
		t.Fatalf("must not fork: %q", entry)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
opened: new window; resuming `+req.SessionID+`
`)
}
```
