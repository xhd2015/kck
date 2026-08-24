## Expected

- Success; sent line; SendText once.

```go
import "github.com/xhd2015/doctest/assert"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.SendCalls) != 1 || resp.SendCalls[0].Text != "hello" {
		t.Fatalf("SendCalls = %#v", resp.SendCalls)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
sent to session `+req.SessionID+`
`)
}
```
