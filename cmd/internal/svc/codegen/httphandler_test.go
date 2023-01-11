package codegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
