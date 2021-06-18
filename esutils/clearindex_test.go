package esutils

import (
	"context"
	"testing"
)

func TestClearIndex(t *testing.T) {
	es, terminator := setupSubTest()
	defer terminator()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := es.ClearIndex(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("ClearIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
