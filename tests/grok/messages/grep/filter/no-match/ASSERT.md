## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if !strings.Contains(resp.Stdout, "(no matching messages)") {
		t.Fatalf("want (no matching messages):\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "(no messages)\n") && !strings.Contains(resp.Stdout, "matching") {
		t.Fatalf("must distinguish no-match from empty session:\n%s", resp.Stdout)
	}
}
```
