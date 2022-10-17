package codegen

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/astutils"
	"path/filepath"
	"testing"
)

func TestExprStringP(t *testing.T) {
	Convey("Test ExprStringP", t, func() {
		So(func() {
			astutils.BuildStructCollector(filepath.Join(testDir, "vo", "vo2.go"), ExprStringP)
		}, ShouldNotPanic)
		So(func() {
			astutils.BuildStructCollector(filepath.Join(testDir, "vop", "vo3.go"), ExprStringP)
		}, ShouldPanic)
		So(func() {
			astutils.BuildStructCollector(filepath.Join(testDir, "vop", "vo4.go"), ExprStringP)
		}, ShouldPanic)
	})
}
