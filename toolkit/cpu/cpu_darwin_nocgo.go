//go:build darwin && !cgo
// +build darwin,!cgo

package cpu

import "github.com/unionj-cloud/go-doudou/toolkit/internal/common"

func perCPUTimes() ([]TimesStat, error) {
	return []TimesStat{}, common.ErrNotImplementedError
}

func allCPUTimes() ([]TimesStat, error) {
	return []TimesStat{}, common.ErrNotImplementedError
}
