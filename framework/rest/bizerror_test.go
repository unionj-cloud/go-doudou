package rest_test

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
)

func TestWithStatusCode(t *testing.T) {
	Convey("Create a BizError with 401 http status code", t, func() {
		bizError := rest.NewBizError(errors.New("Unauthorised"), rest.WithStatusCode(401))
		So(bizError, ShouldNotBeZeroValue)
		So(bizError.StatusCode, ShouldEqual, 401)

		Convey("Should output Unauthorised", func() {
			So(bizError.Error(), ShouldEqual, "Unauthorised")
		})
	})
}

func TestWithErrCode(t *testing.T) {
	Convey("Create a BizError with 100401 business error code", t, func() {
		bizError := rest.NewBizError(errors.New("Unauthorised"), rest.WithErrCode(100401))
		So(bizError, ShouldNotBeZeroValue)
		So(bizError.ErrCode, ShouldEqual, 100401)

		Convey("Should output 100401 Unauthorised", func() {
			So(bizError.String(), ShouldEqual, "100401 Unauthorised")
		})
	})
}

func TestWithCause(t *testing.T) {
	Convey("Create a BizError with a cause error", t, func() {
		cause := errors.New("root cause")
		bizError := rest.NewBizError(errors.New("Wrapped error"), rest.WithCause(cause))
		So(bizError, ShouldNotBeZeroValue)
		So(bizError.Cause, ShouldEqual, cause)
	})
}

func TestBizErrorString(t *testing.T) {
	Convey("Test BizError String method", t, func() {
		Convey("With ErrCode > 0", func() {
			bizError := rest.NewBizError(errors.New("Error message"), rest.WithErrCode(100))
			So(bizError.String(), ShouldEqual, "100 Error message")
		})

		Convey("With ErrCode = 0", func() {
			bizError := rest.NewBizError(errors.New("Error message"), rest.WithErrCode(0))
			So(bizError.String(), ShouldEqual, "Error message")
		})
	})
}

func TestBizErrorError(t *testing.T) {
	Convey("Test BizError Error method", t, func() {
		bizError := rest.NewBizError(errors.New("Error message"))
		So(bizError.Error(), ShouldEqual, "Error message")
	})
}

func TestUnWrap(t *testing.T) {
	Convey("Test UnWrap function", t, func() {
		Convey("With nested BizErrors", func() {
			rootErr := errors.New("root error")
			level1 := rest.NewBizError(errors.New("level 1"), rest.WithCause(rootErr))
			level2 := rest.NewBizError(errors.New("level 2"), rest.WithCause(level1))

			unwrapped := rest.UnWrap(level2)
			So(unwrapped, ShouldEqual, rootErr)
		})

		Convey("With non-BizError", func() {
			err := errors.New("regular error")
			unwrapped := rest.UnWrap(err)
			So(unwrapped, ShouldEqual, err)
		})

		Convey("With BizError without cause", func() {
			err := rest.NewBizError(errors.New("error without cause"))
			unwrapped := rest.UnWrap(err)
			So(unwrapped.Error(), ShouldEqual, "error without cause")
		})
	})
}

func TestPanicBadRequestErr(t *testing.T) {
	Convey("Test PanicBadRequestErr function", t, func() {
		Convey("With nil error", func() {
			// 不应该panic
			So(func() { rest.PanicBadRequestErr(nil) }, ShouldNotPanic)
		})

		Convey("With non-nil error", func() {
			err := errors.New("bad request")
			So(func() { rest.PanicBadRequestErr(err) }, ShouldPanic)
		})
	})
}

func TestPanicInternalServerError(t *testing.T) {
	Convey("Test PanicInternalServerError function", t, func() {
		Convey("With nil error", func() {
			// 不应该panic
			So(func() { rest.PanicInternalServerError(nil) }, ShouldNotPanic)
		})

		Convey("With non-nil error", func() {
			err := errors.New("internal server error")
			So(func() { rest.PanicInternalServerError(err) }, ShouldPanic)
		})
	})
}

func TestHandleBadRequestErr(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		bizErr, ok := r.(rest.BizError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, bizErr.StatusCode)
		assert.Equal(t, "bad request", bizErr.ErrMsg)
	}()

	rest.HandleBadRequestErr(errors.New("bad request"))
}

func TestHandleInternalServerError(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		bizErr, ok := r.(rest.BizError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, bizErr.StatusCode)
		assert.Equal(t, "internal server error", bizErr.ErrMsg)
	}()

	rest.HandleInternalServerError(errors.New("internal server error"))
}
