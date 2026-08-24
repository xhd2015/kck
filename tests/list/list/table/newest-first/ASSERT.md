## Expected

- Exit 0.
- Both `s-new` and `s-old` appear in stdout.
- First occurrence of `s-new` is before first occurrence of `s-old`.

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "s-new", "stdout")
	assertContains(t, out, "s-old", "stdout")
	iNew := strings.Index(out, "s-new")
	iOld := strings.Index(out, "s-old")
	if iNew < 0 || iOld < 0 || iNew > iOld {
		t.Fatalf("want s-new before s-old; s-new@%d s-old@%d\n%s", iNew, iOld, out)
	}
}
```
