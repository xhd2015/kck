## Expected

- `kck codex --help` lists list.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if !strings.Contains(resp.Stdout, "list") {
		t.Fatalf("codex help must list list:\n%s", resp.Stdout)
	}
}
```
