package esutils

import (
	"context"
	"testing"
)

func TestNewIndex(t *testing.T) {
	es := setupSubTest("test_newindex")

	type args struct {
		mapping MappingPayload
	}
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantErr    bool
	}{
		{
			name: "",
			args: args{
				mapping: MappingPayload{
					Base{
						Index: "notexists",
						Type:  "notexists",
					},
					[]Field{
						{
							Name: "createAt",
							Type: DATE,
						},
					},
				},
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "",
			args: args{
				mapping: MappingPayload{
					Base{
						Index: "notexists1",
						Type:  "notexists1",
					},
					[]Field{
						{
							Name: "createAt",
							Type: "shoulderrortype",
						},
					},
				},
			},
			wantExists: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es.SetIndex(tt.args.mapping.Index)
			es.SetType(tt.args.mapping.Type)
			gotExists, err := es.NewIndex(context.Background(), NewMapping(tt.args.mapping))
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
