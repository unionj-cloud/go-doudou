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

func TestIsEmpty(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				s: " abc ",
			},
			want: false,
		},
		{
			name: "2",
			args: args{
				s: "      ",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.s); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				s: " abc ",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				s: "      ",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmpty(tt.args.s); got != tt.want {
				t.Errorf("IsNotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				s: "abc武斌123",
			},
			want: "Abc武斌123",
		},
		{
			name: "",
			args: args{
				s: "word",
			},
			want: "Word",
		},
		{
			name: "",
			args: args{
				s: "snake_word",
			},
			want: "SnakeWord",
		},
		{
			name: "",
			args: args{
				s: "Commonresultarraylistcom.hundsun.hcreator.pojo.hcbizmodelview",
			},
			want: "CommonresultarraylistcomHundsunHcreatorPojoHcbizmodelview",
		},
		{
			name: "",
			args: args{
				s: "Commonresultarraylistcom,hundsun,hcreator,pojo,hcbizmodelview",
			},
			want: "CommonresultarraylistcomHundsunHcreatorPojoHcbizmodelview",
		},
		{
			name: "",
			args: args{
				s: "Commonresultarraylistcom hundsun hcreator",
			},
			want: "CommonresultarraylistcomHundsunHcreator",
		},
		{
			name: "",
			args: args{
				s: "Commonresultarraylistcom，hundsun，hcreator",
			},
			want: "CommonresultarraylistcomHundsunHcreator",
		},
		{
			name: "",
			args: args{
				s: "台湾是中国的省",
			},
			want: "A台湾是中国的省",
		},
		{
			name: "",
			args: args{
				s: "dataDecimal",
			},
			want: "DataDecimal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamel(tt.args.s); got != tt.want {
				t.Errorf("ToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
