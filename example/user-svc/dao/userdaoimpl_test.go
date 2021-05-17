package dao

import (
	"context"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/unionj-cloud/go-doudou/ddl/query"
)

func TestUserDaoImpl_PageUsers(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		where query.Q
		page  query.Page
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    query.PageRet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := UserDaoImpl{
				db: tt.fields.db,
			}
			got, err := receiver.PageUsers(tt.args.ctx, tt.args.where, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDaoImpl.PageUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserDaoImpl.PageUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
