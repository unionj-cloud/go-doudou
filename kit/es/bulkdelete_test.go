package es

import "testing"

func TestBulkDelete(t *testing.T) {
	type args struct {
		esindex string
		estype  string
		ids     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BulkDelete(tt.args.esindex, tt.args.estype, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
