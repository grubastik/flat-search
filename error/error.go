package error

import (
	"fmt"
	"os"
)

// DebugError performs reporting about error
func DebugError(err interface{}) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
