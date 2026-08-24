## Expected

- Exit 0.
- Stdout contains `AGENT_RUN` header and `yes` on the data row for managed grok.
- Stdout contains agent session id `agent-sess-ar-grok`.
- Footer: `1 sessions · 1 needs attention · 0 sendable`.

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
	assertContains(t, resp.Stdout, "agent-sess-ar-grok", "agent sid")
	var dataLine string
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.Contains(line, "agent-sess-ar-grok") && strings.Contains(line, "busy") {
			dataLine = line
			break
		}
	}
	if dataLine == "" {
		t.Fatalf("missing data row with agent-sess-ar-grok; stdout:\n%s", resp.Stdout)
	}
	// Fixed-width: STATE busy, REASON -, AGENT_RUN yes.
	if !strings.Contains(dataLine, "busy") || !regexp.MustCompile(`busy\s+-\s+yes\s`).MatchString(dataLine) {
		t.Fatalf("data row want AGENT_RUN yes after busy/-; line=%q", dataLine)
	}
	assertContains(t, resp.Stdout, "1 sessions · 1 needs attention · 0 sendable", "footer")
}
```
