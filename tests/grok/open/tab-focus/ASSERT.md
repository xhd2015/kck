## Expected

- `kck grok open --tab 2` focuses the resolved tab.

## Expected Output

```text
focused: window 100, tab 2
```

```go
import "github.com/xhd2015/doctest/assert"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.Focused) != 1 || resp.Focused[0] != "w0t2p0:TAB2-UUID" {
		t.Fatalf("Focused=%v", resp.Focused)
	}
	if len(resp.Opened) != 0 {
		t.Fatalf("Opened=%v", resp.Opened)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
focused: window 100, tab 2
`)
}
```
