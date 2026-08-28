## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	got := resp.ErrText + resp.Stderr
	assertContains(t, got, "invalid --cron:", "err")
	if len(resp.SendCalls) != 0 {
		t.Fatalf("SendText must not run on bad cron; got %d", len(resp.SendCalls))
	}
}
```
