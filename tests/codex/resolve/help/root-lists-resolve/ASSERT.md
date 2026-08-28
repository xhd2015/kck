## Expected

- Root help mentions codex resolve.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "codex resolve") {
		t.Fatalf("root help must mention codex resolve:\n%s", resp.Stdout)
	}
}
```
