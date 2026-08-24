## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{"kck grok info", "Active", "--no-pid"} {
		if !strings.Contains(out, want) {
			t.Fatalf("info help missing %q:\n%s", want, out)
		}
	}
}
```
