package es

import "testing"

func TestDeleteIndex(t *testing.T) {
	const index = "team3_voice_analysis_wb"

	teardownSubTest := SetupSubTest(index, t)
	defer teardownSubTest(t)

	type args struct {
		index string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				index: index,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteIndex(tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
