package cmd_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/cmd"
	"testing"
)

func Test_svcCmd(t *testing.T) {
	Convey("Should not panic when run svc command", t, func() {
		So(func() {
			ExecuteCommandC(cmd.GetRootCmd(), []string{"svc"}...)
		}, ShouldNotPanic)
	})
}
