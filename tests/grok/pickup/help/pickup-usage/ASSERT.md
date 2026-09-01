# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "--session-id", "stdout")
	assertContains(t, resp.Stdout, "kck-pickup-a-session", "stdout")
	assertContains(t, resp.Stdout, "--here", "stdout")
	assertContains(t, resp.Stdout, "--new-window", "stdout")
	assertContains(t, resp.Stdout, "--no-agent-run", "stdout")
	if strings.Contains(resp.Stdout, "--skill") {
		t.Fatalf("--skill must be removed from help:\n%s", resp.Stdout)
	}
}
```
