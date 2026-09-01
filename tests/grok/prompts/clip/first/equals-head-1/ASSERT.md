## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "p1", "stdout")
	assertContains(t, out, "(...4 omitted...)", "stdout")
	for _, drop := range []string{"p2", "p3", "p4", "p5"} {
		assertNotContains(t, out, drop, "stdout")
	}
}
```
