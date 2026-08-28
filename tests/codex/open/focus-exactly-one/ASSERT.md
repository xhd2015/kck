## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assertContains(t, resp.Stdout, "focused: window 3, tab 2", "stdout")
	if len(resp.Focused) != 1 || resp.Focused[0] != "w2t2p0" {
		t.Fatalf("Focused=%v", resp.Focused)
	}
}
```
