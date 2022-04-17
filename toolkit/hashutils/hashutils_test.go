package hashutils

import "testing"

func TestBase64(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				input: "my secret key",
			},
			want: "bXkgc2VjcmV0IGtleQ==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64(tt.args.input); got != tt.want {
				t.Errorf("Base64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecret2Password(t *testing.T) {
	type args struct {
		username string
		secret   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				username: "wubin",
				secret:   "my secret",
			},
			want: "cff23c519b29a0e0c0304ff1a3d795f171b9c919",
		},
		{
			name: "",
			args: args{
				username: "wubin",
				secret:   "",
			},
			want: "f85610573ac9cda1a0e27e27406e9125e0e2403d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Secret2Password(tt.args.username, tt.args.secret); got != tt.want {
				t.Errorf("Secret2Password() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUIDByString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				input: "http://www.qingdao.gov.cn/zwgk/xxgk/fgw/gkml/gwfg/bmgw/",
			},
			want: "f64e7a2c-7c3c-574b-8353-e10efee0efc5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UUIDByString(tt.args.input); got != tt.want {
				t.Errorf("UUIDByString() = %v, want %v", got, tt.want)
			}
		})
	}
}
