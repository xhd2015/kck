## Expected

- Success; SendText `"\x03"` with NoCtrlU+NoSubmit.

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
	if c.Text != "\x03" {
		t.Fatalf("Text = %q", c.Text)
	}
	if !c.Opts.NoCtrlU || !c.Opts.NoSubmit {
		t.Fatalf("opts = %+v", c.Opts)
	}
}
```
