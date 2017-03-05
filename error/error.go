package error

import (
    "fmt"
    "os"
)

func DebugError(err interface{}) {
    if err != nil {
        fmt.Println(err)
	    os.Exit(1)
    }
}
