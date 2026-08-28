## Expected

- Root help mentions `codex list`.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if !strings.Contains(resp.Stdout, "codex list") {
		t.Fatalf("root help must mention codex list:\n%s", resp.Stdout)
	}
}
```
