package maputils_test

import (
	"bytes"
	"encoding/json"
	"github.com/goccy/go-reflect"
	"testing"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/maputils"
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

func TestMerge(t *testing.T) {
	for _, tuple := range []struct {
		src      string
		dst      string
		expected string
	}{
		{
			src:      `{}`,
			dst:      `{}`,
			expected: `{}`,
		},
		{
			src:      `{"b":2}`,
			dst:      `{"a":1}`,
			expected: `{"a":1,"b":2}`,
		},
		{
			src:      `{"a":0}`,
			dst:      `{"a":1}`,
			expected: `{"a":0}`,
		},
		{
			src:      `{"a":{       "y":2}}`,
			dst:      `{"a":{"x":1       }}`,
			expected: `{"a":{"x":1, "y":2}}`,
		},
		{
			src:      `{"a":{"x":2}}`,
			dst:      `{"a":{"x":1}}`,
			expected: `{"a":{"x":2}}`,
		},
		{
			src:      `{"a":{       "y":7, "z":8}}`,
			dst:      `{"a":{"x":1, "y":2       }}`,
			expected: `{"a":{"x":1, "y":7, "z":8}}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":"xxx"} }, "a":3 }}`,
			expected: `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2]} }, "a":3 }}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":[3]} }, "a":3 }}`,
			expected: `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[3,1,2]} }, "a":3 }}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":[3,1]} }, "a":3 }}`,
			expected: `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[3,1,2]} }, "a":3 }}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,{"c":4,"d":5}]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":[3,1,{"c":2}]} }, "a":3 }}`,
			expected: `{"1":{"2":{"3":{"a":"A","b":3,"n":[3,1,{"c":2},{"c":4,"d":5}]}},"a":3,"b":1}}`,
		},
	} {
		var dst map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.dst), &dst); err != nil {
			t.Error(err)
			continue
		}

		var src map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.src), &src); err != nil {
			t.Error(err)
			continue
		}

		var expected map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.expected), &expected); err != nil {
			t.Error(err)
			continue
		}

		got := maputils.Merge(dst, src)
		assert(t, expected, got)
	}
}

func assert(t *testing.T, expected, got map[string]interface{}) {
	expectedBuf, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
		return
	}
	gotBuf, err := json.Marshal(got)
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(expectedBuf, gotBuf) != 0 {
		t.Errorf("expected %s, got %s", string(expectedBuf), string(gotBuf))
		return
	}
}

func TestMergeOverwriteSlice(t *testing.T) {
	for _, tuple := range []struct {
		src      string
		dst      string
		expected string
	}{
		{
			src:      `{}`,
			dst:      `{}`,
			expected: `{}`,
		},
		{
			src:      `{"b":2}`,
			dst:      `{"a":1}`,
			expected: `{"a":1,"b":2}`,
		},
		{
			src:      `{"a":0}`,
			dst:      `{"a":1}`,
			expected: `{"a":0}`,
		},
		{
			src:      `{"a":{       "y":2}}`,
			dst:      `{"a":{"x":1       }}`,
			expected: `{"a":{"x":1, "y":2}}`,
		},
		{
			src:      `{"a":{"x":2}}`,
			dst:      `{"a":{"x":1}}`,
			expected: `{"a":{"x":2}}`,
		},
		{
			src:      `{"a":{       "y":7, "z":8}}`,
			dst:      `{"a":{"x":1, "y":2       }}`,
			expected: `{"a":{"x":1, "y":7, "z":8}}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":"xxx"} }, "a":3 }}`,
			expected: `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2]} }, "a":3 }}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":[3]} }, "a":3 }}`,
			expected: `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2]} }, "a":3 }}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":[3,1]} }, "a":3 }}`,
			expected: `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2]} }, "a":3 }}`,
		},
		{
			src:      `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,{"c":4,"d":5}]} }        }}`,
			dst:      `{"1": {        "2": { "3": {"a":"A",        "n":[3,1,{"c":2}]} }, "a":3 }}`,
			expected: `{"1":{"2":{"3":{"a":"A","b":3,"n":[1,{"c":4,"d":5}]}},"a":3,"b":1}}`,
		},
	} {
		var dst map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.dst), &dst); err != nil {
			t.Error(err)
			continue
		}

		var src map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.src), &src); err != nil {
			t.Error(err)
			continue
		}

		var expected map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.expected), &expected); err != nil {
			t.Error(err)
			continue
		}

		got := maputils.MergeOverwriteSlice(dst, src)
		assert(t, expected, got)
	}
}
