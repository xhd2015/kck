## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := strings.ToLower(resp.Stdout)
	assertContains(t, out, "timeout", "stdout")
	assertContains(t, out, "updates.jsonl", "stdout")
	assertContains(t, out, "turn_completed", "stdout")
}
```
