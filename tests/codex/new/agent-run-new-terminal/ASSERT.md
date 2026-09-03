# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	providerID := "019f283a-aaaa-7bbb-cccc-dddddddddddd"
	assertContains(t, resp.Stdout, "opened: new terminal; new codex session", "stdout")
	assertContains(t, resp.Stdout, "session-id: "+providerID, "stdout")
	if len(resp.NewTerminal) != 1 {
		t.Fatalf("NewTerminal=%v, want 1", resp.NewTerminal)
	}
	entry := resp.NewTerminal[0]
	if !strings.Contains(entry, "codex-tty") {
		t.Fatalf("want codex-tty, got %q", entry)
	}
	if !strings.Contains(entry, "$brainstorm extract TODOs") {
		t.Fatalf("want $brainstorm prompt, got %q", entry)
	}
	if !strings.Contains(entry, "--new-terminal") {
		t.Fatalf("want --new-terminal, got %q", entry)
	}
}
```
