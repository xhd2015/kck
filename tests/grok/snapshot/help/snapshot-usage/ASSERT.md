## Expected

- Snapshot help is kck-flavored and documents tab source + Contents + index.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{"kck grok snapshot", "--tab", "--tab-index", "--index", "--json", "--dry-run", "--iterm", "agent-run", "source"} {
		if !strings.Contains(out, want) {
			t.Fatalf("snapshot help missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "agent-pro grok session snapshot") {
		t.Fatalf("snapshot help must not be agent-pro branded:\n%s", out)
	}
}
```
