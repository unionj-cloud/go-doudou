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
			Data:  "a",
		},
		{
			Value: 12,
			Data:  "b",
		},
		{
			Value: 7,
			Data:  "c",
		},
	}
	numberutils.LargestRemainder(input, 0)
	for _, item := range input {
		fmt.Printf("%v\t%v\t%v\n", item.Data, item.Value, item.Percent)
	}
}
