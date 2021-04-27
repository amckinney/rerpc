// +build tools

package tools

import (
	_ "golang.org/x/tools/cmd/stringer" // for code gen
	_ "google.golang.org/grpc"          // for code gen
)
