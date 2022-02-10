package ratelimit

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Limit
		wantErr bool
	}{
		{
			name: "",
			args: args{
				value: "0.0055-S-20",
			},
			want: Limit{
				Rate:   0.0055,
				Burst:  20,
				Period: time.Second,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				value: "1000-H",
			},
			want: Limit{
				Rate:   1000,
				Burst:  1,
				Period: time.Hour,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
