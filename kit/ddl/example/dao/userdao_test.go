package dao

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/unionj-cloud/go-doudou/kit/ddl/example/domain"
	"reflect"
	"testing"
)

func TestUserDao_UpsertUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "upsert",
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					Name:   "wu",
					Phone:  "13552053960",
					Age:    35,
					No:     1988,
					School: "havard",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserDao{
				NewUserDaoGen(db),
			}
			got, err := u.UpsertUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpsertUser() got = %v, want %v", got, tt.want)
			}
			fmt.Println(tt.args.user.ID)
		})
	}
}

func TestUserDao_UpsertUser1(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "upsert",
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					Name:   "david",
					Phone:  "13552053960",
					Age:    35,
					No:     1997,
					School: "havard",
				},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserDao{
				NewUserDaoGen(db),
			}
			got, err := u.UpsertUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpsertUser() got = %v, want %v", got, tt.want)
			}
			fmt.Println(tt.args.user.ID)
		})
	}
}

func TestUserDao_GetUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "getUser",
			args: args{
				ctx: context.TODO(),
				id:  6,
			},
			want:    "wu",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserDao{
				NewUserDaoGen(db),
			}
			got, err := u.GetUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Name, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
