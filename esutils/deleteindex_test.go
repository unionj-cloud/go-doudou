package esutils

import (
	"context"
	"testing"
)

func TestDeleteIndex(t *testing.T) {
	es, terminator := setupSubTest()
	defer terminator()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := es.DeleteIndex(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteIndex1(t *testing.T) {
	es, terminator := setupSubTest()
	defer terminator()

	es.SetIndex("notexistsindex")

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := es.DeleteIndex(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
