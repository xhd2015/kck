## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := strings.ToLower(resp.Stdout)
	if !strings.Contains(out, "wait") {
		t.Fatalf("codex help must list wait:\n%s", resp.Stdout)
	}
}
```
