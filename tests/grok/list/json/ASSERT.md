## Expected

- JSON includes session_id.

```go
import (
	"encoding/json"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if strings.Contains(resp.Stdout, "\x1b[") {
		t.Fatalf("json must not contain ANSI")
	}
	var env struct {
		Sessions []struct {
			SessionID string `json:"session_id"`
		} `json:"sessions"`
		Summary struct {
			Count int `json:"count"`
		} `json:"summary"`
	}
	if err := json.Unmarshal([]byte(resp.Stdout), &env); err != nil {
		t.Fatalf("json: %v\n%s", err, resp.Stdout)
	}
	if env.Summary.Count != 1 || len(env.Sessions) != 1 || env.Sessions[0].SessionID != fixtureKckListSID {
		t.Fatalf("unexpected json: %+v", env)
	}
}
```
