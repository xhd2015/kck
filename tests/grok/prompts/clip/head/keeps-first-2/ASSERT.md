## Expected

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "p1", "stdout")
	assertContains(t, out, "p2", "stdout")
	assertContains(t, out, "(...3 omitted...)", "stdout")
	for _, drop := range []string{"p3", "p4", "p5"} {
		assertNotContains(t, out, drop, "stdout")
	}
	// marker after kept prompts
	i1 := strings.Index(out, "p1")
	i2 := strings.Index(out, "p2")
	im := strings.Index(out, "(...3 omitted...)")
	if i1 < 0 || i2 < 0 || im < 0 || !(i1 < i2 && i2 < im) {
		t.Fatalf("want p1 then p2 then trailing marker:\n%s", out)
	}
}
```
