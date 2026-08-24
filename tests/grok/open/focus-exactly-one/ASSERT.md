## Expected

- Focuses sole tab; does not open a window.

## Expected Output

```text
focused: window 3, tab 1
```

```go
import "github.com/xhd2015/doctest/assert"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.Focused) != 1 || resp.Focused[0] != "w2t1p0" {
		t.Fatalf("Focused=%v", resp.Focused)
	}
	if len(resp.Opened) != 0 {
		t.Fatalf("Opened=%v", resp.Opened)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
focused: window 3, tab 1
`)
}
```
