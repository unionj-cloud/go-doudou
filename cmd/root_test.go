package cmd_test

import (
	"testing"
	"time"

	"github.com/unionj-cloud/go-doudou/v2/cmd"
)

func TestRootCmd(t *testing.T) {
	go cmd.GetRootCmd().Run(nil, nil)
	time.Sleep(2 * time.Second)
}
