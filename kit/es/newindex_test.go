package es

import "testing"

func TestNewIndex(t *testing.T) {
	type args struct {
		index   string
		mapping string
	}
	index := "team3_voice_analysis_wb"
	mapping := NewMapping(MappingPayload{
		Base{
			Index: "team3_voice_analysis_wb",
			Type:  "team3_voice_analysis_wb",
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
				index:   index,
				mapping: mapping,
			},
			wantExists: false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := NewIndex(tt.args.index, tt.args.mapping)
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
