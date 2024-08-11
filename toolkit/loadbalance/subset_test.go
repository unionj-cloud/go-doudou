package loadbalance

import (
	"github.com/goccy/go-reflect"
	"testing"
)

func TestSubset(t *testing.T) {
	type args struct {
		backends   []string
		clientId   int
		subsetSize int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				backends: []string{
					"192.168.0.0",
					"192.168.0.1",
					"192.168.0.2",
					"192.168.0.3",
					"192.168.0.4",
					"192.168.0.5",
					"192.168.0.6",
					"192.168.0.7",
					"192.168.0.8",
					"192.168.0.9",
					"192.168.0.10",
					"192.168.0.11",
				},
				clientId:   8,
				subsetSize: 3,
			},
			want: []string{"192.168.0.5", "192.168.0.7", "192.168.0.3"},
		},
		{
			name: "",
			args: args{
				backends: []string{
					"192.168.0.0",
					"192.168.0.1",
					"192.168.0.2",
					"192.168.0.3",
					"192.168.0.4",
					"192.168.0.5",
					"192.168.0.6",
					"192.168.0.7",
					"192.168.0.8",
					"192.168.0.9",
					"192.168.0.10",
					"192.168.0.11",
				},
				clientId:   9,
				subsetSize: 3,
			},
			want: []string{"192.168.0.9", "192.168.0.8", "192.168.0.10"},
		},
		{
			name: "",
			args: args{
				backends: []string{
					"192.168.0.0",
					"192.168.0.1",
					"192.168.0.2",
					"192.168.0.3",
					"192.168.0.4",
					"192.168.0.5",
					"192.168.0.6",
					"192.168.0.7",
					"192.168.0.8",
					"192.168.0.9",
					"192.168.0.10",
					"192.168.0.11",
				},
				clientId:   10,
				subsetSize: 3,
			},
			want: []string{"192.168.0.6", "192.168.0.4", "192.168.0.1"},
		},
		{
			name: "",
			args: args{
				backends: []string{
					"192.168.0.0",
					"192.168.0.1",
					"192.168.0.2",
					"192.168.0.3",
					"192.168.0.4",
					"192.168.0.5",
					"192.168.0.6",
					"192.168.0.7",
					"192.168.0.8",
					"192.168.0.9",
					"192.168.0.10",
					"192.168.0.11",
				},
				clientId:   11,
				subsetSize: 3,
			},
			want: []string{"192.168.0.0", "192.168.0.11", "192.168.0.2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Subset(tt.args.backends, tt.args.clientId, tt.args.subsetSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subset() = %v, want %v", got, tt.want)
			}
		})
	}
}
