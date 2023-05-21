package tests_test

import (
	"context"
	"sync"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/tests/.expect/dal_test/query"
)

var useOnce sync.Once
var ctx = context.Background()

func CRUDInit() {
	query.Use(DB)
	query.SetDefault(DB)
}
