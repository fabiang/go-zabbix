package zabbix

import (
	"fmt"
	"os"
	"strconv"
)

// debug caches the value of environment variable ZBX_DEBUG from program start.
var debug bool

func init() {
	debug, _ = strconv.ParseBool(os.Getenv("ZBX_DEBUG"))
}

// dprintf prints formatted debug message to STDERR if the ZBX_DEBUG environment
// variable is set to "1".
func dprintf(format string, a ...interface{}) {
	if debug {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}
