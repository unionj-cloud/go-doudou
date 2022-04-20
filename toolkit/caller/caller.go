package caller

import (
	"fmt"
	"runtime"
)

type Caller struct {
	Name string
	File string
	Line int
}

func (c Caller) String() string {
	return fmt.Sprintf("called from %s on %s#%d", c.Name, c.File, c.Line)
}

func NewCaller() Caller {
	var caller Caller
	pc, file, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		caller.File = file
		caller.Line = line
		caller.Name = details.Name()
	}
	return caller
}
