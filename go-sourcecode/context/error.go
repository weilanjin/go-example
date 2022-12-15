package context

import "errors"

var Canceled = errors.New("context canceled")

var DeadlineExceeded error = deadlineExceededError{}

// 实现了 error 接口
type deadlineExceededError struct{}

func (deadlineExceededError) Error() string { return "context deadline exceeded" }