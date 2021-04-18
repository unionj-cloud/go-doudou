package name

import "testing"

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
				File:     "/Users/wubin1989/workspace/cloud/go-doudou/example/doudou/vo/vo.go",
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
