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
		"kck grok prompts",
		"--first",
		"--main",
		"--head",
		"--tail",
		"--grep",
		"--this-window",
		"--this-space",
		"--this-tab",
		"--session-id",
		"--tab",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("prompts help missing %q:\n%s", want, out)
		}
	}
}
```
