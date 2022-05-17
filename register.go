package xk6outputerror

import (
	"github.com/Tinkoff/xk6-output-error/pkg"
	"go.k6.io/k6/output"
)

func init() {
	output.RegisterExtension(pkg.NameOutput, pkg.New)
}
