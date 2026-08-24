## Expected

- `--tab 2` sends to resolved pane.

```go
import (
	"strings"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.SendCalls) != 1 {
		t.Fatalf("SendCalls = %#v", resp.SendCalls)
	}
	if !strings.Contains(resp.SendCalls[0].SessionID, "TAB2-UUID") {
		t.Fatalf("session = %q", resp.SendCalls[0].SessionID)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
sent to session `+fixtureKckTabSendSessionID+`
`)
}
```
