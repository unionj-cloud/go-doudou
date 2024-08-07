package numberutils

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

type Percentage struct {
	Value            int
	Percent          float64
	PercentFormatted string
	Data             interface{}
}

type decimalPercentage struct {
	Percentage *Percentage
	Percent    decimal.Decimal
	Remainder  decimal.Decimal
}

// LargestRemainder https://en.wikipedia.org/wiki/Largest_remainder_method
func LargestRemainder(percentages []Percentage, places int32) {
	if len(percentages) == 0 {
		return
	}
	decimalPercentages := make([]decimalPercentage, len(percentages))
	for i := range percentages {
		decimalPercentages[i] = decimalPercentage{
			Percentage: &percentages[i],
		}
	}
	var sum int
	for _, item := range decimalPercentages {
		sum += item.Percentage.Value
	}
	if sum > 0 {
		for i, item := range decimalPercentages {
			raw := decimal.NewFromFloat(float64(item.Percentage.Value*100) / float64(sum))
			decimalPercentages[i].Percent = raw.RoundFloor(places)
			decimalPercentages[i].Remainder = raw.Sub(decimal.NewFromFloat(percentages[i].Percent))
		}
		var curSum decimal.Decimal
		for _, item := range decimalPercentages {
			curSum = curSum.Add(item.Percent)
		}
		offset := decimal.New(1, -places)
		limit := decimal.NewFromInt(100)
		for curSum.LessThan(limit) {
			sort.Slice(decimalPercentages, func(i, j int) bool {
				return decimalPercentages[j].Remainder.LessThan(decimalPercentages[i].Remainder)
			})
			decimalPercentages[0].Percent = decimalPercentages[0].Percent.Add(offset)
			decimalPercentages[0].Remainder = decimal.Decimal{}
			curSum = curSum.Add(offset)
		}
	}
	for _, item := range decimalPercentages {
		item.Percentage.Percent, _ = item.Percent.Float64()
		item.Percentage.PercentFormatted = fmt.Sprintf("%."+fmt.Sprint(places)+"f%%", item.Percentage.Percent)
	}
}
