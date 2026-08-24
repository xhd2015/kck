## Expected

- Prefer agent-run resume window for managed Grok id.

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
	if len(resp.Opened) != 1 || !strings.Contains(resp.Opened[0], "agent-run") {
		t.Fatalf("Opened=%v", resp.Opened)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
opened: new window; agent-run resume ar-kck-open
`)
}
```
