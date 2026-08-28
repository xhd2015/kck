# Scenario

```go
import (
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"codex", "open", "019f283a-ffff-7fff-ffff-ffffffffff99"}
	return nil
}
```
