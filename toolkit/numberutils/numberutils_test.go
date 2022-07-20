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
	numberutils.LargestRemainder(input, 3)
	for _, item := range input {
		fmt.Printf("%v\t%v\t%v\t%s\n", item.Data, item.Value, item.Percent, item.PercentFormatted)
	}
}

func TestLargestRemainderZeroSum(t *testing.T) {
	input := []numberutils.Percentage{
		{
			Value: 1,
			Data:  "a",
		},
		{
			Value: 0,
			Data:  "b",
		},
		{
			Value: 0,
			Data:  "c",
		},
	}
	numberutils.LargestRemainder(input, 3)
	for _, item := range input {
		fmt.Printf("%v\t%v\t%v\t%s\n", item.Data, item.Value, item.Percent, item.PercentFormatted)
	}
}
