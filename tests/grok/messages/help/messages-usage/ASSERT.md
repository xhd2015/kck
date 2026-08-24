## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{
		"kck grok messages",
		"--limit",
		"--offset-from-end",
		"--offset-from-end 32",
		"4096",
		"128",
		"512",
		"8192",
		"--json",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("messages help missing %q:\n%s", want, out)
		}
	}
}
```
