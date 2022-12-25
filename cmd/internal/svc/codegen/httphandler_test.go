package codegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_pattern(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				method: "GetBooks",
			},
			want: "books",
		},
		{
			name: "2",
			args: args{
				method: "PageUsers",
			},
			want: "page/users",
		},
		{
			name: "",
			args: args{
				method: "Get",
			},
			want: "",
		},
		{
			name: "",
			args: args{
				method: "GetShelves_ShelfBooks_Book",
			},
			want: "shelves/:shelf/books/:book",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, pattern(tt.args.method), "pattern(%v)", tt.args.method)
		})
	}
}

// Shelves_ShelfBooks_Book
func Test_httpMethod(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				method: "GetShelves_ShelfBooks_Book",
			},
			want: "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, httpMethod(tt.args.method), "httpMethod(%v)", tt.args.method)
		})
	}
}
