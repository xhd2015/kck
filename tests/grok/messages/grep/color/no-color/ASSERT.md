## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if strings.ContainsRune(resp.Stdout, '\x1b') {
		t.Fatalf("--no-color must emit no ANSI:\n%q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "[assistant] : a1", "stdout")
}
```
