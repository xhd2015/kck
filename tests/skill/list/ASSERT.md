## Expected Output

```
---
version: 3
---
kck
info
list
messages
open
overview
pickup
resolve
send
snapshot
status
kck-pickup-a-session
```

## Exit Code

- 0

```go
import (
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	assert.Output(t, resp.Stdout, `---
version: 3
---
kck
info
list
messages
open
overview
pickup
resolve
send
snapshot
status
kck-pickup-a-session
`)
}
```
