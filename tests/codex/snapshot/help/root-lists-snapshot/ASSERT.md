## Expected

- Root help mentions codex snapshot.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "codex snapshot") {
		t.Fatalf("root help must mention codex snapshot:\n%s", resp.Stdout)
	}
}
```
