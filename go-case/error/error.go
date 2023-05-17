package er

import (
	"fmt"
	"runtime/debug"
)

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]any // miscellaneous（各种各样） information
}

func WrapError(err error, messagef string, msgArgs ...any) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]any),
	}
}

func (err MyError) Error() string {
	return err.Message
}
