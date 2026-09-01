## Expected

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "main-prompt", "stdout")
	assertContains(t, out, "019f283a-prmt-7ccc-cccc-cccccccccc11", "stdout")
	assertNotContains(t, out, "sub-prompt", "stdout")
	assertNotContains(t, out, "019f283a-prmt-7ccc-cccc-cccccccccc22", "stdout")
}
```
