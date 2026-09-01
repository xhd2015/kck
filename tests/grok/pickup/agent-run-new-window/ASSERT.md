# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "opened: new window; pickup from session "+req.SessionID, "stdout")
	if len(resp.Foreground) != 0 {
		t.Fatalf("new-window must not RunForeground: %v", resp.Foreground)
	}
	if len(resp.Opened) != 1 {
		t.Fatalf("Opened=%v, want 1", resp.Opened)
	}
	entry := resp.Opened[0]
	if !strings.Contains(entry, "agent-run") {
		t.Fatalf("want agent-run open, got %q", entry)
	}
	if !strings.Contains(entry, "--no-submit") {
		t.Fatalf("want --no-submit, got %q", entry)
	}
	if !strings.Contains(entry, draftSkillDisplay(req)) {
		t.Fatalf("draft missing tilde skill path in %q", entry)
	}
}
```
