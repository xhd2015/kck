## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "Chat history (showing last 2 of", "stdout")
	assertContains(t, out, "[user] : u2", "stdout")
	assertContains(t, out, "[assistant] : a2", "stdout")
	if strings.Contains(out, "u0") || strings.Contains(out, "a0") || strings.Contains(out, "u1") {
		t.Fatalf("limit 2 must not include older messages:\n%s", out)
	}
}
```
