package ddhttp

import (
	"fmt"
)

// BizError is used for business error implemented error interface
// StatusCode will be set to http response status code
// ErrCode is used for business error code
// ErrMsg is custom error message
type BizError struct {
	StatusCode int
	ErrCode    int
	ErrMsg     string
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

// NewBizError is factory function for creating an instance of BizError struct
func NewBizError(err error, opts ...BizErrorOption) *BizError {
	bz := &BizError{
		StatusCode: 500,
		ErrMsg:     err.Error(),
	}
	for _, fn := range opts {
		fn(bz)
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
	return b.String()
}
