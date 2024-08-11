package cast

import (
	"github.com/goccy/go-reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestToInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want: 21,
		},
		{
			name: "",
			args: args{
				s: "not_int",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.args.s); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToIntE(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: "003",
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToIntE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToIntE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToIntE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt8E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int8
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt8E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt8E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt8E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt16E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int16
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt16E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt16E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt16E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt32E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt32E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt32E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt32E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt64E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUintE(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUintE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUintE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUintE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint8E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint8
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint8E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint8E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUint8E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint16E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint16E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint16E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUint16E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint32E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint32E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint32E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUint32E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUint64E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint64E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint64E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUint64E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat32E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    float32
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21.7812",
			},
			want:    21.7812,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat32E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat32E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat32E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64E(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "21",
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64E(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64E() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat64E() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToErrorE(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 error
	}{
		{
			name: "",
			args: args{
				s: "error",
			},
			want:  "error",
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ToErrorE(tt.args.s)
			if got.Error() != tt.want {
				t.Errorf("ToErrorE() got = %v, want %v", got.Error(), tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToErrorE() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestToBoolE(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "true",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				s: "21a",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				s: "",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBoolE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBoolE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToBoolE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestToComplex64E(t *testing.T) {
//	type args struct {
//		s string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    complex64
//		wantErr bool
//	}{
//		{
//			name: "",
//			args: args{
//				s: "21",
//			},
//			want:    21,
//			wantErr: false,
//		},
//		{
//			name: "",
//			args: args{
//				s: "21a",
//			},
//			want:    0,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ToComplex64E(tt.args.s)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ToComplex64E() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("ToComplex64E() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestToComplex128E(t *testing.T) {
//	type args struct {
//		s string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    complex128
//		wantErr bool
//	}{
//		{
//			name: "",
//			args: args{
//				s: "21",
//			},
//			want:    21,
//			wantErr: false,
//		},
//		{
//			name: "",
//			args: args{
//				s: "21a",
//			},
//			want:    0,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ToComplex128E(tt.args.s)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ToComplex128E() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("ToComplex128E() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestToRuneSliceE(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "text",
			},
			want:    []rune("text"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToRuneSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToRuneSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRuneSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToByteSliceE(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "",
			args: args{
				s: "text",
			},
			want:    []byte("text"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToByteSliceE(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToByteSliceE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToByteSliceE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBoolOrDefault(t *testing.T) {
	type args struct {
		s string
		d bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				s: "true",
				d: false,
			},
			want: true,
		},
		{
			name: "",
			args: args{
				s: "21a",
				d: true,
			},
			want: true,
		},
		{
			name: "",
			args: args{
				s: "",
				d: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBoolOrDefault(tt.args.s, tt.args.d); got != tt.want {
				t.Errorf("ToBoolOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToIntOrDefault(t *testing.T) {
	type args struct {
		s string
		d int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				s: "not_int",
				d: 10,
			},
			want: 10,
		},
		{
			name: "",
			args: args{
				s: "5",
				d: 10,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToIntOrDefault(tt.args.s, tt.args.d); got != tt.want {
				t.Errorf("ToIntOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToDecimal(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "",
			args: args{
				s: "2.43",
			},
			want: decimal.NewFromFloat(2.43),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToDecimal(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}
