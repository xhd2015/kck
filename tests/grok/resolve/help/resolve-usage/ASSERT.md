## Expected

- Resolve help is kck-flavored and documents pid / tab / dry-run / verbose / json.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{
		"kck grok resolve",
		"--pid",
		"--tab",
		"--tab-index",
		"--dry-run",
		"-v",
		"--json",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("resolve help missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "agent-pro grok session resolve") {
		t.Fatalf("resolve help must not be agent-pro branded:\n%s", out)
	}
}
```
