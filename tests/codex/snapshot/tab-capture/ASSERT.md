## Expected

- Captures tab pane text to stdout.

```go
import "github.com/xhd2015/doctest/assert"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assert.Output(t, resp.Stdout, `---
version: 3
---
kck tab pane
`)
}
```
