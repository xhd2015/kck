# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "opened: here; pickup from session "+req.SessionID, "stdout")
	if len(resp.Opened) != 0 {
		t.Fatalf("here must not open window: %v", resp.Opened)
	}
	if len(resp.Foreground) != 1 {
		t.Fatalf("Foreground=%v, want 1", resp.Foreground)
	}
	entry := resp.Foreground[0]
	if !strings.Contains(entry, "agent-run") || !strings.Contains(entry, "codex-tty") {
		t.Fatalf("want agent-run codex-tty here, got %q", entry)
	}
	if !strings.Contains(entry, "--no-submit") {
		t.Fatalf("want --no-submit, got %q", entry)
	}
	if !strings.Contains(entry, "extract TODOs") {
		t.Fatalf("draft missing msg in %q", entry)
	}
}
```
