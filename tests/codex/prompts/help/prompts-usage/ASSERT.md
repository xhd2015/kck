## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{"kck codex prompts", "--first", "--grep", "--this-window"} {
		if !strings.Contains(out, want) {
			t.Fatalf("codex prompts help missing %q:\n%s", want, out)
		}
	}
}
```
