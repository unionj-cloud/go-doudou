package numberutils_test

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/toolkit/numberutils"
	"testing"
)

func TestLargestRemainder(t *testing.T) {
	input := []numberutils.Percentage{
		{
			Value: 10,
		},
		{
			Value: 12,
		},
		{
			Value: 7,
		},
	}
	numberutils.LargestRemainder(input, 0)
	for _, item := range input {
		fmt.Printf("%v\t%v\n", item.Value, item.Percent)
	}
}
