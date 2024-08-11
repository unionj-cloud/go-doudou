package cast

import (
	"fmt"
	"github.com/goccy/go-reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestToIntSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []int{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToIntSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToIntSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToIntSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt8SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int8
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []int8{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt8SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt8SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInt8SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt16SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int16
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []int16{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt16SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt16SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInt16SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt32SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int32
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []int32{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt32SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt32SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInt32SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []int64{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInt64SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUintSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []uint
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []uint{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUintSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUintSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUintSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint8SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []uint8
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []uint8{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint8SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint8SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUint8SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint16SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []uint16
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []uint16{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint16SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint16SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUint16SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint32SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []uint32
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []uint32{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint32SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint32SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUint32SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint64SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3", "4"},
			},
			want:    []uint64{2, 3, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint64SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint64SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUint64SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat32SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []float32
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3.1314926", "4"},
			},
			want:    []float32{2, 3.1314926, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat32SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat32SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFloat32SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64SliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []float64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"2", "3.1314926", "4"},
			},
			want:    []float64{2, 3.1314926, 4},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"2", "3q", "4"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64SliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64SliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToFloat64SliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleToErrorSliceE() {
	got, _ := ToErrorSliceE([]string{"test1", "test2"})
	fmt.Println(got)

	//Output:
	//[test1 test2]
}

func TestToBoolSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []bool
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"true", "false"},
			},
			want:    []bool{true, false},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{"T", "fff"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBoolSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBoolSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBoolSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestToComplex64SliceE(t *testing.T) {
//	type args struct {
//		s []string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []complex64
//		wantErr bool
//	}{
//		{
//			name: "",
//			args: args{
//				s: []string{"2", "3.1314926", "4"},
//			},
//			want:    []complex64{2, 3.1314926, 4},
//			wantErr: false,
//		},
//		{
//			name: "",
//			args: args{
//				s: []string{"2", "3q", "4"},
//			},
//			want:    nil,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ToComplex64SliceE(tt.args.s)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ToComplex64SliceE() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ToComplex64SliceE() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestToComplex128SliceE(t *testing.T) {
//	type args struct {
//		s []string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []complex128
//		wantErr bool
//	}{
//		{
//			name: "",
//			args: args{
//				s: []string{"2", "3.1314926", "4"},
//			},
//			want:    []complex128{2, 3.1314926, 4},
//			wantErr: false,
//		},
//		{
//			name: "",
//			args: args{
//				s: []string{"2", "3q", "4"},
//			},
//			want:    nil,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ToComplex128SliceE(tt.args.s)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ToComplex128SliceE() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ToComplex128SliceE() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestToRuneSliceSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]rune
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"test1"},
			},
			want:    [][]rune{[]rune("test1")},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToRuneSliceSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToRuneSliceSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRuneSliceSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToByteSliceSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"test1"},
			},
			want:    [][]byte{[]byte("test1")},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToByteSliceSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToByteSliceSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToByteSliceSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInterfaceSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{"test1"},
			},
			want:    []interface{}{"test1"},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInterfaceSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInterfaceSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInterfaceSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToErrorSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []error
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToErrorSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToErrorSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToErrorSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDecimalSlice(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []decimal.Decimal
	}{
		{
			name: "",
			args: args{
				s: []string{
					"2.43",
					"17.89",
				},
			},
			want: []decimal.Decimal{
				decimal.NewFromFloat(2.43),
				decimal.NewFromFloat(17.89),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToDecimalSlice(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDecimalSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDecimalSliceE(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		want    []decimal.Decimal
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: []string{
					"2.43",
					"17.89",
				},
			},
			want: []decimal.Decimal{
				decimal.NewFromFloat(2.43),
				decimal.NewFromFloat(17.89),
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: []string{
					"2.43",
					"17d.89",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDecimalSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDecimalSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ToDecimalSliceE() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
