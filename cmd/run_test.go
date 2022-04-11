package cmd_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/cmd"
	"github.com/unionj-cloud/go-doudou/cmd/internal/svc"
	"github.com/unionj-cloud/go-doudou/cmd/mock"
	"testing"
)

func Test_runCmd(t *testing.T) {
	Convey("Should not panic when run svc run command", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		s := mock.NewMockISvc(ctrl)
		s.
			EXPECT().
			Run(gomock.Any()).
			AnyTimes().
			Return()

		cmd.RunSvc = func(dir string, opts ...svc.SvcOption) svc.ISvc {
			return s
		}

		So(func() {
			ExecuteCommandC(cmd.GetRootCmd(), []string{"svc", "run"}...)
		}, ShouldNotPanic)
	})
}
