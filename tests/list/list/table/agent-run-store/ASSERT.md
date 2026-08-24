## Expected

- Exit 0.
- Header has AGENT_RUN and AGENT_SID.
- Data row: store-sess-ar, AGENT_RUN yes, AGENT_SID grok-native-99.

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertTrailingNewline(t, resp.Stdout, "list stdout")
	upper := strings.ToUpper(resp.Stdout)
	assertContains(t, upper, "AGENT_RUN", "header")
	assertContains(t, upper, "AGENT_SID", "header")
	assertContains(t, resp.Stdout, "store-sess-ar", "session id")
	assertContains(t, resp.Stdout, "grok-native-99", "runner session id")
	var dataLine string
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.Contains(line, "store-sess-ar") {
			dataLine = line
			break
		}
	}
	if dataLine == "" {
		t.Fatalf("missing data row; stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(dataLine, "yes") {
		t.Fatalf("store row want AGENT_RUN yes; line=%q", dataLine)
	}
	if !strings.Contains(dataLine, "grok-native-99") {
		t.Fatalf("store row want AGENT_SID grok-native-99; line=%q", dataLine)
	}
}
```