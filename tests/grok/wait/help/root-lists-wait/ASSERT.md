## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if !strings.Contains(strings.ToLower(resp.Stdout), "grok wait") {
		t.Fatalf("root help must mention grok wait:\n%s", resp.Stdout)
	}
}
```
