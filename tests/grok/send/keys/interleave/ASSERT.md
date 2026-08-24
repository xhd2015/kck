## Expected

- Payload Up + pick + Enter + positional tail; NoSubmit true.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.SendCalls) != 1 {
		t.Fatalf("SendCalls = %#v", resp.SendCalls)
	}
	c := resp.SendCalls[0]
	want := "\x1b[Apick\ntail"
	if c.Text != want {
		t.Fatalf("Text = %q want %q", c.Text, want)
	}
	if !c.Opts.NoSubmit {
		t.Fatalf("opts = %+v", c.Opts)
	}
}
```
