package cast

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func ToIntSliceE(s []string) ([]int, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []int", s)
	}
	var ret []int
	for _, item := range s {
		i, err := ToIntE(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []int because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToInt8SliceE(s []string) ([]int8, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []int8", s)
	}
	var ret []int8
	for _, item := range s {
		i, err := ToInt8E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []int8 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToInt16SliceE(s []string) ([]int16, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []int16", s)
	}
	var ret []int16
	for _, item := range s {
		i, err := ToInt16E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []int16 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToInt32SliceE(s []string) ([]int32, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []int32", s)
	}
	var ret []int32
	for _, item := range s {
		i, err := ToInt32E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []int32 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToInt64SliceE(s []string) ([]int64, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []int64", s)
	}
	var ret []int64
	for _, item := range s {
		i, err := ToInt64E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []int64 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToUintSliceE(s []string) ([]uint, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []uint", s)
	}
	var ret []uint
	for _, item := range s {
		i, err := ToUintE(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []uint because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToUint8SliceE(s []string) ([]uint8, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []uint8", s)
	}
	var ret []uint8
	for _, item := range s {
		i, err := ToUint8E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []uint8 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToUint16SliceE(s []string) ([]uint16, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []uint16", s)
	}
	var ret []uint16
	for _, item := range s {
		i, err := ToUint16E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []uint16 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToUint32SliceE(s []string) ([]uint32, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []uint32", s)
	}
	var ret []uint32
	for _, item := range s {
		i, err := ToUint32E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []uint32 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToUint64SliceE(s []string) ([]uint64, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []uint64", s)
	}
	var ret []uint64
	for _, item := range s {
		i, err := ToUint64E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []uint64 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToFloat32SliceE(s []string) ([]float32, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []float32", s)
	}
	var ret []float32
	for _, item := range s {
		i, err := ToFloat32E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []float32 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToFloat64SliceE(s []string) ([]float64, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []float64", s)
	}
	var ret []float64
	for _, item := range s {
		i, err := ToFloat64E(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []float64 because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func ToErrorSliceE(s []string) ([]error, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []error", s)
	}
	var ret []error
	for _, item := range s {
		i, _ := ToErrorE(item)
		ret = append(ret, i)
	}
	return ret, nil
}

func ToBoolSliceE(s []string) ([]bool, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []bool", s)
	}
	var ret []bool
	for _, item := range s {
		i, err := ToBoolE(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []bool because of error %s", s, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

//func ToComplex64SliceE(s []string) ([]complex64, error) {
//	if s == nil {
//		return nil, fmt.Errorf("unable to cast string slice %#v to []complex64", s)
//	}
//	var ret []complex64
//	for _, item := range s {
//		i, err := ToComplex64E(item)
//		if err != nil {
//			return nil, fmt.Errorf("unable to cast string slice %#v to []complex64 because of error %s", s, err)
//		}
//		ret = append(ret, i)
//	}
//	return ret, nil
//}
//
//func ToComplex128SliceE(s []string) ([]complex128, error) {
//	if s == nil {
//		return nil, fmt.Errorf("unable to cast string slice %#v to []complex128", s)
//	}
//	var ret []complex128
//	for _, item := range s {
//		i, err := ToComplex128E(item)
//		if err != nil {
//			return nil, fmt.Errorf("unable to cast string slice %#v to []complex128 because of error %s", s, err)
//		}
//		ret = append(ret, i)
//	}
//	return ret, nil
//}

func ToRuneSliceSliceE(s []string) ([][]rune, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to [][]rune", s)
	}
	var ret [][]rune
	for _, item := range s {
		i, _ := ToRuneSliceE(item)
		ret = append(ret, i)
	}
	return ret, nil
}

func ToByteSliceSliceE(s []string) ([][]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to [][]byte", s)
	}
	var ret [][]byte
	for _, item := range s {
		i, _ := ToByteSliceE(item)
		ret = append(ret, i)
	}
	return ret, nil
}

func ToInterfaceSliceE(s []string) ([]interface{}, error) {
	if s == nil {
		return nil, fmt.Errorf("unable to cast string slice %#v to []interface{}", s)
	}
	var ret []interface{}
	for _, item := range s {
		ret = append(ret, item)
	}
	return ret, nil
}

func ToDecimalSlice(s []string) []decimal.Decimal {
	var ret []decimal.Decimal
	for _, item := range s {
		d, _ := decimal.NewFromString(item)
		ret = append(ret, d)
	}
	return ret
}

func ToDecimalSliceE(s []string) ([]decimal.Decimal, error) {
	var ret []decimal.Decimal
	for _, item := range s {
		d, err := ToDecimalE(item)
		if err != nil {
			return nil, fmt.Errorf("unable to cast string slice %#v to []decimal.Decimal because of error %s", s, err)
		}
		ret = append(ret, d)
	}
	return ret, nil
}
