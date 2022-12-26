package codegen

import "testing"

func TestPattern2Method(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				pattern: "/shelves/{shelf}/books/{book}",
			},
			want: "Shelves_ShelfBooks_Book",
		},
		{
			name: "",
			args: args{
				pattern: "/goodFood/{bigApple}/books/{myBird}",
			},
			want: "Goodfood_BigappleBooks_Mybird",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pattern2Method(tt.args.pattern); got != tt.want {
				t.Errorf("Pattern2Method() = %v, want %v", got, tt.want)
			}
		})
	}
}
