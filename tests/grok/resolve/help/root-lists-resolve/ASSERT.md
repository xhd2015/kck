## Expected

- Root help mentions grok resolve.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "grok resolve") {
		t.Fatalf("root help must mention grok resolve:\n%s", resp.Stdout)
	}
}
```
