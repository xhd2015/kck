## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	// 1785493077/1785493078s → 2026-07-31 18:17:57 / 18:17:58 in +08:00
	assertContains(t, out, "[2026-07-31 18:17:57] [user] : u2", "stdout")
	assertContains(t, out, "[2026-07-31 18:17:58] [assistant] : a2", "stdout")
	if strings.Contains(out, "[—]") {
		t.Fatalf("fixture has wire times; did not expect [—]:\n%s", out)
	}
}
```
