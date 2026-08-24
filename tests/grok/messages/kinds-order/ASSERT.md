## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	for _, want := range []string{
		"[user] : u1",
		"[thinking] : th1",
		"[tool] :",
		"tool=shell",
		"[assistant] : a1",
	} {
		assertContains(t, out, want, "stdout")
	}
	iUser := strings.Index(out, "[user] : u1")
	iThink := strings.Index(out, "[thinking] : th1")
	iTool := strings.Index(out, "[tool] :")
	iAsst := strings.Index(out, "[assistant] : a1")
	if !(iUser < iThink && iThink < iTool && iTool < iAsst) {
		t.Fatalf("kinds out of order: user=%d think=%d tool=%d asst=%d\n%s",
			iUser, iThink, iTool, iAsst, out)
	}
}
```
