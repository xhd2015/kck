## Expected

- Exit 0; valid JSON; no ANSI.
- JSON text includes iterm value `w=42 t=3` and/or `"iterm"` key with that value.

## Exit Code

0

```go
import (
	"encoding/json"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertNoANSI(t, resp.Stdout, "json stdout")
	raw := strings.TrimSpace(resp.Stdout)
	var anyJSON any
	if err := json.Unmarshal([]byte(raw), &anyJSON); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "w=42 t=3") && !strings.Contains(resp.Stdout, `"iterm"`) {
		t.Fatalf("JSON must include iterm field or w=42 t=3; got:\n%s", resp.Stdout)
	}
	// Strong preference: formatted iterm string present.
	if !strings.Contains(resp.Stdout, "w=42 t=3") {
		t.Fatalf("JSON should include formatted iterm w=42 t=3; got:\n%s", resp.Stdout)
	}
}
```
