## Expected

- Two stdout lines; opener + SendText.

```go
import (
	"strings"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.Opened) != 1 {
		t.Fatalf("Opened = %v, want 1", resp.Opened)
	}
	if !strings.Contains(resp.Opened[0], "--resume") {
		t.Fatalf("Opened missing resume: %q", resp.Opened[0])
	}
	if len(resp.SendCalls) != 1 {
		t.Fatalf("SendCalls = %#v", resp.SendCalls)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
opened: new window; resuming `+req.SessionID+`
sent to session `+req.SessionID+`
`)
}
```
