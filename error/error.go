package error

import (
	"os"

    "github.com/davecgh/go-spew/spew"
)

// DebugError performs reporting about error
func DebugError(err interface{}) {
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}
}
