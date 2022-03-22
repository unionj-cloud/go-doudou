package maputils_test

import (
	"github.com/unionj-cloud/go-doudou/toolkit/maputils"
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	type args struct {
		new map[string]interface{}
		old map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]maputils.Change
	}{
		{
			name: "",
			args: args{
				new: map[string]interface{}{
					"gdd.port":         6060,
					"gdd.read.timeout": "30s",
				},
				old: map[string]interface{}{
					"gdd.port":         3000,
					"gdd.read.timeout": "30s",
				},
			},
			want: map[string]maputils.Change{
				"gdd.port": {
					OldValue:   3000,
					NewValue:   6060,
					ChangeType: maputils.MODIFIED,
				},
			},
		},
		{
			name: "",
			args: args{
				new: map[string]interface{}{
					"gdd.port":         6060,
					"gdd.read.timeout": "30s",
				},
				old: map[string]interface{}{
					"gdd.port": 6060,
				},
			},
			want: map[string]maputils.Change{
				"gdd.read.timeout": {
					OldValue:   nil,
					NewValue:   "30s",
					ChangeType: maputils.ADDED,
				},
			},
		},
		{
			name: "",
			args: args{
				new: map[string]interface{}{
					"gdd.port": 6060,
				},
				old: map[string]interface{}{
					"gdd.port":         6060,
					"gdd.read.timeout": "30s",
				},
			},
			want: map[string]maputils.Change{
				"gdd.read.timeout": {
					OldValue:   "30s",
					NewValue:   nil,
					ChangeType: maputils.DELETED,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maputils.Diff(tt.args.new, tt.args.old); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
