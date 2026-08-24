## Expected

- Root help mentions grok snapshot.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "grok snapshot") {
		t.Fatalf("root help must mention grok snapshot:\n%s", resp.Stdout)
	}
}
```
