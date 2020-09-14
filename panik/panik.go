package panik

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Panikers() {
	if r := recover(); r != nil {
		log.Println("Jangan Panik is wrapping ", Wrap(r.(error)).StackTrace)
	}
}

// Error is the type that implements the error interface.
// It contains the underlying err and its stacktrace.
type Error struct {
	Err        error
	StackTrace string
}

func (m Error) Error() string {
	return m.Err.Error() + m.StackTrace
}

// Wrap annotates the given error with a stack trace
func Wrap(err error) Error {
	return Error{Err: err, StackTrace: getStackTrace()}
}

func getStackTrace() string {
	stackBuf := make([]uintptr, 5000)
	length := runtime.Callers(3, stackBuf[:])
	stack := stackBuf[:length]

	trace := ""
	frames := runtime.CallersFrames(stack)
	i := 0
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			if i > 0 {
				trace = trace + fmt.Sprintf("\n\tFile: %s, Line: %d. Function: %s", frame.File, frame.Line, frame.Function)
			}
		}
		if !more {
			break
		}
		i++
	}
	return trace
}
