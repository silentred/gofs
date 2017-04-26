package store

import (
	"fmt"
)

var (
	errNotAlign          = fmt.Errorf("not align to 8")
	errInvalidNeddleByte = fmt.Errorf("invalid bytes to needle")
)
