package esutils

import (
	"context"
	"testing"
)

func TestNewIndex(t *testing.T) {
	es := setupSubTest("test_newindex")
	es.esIndex = "notexists"

	type args struct {
		mapping string
	}
	mapping := NewMapping(MappingPayload{
		Base{
			Index: es.esIndex,
			Type:  es.esType,
		},
		[]Field{
			{
				Name: "createAt",
				Type: DATE,
			},
		},
	})
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantErr    bool
	}{
		{
			name: "1",
			args: args{
				mapping: mapping,
			},
			wantExists: false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := es.NewIndex(context.Background(), tt.args.mapping)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExists != tt.wantExists {
				t.Errorf("NewIndex() = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}
