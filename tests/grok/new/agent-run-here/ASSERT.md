# Assert

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if strings.TrimSpace(resp.Stdout) != "" {
		t.Fatalf("--here must be silent on stdout; got %q", resp.Stdout)
	}
	if strings.TrimSpace(resp.Stderr) != "" {
		t.Fatalf("--here must be silent on stderr; got %q", resp.Stderr)
	}
	if len(resp.NewTerminal) != 0 {
		t.Fatalf("here must not RunNewTerminal: %v", resp.NewTerminal)
	}
	if len(resp.WaitedIDs) != 0 {
		t.Fatalf("here must not WaitSession: %v", resp.WaitedIDs)
	}
	if len(resp.Foreground) != 1 {
		t.Fatalf("Foreground=%v, want 1", resp.Foreground)
	}
	entry := resp.Foreground[0]
	if !strings.Contains(entry, "agent-run") {
		t.Fatalf("want agent-run, got %q", entry)
	}
	if strings.Contains(entry, "--new-terminal") {
		t.Fatalf("here must omit --new-terminal, got %q", entry)
	}
	if !strings.Contains(entry, "--no-submit") {
		t.Fatalf("want --no-submit, got %q", entry)
	}
	if !strings.Contains(entry, "/brainstorm fix flaky auth") {
		t.Fatalf("want brainstorm prompt, got %q", entry)
	}
}
```
