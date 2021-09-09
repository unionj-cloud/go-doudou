package cmd

import (
	"testing"
	"time"
)

func TestRootCmd(t *testing.T) {
	go rootCmd.Run(nil, nil)
	time.Sleep(2 * time.Second)
}
