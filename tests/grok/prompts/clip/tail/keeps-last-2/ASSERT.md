## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "(...3 omitted...)", "stdout")
	assertContains(t, out, "p4", "stdout")
	assertContains(t, out, "p5", "stdout")
	for _, drop := range []string{"p1", "p2", "p3"} {
		assertNotContains(t, out, drop, "stdout")
	}
	im := strings.Index(out, "(...3 omitted...)")
	i4 := strings.Index(out, "p4")
	i5 := strings.Index(out, "p5")
	if im < 0 || i4 < 0 || i5 < 0 || !(im < i4 && i4 < i5) {
		t.Fatalf("want leading marker then p4 then p5:\n%s", out)
	}
}
```
