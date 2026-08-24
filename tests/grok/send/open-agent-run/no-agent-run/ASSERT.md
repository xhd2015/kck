## Expected

- `--no-agent-run` forces bare resume + SendText.

```go
import (
	"strings"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.AgentRunCalls) != 0 {
		t.Fatalf("AgentRunCalls=%v", resp.AgentRunCalls)
	}
	if len(resp.Opened) != 1 || !strings.Contains(resp.Opened[0], "--resume") {
		t.Fatalf("Opened=%v", resp.Opened)
	}
	if len(resp.SendCalls) != 1 {
		t.Fatalf("SendCalls=%#v", resp.SendCalls)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
opened: new window; resuming `+req.SessionID+`
sent to session `+req.SessionID+`
`)
}
```
