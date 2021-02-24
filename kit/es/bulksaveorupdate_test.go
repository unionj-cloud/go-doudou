package es

import (
	"testing"
	"time"
	"github.com/unionj-cloud/papilio/kit/constants"
)

func TestBulkSaveOrUpdate(t *testing.T) {

	const index = "team3_voice_analysis_wb"

	teardownSubTest := SetupSubTest(index, t)
	defer teardownSubTest(t)

	type args struct {
		esindex string
		estype  string
		docs    []map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				esindex: index,
				estype:  index,
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
			if err := BulkSaveOrUpdate(tt.args.esindex, tt.args.estype, tt.args.docs); (err != nil) != tt.wantErr {
				t.Errorf("BulkSaveOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
