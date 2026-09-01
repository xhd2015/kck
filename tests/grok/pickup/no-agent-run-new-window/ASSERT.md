# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "opened: new window; pickup from session "+req.SessionID, "stdout")
	if len(resp.Opened) != 1 {
		t.Fatalf("Opened=%v, want 1", resp.Opened)
	}
	cmd := resp.Opened[0]
	if i := strings.Index(cmd, "|"); i >= 0 {
		cmd = cmd[i+1:]
	}
	if strings.Contains(cmd, "agent-run") {
		t.Fatalf("no-agent-run must not use agent-run: %q", resp.Opened[0])
	}
	if !strings.Contains(cmd, "grok") {
		t.Fatalf("want bare grok open: %q", resp.Opened[0])
	}
	if len(resp.Staged) != 1 {
		t.Fatalf("Staged=%v, want 1", resp.Staged)
	}
	assertContains(t, resp.Staged[0], "read "+draftSkillDisplay(req), "draft")
	assertContains(t, resp.Staged[0], "session-id: "+req.SessionID, "draft")
	assertContains(t, resp.Staged[0], "continue the refactor", "draft")
}
```
