package registry

import (
	"reflect"
	"testing"
)

func Test_seeds(t *testing.T) {
	type args struct {
		seedstr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				seedstr: "seed-01,seed-02,seed-03",
			},
			want: []string{"seed-01:56199", "seed-02:56199", "seed-03:56199"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seeds(tt.args.seedstr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seeds() = %v, want %v", got, tt.want)
			}
		})
	}
}
