## Expected

- Help documents session source and send/open flags.

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	lower := strings.ToLower(resp.Stdout)
	for _, want := range []string{
		"--session-id", "--tab", "--no-submit", "--focus", "--no-ctrl-u", "--open", "--cron",
		"--enter", "--up", "--ctrl-c", "--esc", "--text",
	} {
		if !strings.Contains(lower, want) {
			t.Fatalf("send help must document %q:\n%s", want, resp.Stdout)
		}
	}
}
```
