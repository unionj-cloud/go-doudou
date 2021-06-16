package stringutils

import "testing"

func TestContains(t *testing.T) {
	type args struct {
		s      string
		substr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				s:      "filedownloadUser",
				substr: "Download",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsI(tt.args.s, tt.args.substr); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasPrefixI(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				s:      "VARCHAR(255)",
				prefix: "var",
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				s:      "VARCHAR(255)",
				prefix: "VAR",
			},
			want: true,
		},
		{
			name: "3",
			args: args{
				s:      "VARCHAR(255)",
				prefix: "CHA",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPrefixI(tt.args.s, tt.args.prefix); got != tt.want {
				t.Errorf("HasPrefixI() = %v, want %v", got, tt.want)
			}
		})
	}
}
