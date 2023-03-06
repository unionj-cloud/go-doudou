package rest

import (
	"fmt"
	"net/http"
)

// BizError is used for business error implemented error interface
// StatusCode will be set to http response status code
// ErrCode is used for business error code
// ErrMsg is custom error message
type BizError struct {
	StatusCode int
	ErrCode    int
	ErrMsg     string
	Cause      error
}

type BizErrorOption func(bizError *BizError)

func WithStatusCode(statusCode int) BizErrorOption {
	return func(bizError *BizError) {
		bizError.StatusCode = statusCode
	}
}

func WithErrCode(errCode int) BizErrorOption {
	return func(bizError *BizError) {
		bizError.ErrCode = errCode
	}
}

func WithCause(cause error) BizErrorOption {
	return func(bizError *BizError) {
		bizError.Cause = cause
	}
}

// NewBizError is factory function for creating an instance of BizError struct
func NewBizError(err error, opts ...BizErrorOption) BizError {
	bz := BizError{
		ErrCode:    1,
		StatusCode: http.StatusInternalServerError,
		ErrMsg:     err.Error(),
	}
	for _, fn := range opts {
		fn(&bz)
	}
	return bz
}

// String function is used for printing string representation of a BizError instance
func (b BizError) String() string {
	if b.ErrCode > 0 {
		return fmt.Sprintf("%d %s", b.ErrCode, b.ErrMsg)
	}
	return b.ErrMsg
}

// Error is used for implementing error interface
func (b BizError) Error() string {
	return b.ErrMsg
}

func HandleBadRequestErr(err error) {
	panic(NewBizError(err, WithStatusCode(http.StatusBadRequest)))
}

func HandleInternalServerError(err error) {
	panic(NewBizError(err))
}
