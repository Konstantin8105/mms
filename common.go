package mms

import (
	"fmt"
	"runtime"
)

// Debug is flag for switching to debug mode
var Debug bool = false

func called() string {
	function, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf(
		"\tFile     : %s\n"+
			"\tFunction : %s\n"+
			"\tLine     : %d\n",
		file,
		runtime.FuncForPC(function).Name(),
		line,
	)
}
