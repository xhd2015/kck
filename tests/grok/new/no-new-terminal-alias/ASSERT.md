# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if strings.TrimSpace(resp.Stdout) != "" {
		t.Fatalf("--no-new-terminal must be silent; got %q", resp.Stdout)
	}
	if len(resp.Foreground) != 1 {
		t.Fatalf("Foreground=%v, want 1", resp.Foreground)
	}
	if strings.Contains(resp.Foreground[0], "--new-terminal") {
		t.Fatalf("must omit --new-terminal: %q", resp.Foreground[0])
	}
}
```
