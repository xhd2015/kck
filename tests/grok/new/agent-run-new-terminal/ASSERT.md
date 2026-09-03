# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	arID := expectedSessionID("fix flaky auth")
	providerID := "01a064d2-70ec-7162-b36b-8a50ba323569"
	assertContains(t, resp.Stdout, "opened: new terminal; new grok session", "stdout")
	assertContains(t, resp.Stdout, "session-id: "+providerID, "stdout")
	assertNotContains(t, resp.Stdout, "session-id: "+arID, "stdout")
	if len(resp.Foreground) != 0 {
		t.Fatalf("new-terminal must not RunForeground: %v", resp.Foreground)
	}
	if len(resp.NewTerminal) != 1 {
		t.Fatalf("NewTerminal=%v, want 1", resp.NewTerminal)
	}
	entry := resp.NewTerminal[0]
	if !strings.Contains(entry, "agent-run") {
		t.Fatalf("want agent-run, got %q", entry)
	}
	if !strings.Contains(entry, "--new-terminal") {
		t.Fatalf("want --new-terminal, got %q", entry)
	}
	if !strings.Contains(entry, "--no-submit") {
		t.Fatalf("want --no-submit, got %q", entry)
	}
	if !strings.Contains(entry, "grok-tty") {
		t.Fatalf("want grok-tty, got %q", entry)
	}
	if !strings.Contains(entry, "/brainstorm fix flaky auth") {
		t.Fatalf("want brainstorm prompt, got %q", entry)
	}
	if !strings.Contains(entry, arID) {
		t.Fatalf("want agent-run session id %s in argv %q", arID, entry)
	}
	if len(resp.WaitedIDs) != 1 || resp.WaitedIDs[0] != arID {
		t.Fatalf("WaitedIDs=%v want [%s]", resp.WaitedIDs, arID)
	}
}
```
