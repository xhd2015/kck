# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.NewTerminal) != 1 {
		t.Fatalf("NewTerminal=%v, want 1", resp.NewTerminal)
	}
	entry := resp.NewTerminal[0]
	if strings.Contains(entry, "--no-submit") {
		t.Fatalf("--submit must omit --no-submit: %q", entry)
	}
	if !strings.Contains(entry, "--open") {
		t.Fatalf("want --open, got %q", entry)
	}
}
```
