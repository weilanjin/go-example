package go126

import "errors"

type MyError struct {
	msg string
}

func (e MyError) Error() string {
	return e.msg
}

func asType() {
	var err error = MyError{msg: "example error"}
	myError, ok := errors.AsType[MyError](err)
	if ok {
		println("Matched MyError:", myError.Error())
	}
}
