## Expected

- Exit 0.
- Data row for `agent-sess-grok-1` has AGENT_RUN `no` (bare grok, not under agent-run).
- AGENT_SID shows `agent-sess-grok-1`.

## Exit Code

0

```go
import (
	"regexp"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertTrailingNewline(t, resp.Stdout, "live stdout")
	assertContains(t, resp.Stdout, "AGENT_RUN", "header")
	assertContains(t, resp.Stdout, "agent-sess-grok-1", "agent sid")
	// Fixed-width: STATE busy, REASON -, AGENT_RUN no, then AGENT_SID.
	re := regexp.MustCompile(`(?m)^.*agent-sess-grok-1.*\bbusy\b\s+-\s+no\s+agent-sess-grok-1`)
	if !re.MatchString(resp.Stdout) {
		// Fallback: single line scan
		var dataLine string
		for _, line := range strings.Split(resp.Stdout, "\n") {
			if strings.Contains(line, "agent-sess-grok-1") && strings.Contains(line, "busy") {
				dataLine = line
				break
			}
		}
		if dataLine == "" {
			t.Fatalf("missing busy data row; stdout:\n%s", resp.Stdout)
		}
		// AGENT_RUN must be no (not yes). After reason "-", next cell is agent_run.
		// Reject if " -                yes" pattern for agent_run.
		if regexp.MustCompile(`busy\s+-\s+yes\s`).MatchString(dataLine) {
			t.Fatalf("AGENT_RUN should be no for bare grok; line=%q", dataLine)
		}
		if !regexp.MustCompile(`busy\s+-\s+no\s`).MatchString(dataLine) {
			t.Fatalf("want AGENT_RUN no after busy/-; line=%q\nstdout:\n%s", dataLine, resp.Stdout)
		}
	}
	assertContains(t, resp.Stdout, "1 sessions · 1 needs attention · 0 sendable", "footer")
}
```
