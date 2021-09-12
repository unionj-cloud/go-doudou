package pathutils

import (
	"strings"
	"testing"
)

func TestAbs(t *testing.T) {
	type args struct {
		rel string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				rel: "testfiles",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.args.rel); len(got) == 0 {
				t.Errorf("Abs() got nothing")
			}
		})
	}
}

func TestFixPath(t *testing.T) {
	type args struct {
		dir      string
		fallback string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				dir:      "testfiles",
				fallback: "fallback",
			},
			want:    "",
			wantErr: false,
		},
		{
			args: args{
				dir:      "",
				fallback: "fallback",
			},
			want:    "fallback",
			wantErr: false,
		},
		{
			args: args{
				dir:      "/absolute/path",
				fallback: "",
			},
			want:    "/absolute/path",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FixPath(tt.args.dir, tt.args.fallback)
			if (err != nil) != tt.wantErr {
				t.Errorf("FixPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasSuffix(got, tt.want) {
				t.Errorf("FixPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
