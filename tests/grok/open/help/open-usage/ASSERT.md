## Expected

- Open help is kck-flavored and documents tab source + resume + index.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{"kck grok open", "--tab", "--tab-index", "--index", "grok --resume", "--no-agent-run", "--dry-run", "agent-run"} {
		if !strings.Contains(out, want) {
			t.Fatalf("open help missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "agent-pro grok session open") {
		t.Fatalf("open help must not be agent-pro branded:\n%s", out)
	}
}
```
