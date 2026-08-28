## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	lower := strings.ToLower(out)
	for _, want := range []string{"kck codex status", "--no-pid", "--json"} {
		if !strings.Contains(out, want) {
			t.Fatalf("status help missing %q:\n%s", want, out)
		}
	}
	for _, want := range []string{"path", "rollout"} {
		if !strings.Contains(lower, want) {
			t.Fatalf("status help missing %q:\n%s", want, out)
		}
	}
}
```
