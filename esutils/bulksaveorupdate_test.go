package esutils

import (
	"context"
	"github.com/unionj-cloud/go-doudou/constants"
	"testing"
	"time"
)

func TestBulkSaveOrUpdate(t *testing.T) {
	es, terminator := setupSubTest()
	defer terminator()

	type args struct {
		docs []map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				docs: []map[string]interface{}{
					{
						"createAt": time.Now().In(constants.Loc).Format(constants.FORMATES),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := es.BulkSaveOrUpdate(context.Background(), tt.args.docs); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
