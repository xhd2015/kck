# Assert

```go
import "os"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "Would pickup from "+req.SessionID, "stdout")
	assertContains(t, resp.Stdout, "terminal: current", "stdout")
	assertContains(t, resp.Stdout, "draft: read "+draftSkillDisplay(req), "stdout")
	assertContains(t, resp.Stdout, "session-id: "+req.SessionID, "stdout")
	assertContains(t, resp.Stdout, "summarize decisions", "stdout")
	if len(resp.Opened) != 0 {
		t.Fatalf("dry-run must not open: %v", resp.Opened)
	}
	if len(resp.Foreground) != 0 {
		t.Fatalf("dry-run must not run: %v", resp.Foreground)
	}
	if _, err := os.Stat(req.CachePath); err != nil {
		t.Fatalf("dry-run must hydrate skill cache: %v", err)
	}
}
```
