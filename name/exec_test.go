package name

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"testing"
)

func TestName_Exec(t *testing.T) {
	type fields struct {
		File     string
		Strategy string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "1",
			fields: fields{
				File:     pathutils.Abs("testfiles/vo.go"),
				Strategy: "lowerCaseNamingStrategy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Name{
				File:     tt.fields.File,
				Strategy: tt.fields.Strategy,
			}
			receiver.Exec()
		})
	}
}
