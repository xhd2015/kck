## Expected

- Root help mentions grok send.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "grok send") {
		t.Fatalf("root help must mention grok send:\n%s", resp.Stdout)
	}
	if strings.Contains(lower, "--send message") || strings.Contains(resp.Stdout, "not implemented yet") {
		t.Fatalf("root help must not advertise removed list --send stub:\n%s", resp.Stdout)
	}
}
```
