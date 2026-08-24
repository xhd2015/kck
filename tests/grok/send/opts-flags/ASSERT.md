## Expected

- Focus / NoSubmit / NoCtrlU plumbed.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if len(resp.SendCalls) != 1 {
		t.Fatalf("SendCalls = %#v", resp.SendCalls)
	}
	o := resp.SendCalls[0].Opts
	if !o.Focus || !o.NoSubmit || !o.NoCtrlU {
		t.Fatalf("opts = %+v", o)
	}
}
```
