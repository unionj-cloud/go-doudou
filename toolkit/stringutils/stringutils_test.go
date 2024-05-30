package stringutils

import (
	"testing"
)

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

func TestReplaceStringAtIndex(t *testing.T) {
	in := `INSERT INTO user ("name", "age") VALUES (?, ?);`
	loc := IndexAll(in, "?", -1)

	in2 := `我爱北京天安门，啦啦啦`
	loc2 := IndexAll(in2, "天安门", -1)

	type args struct {
		in      string
		replace string
		start   int
		end     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				in:      in,
				replace: "18",
				start:   loc[1][0],
				end:     loc[1][1],
			},
			want: `INSERT INTO user ("name", "age") VALUES (?, 18);`,
		},
		{
			name: "",
			args: args{
				in:      in2,
				replace: "颐和园",
				start:   loc2[0][0],
				end:     loc2[0][1],
			},
			want: `我爱北京颐和园，啦啦啦`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceStringAtByteIndex(tt.args.in, tt.args.replace, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("ReplaceStringAtIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceStringAtByteIndexBatch(t *testing.T) {
	in := `INSERT INTO user ("name", "age") VALUES (?, ?);`
	loc := IndexAll(in, "?", -1)

	in2 := `我爱北京天安门，啦啦啦`
	loc2 := IndexAll(in2, "天安门", -1)

	type args struct {
		in   string
		args []string
		locs [][]int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				in:   in,
				args: []string{"'wubin'", "18"},
				locs: loc,
			},
			want: `INSERT INTO user ("name", "age") VALUES ('wubin', 18);`,
		},
		{
			name: "",
			args: args{
				in:   in2,
				args: []string{"颐和园"},
				locs: loc2,
			},
			want: `我爱北京颐和园，啦啦啦`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceStringAtByteIndexBatch(tt.args.in, tt.args.args, tt.args.locs); got != tt.want {
				t.Errorf("ReplaceStringAtByteIndexBatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
