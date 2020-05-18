package httpsum

import "errors"

var httpStatusError = errors.New("HTTP Status code error")

var timeoutError = errors.New("HTTP Timeout error")
