## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "keep zebra", "stdout")
	if strings.Contains(resp.Stdout, "[user] : drop") {
		t.Fatalf("drop should be filtered:\n%s", resp.Stdout)
	}
}
```
