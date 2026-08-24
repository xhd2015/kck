## Expected

- Usage documents `--json` and `--limit`.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{"Usage: kck grok list", "--json", "--limit"} {
		if !strings.Contains(out, want) {
			t.Fatalf("list help missing %q:\n%s", want, out)
		}
	}
}
```
