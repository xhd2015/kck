## Expected

- `--open` with `--tab` → Error; no send/open.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertFailure(t, resp)
	assertErrorPrefix(t, resp)
	assertContains(t, resp.Stderr+"\n"+resp.ErrText, "--open", "error")
	assertNoSend(t, resp)
	if len(resp.Opened) != 0 {
		t.Fatalf("Opened = %v, want none", resp.Opened)
	}
}
```
