## Expected

- Prefer live agent-run; single sent line; no grok `--resume`.

```go
import (
	"strings"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.AgentRunCalls) != 1 {
		t.Fatalf("AgentRunCalls=%v", resp.AgentRunCalls)
	}
	if len(resp.Opened) != 0 {
		t.Fatalf("Opened=%v", resp.Opened)
	}
	if len(resp.SendCalls) != 0 {
		t.Fatalf("SendCalls=%#v", resp.SendCalls)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
sent to session `+req.SessionID+`
`)
	if strings.Contains(resp.Stdout, "resuming") {
		t.Fatalf("unexpected resume ack:\n%s", resp.Stdout)
	}
}
```
