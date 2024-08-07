package cast

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func ToInt(s string) int {
	v, err := strconv.ParseInt(s, 0, 0)
	if err == nil {
		return int(v)
	}
	return 0
}

func ToIntE(s string) (int, error) {
	v, err := strconv.ParseInt(s, 0, 0)
	if err == nil {
		return int(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to int", s)
}

func ToInt8E(s string) (int8, error) {
	v, err := strconv.ParseInt(s, 0, 8)
	if err == nil {
		return int8(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to int8", s)
}

func ToInt16E(s string) (int16, error) {
	v, err := strconv.ParseInt(s, 0, 16)
	if err == nil {
		return int16(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to int16", s)
}

func ToInt32E(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 0, 32)
	if err == nil {
		return int32(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to int32", s)
}

func ToInt64E(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to int64", s)
}

func ToUintE(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 0, 0)
	if err == nil {
		return uint(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to uint", s)
}

func ToUint8E(s string) (uint8, error) {
	v, err := strconv.ParseUint(s, 0, 8)
	if err == nil {
		return uint8(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to uint8", s)
}

func ToUint16E(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 0, 16)
	if err == nil {
		return uint16(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to uint16", s)
}

func ToUint32E(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 0, 32)
	if err == nil {
		return uint32(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to uint32", s)
}

func ToUint64E(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to uint64", s)
}

func ToFloat32E(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	if err == nil {
		return float32(v), nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to float32", s)
}

func ToFloat64E(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("unable to cast string %#v to float64", s)
}

func ToErrorE(s string) (error, error) {
	return errors.New(s), nil
}

func ToBoolE(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err == nil {
		return b, nil
	}
	return false, fmt.Errorf("unable to cast string %#v to bool", s)
}

func ToBoolOrDefault(s string, d bool) bool {
	result := d
	if eg, err := ToBoolE(s); err == nil {
		result = eg
	}
	return result
}

func ToIntOrDefault(s string, d int) int {
	result := d
	if eg, err := ToIntE(s); err == nil {
		result = eg
	}
	return result
}

func ToInt64OrDefault(s string, d int64) int64 {
	result := d
	if eg, err := ToInt64E(s); err == nil {
		result = eg
	}
	return result
}

func ToUInt32OrDefault(s string, d uint32) uint32 {
	result := d
	if eg, err := ToUint32E(s); err == nil {
		result = eg
	}
	return result
}

func ToRuneSliceE(s string) ([]rune, error) {
	return []rune(s), nil
}

func ToByteSliceE(s string) ([]byte, error) {
	return []byte(s), nil
}

func ToDecimal(s string) decimal.Decimal {
	ret, _ := decimal.NewFromString(s)
	return ret
}

func ToDecimalE(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}
