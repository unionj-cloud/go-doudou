package dao_bak

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/ddl/example/domain"
	. "github.com/unionj-cloud/go-doudou/kit/ddl/query"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"os"
	"reflect"
	"testing"
)

type DbConfig struct {
	host    string
	port    string
	user    string
	passwd  string
	schema  string
	charset string
}

var (
	db       *sqlx.DB
	dbConfig DbConfig
)

func init() {
	err := godotenv.Load(pathutils.Abs(".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConfig = DbConfig{
		host:    os.Getenv("DB_HOST"),
		port:    os.Getenv("DB_PORT"),
		user:    os.Getenv("DB_USER"),
		passwd:  os.Getenv("DB_PASSWD"),
		schema:  os.Getenv("DB_SCHEMA"),
		charset: os.Getenv("DB_CHARSET"),
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbConfig.user,
		dbConfig.passwd,
		dbConfig.host,
		dbConfig.port,
		dbConfig.schema,
		dbConfig.charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	db, err = sqlx.ConnectContext(context.TODO(), "mysql", conn)
	if err != nil {
		log.Fatalln(err)
	}
	db.MapperFunc(strcase.ToSnake)
}

func TestUserDaoGen_UpsertUser(t *testing.T) {
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
					Age:    18,
					No:     1990,
					School: "havard",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserDao(db)
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

func TestUserDaoGen_UpsertUser1(t *testing.T) {
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
			u := NewUserDao(db)
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

func TestUserDaoGen_GetUser(t *testing.T) {
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
			u := NewUserDao(db)
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

func TestUserDaoGen_DeleteUsers(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		where Q
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "deleteUsers",
			fields: fields{
				db,
			},
			args: args{
				ctx:   context.TODO(),
				where: C().Col("name").Eq("wu").Or(C().Col("school").Eq("havard")).And(C().Col("age").Gte("27")),
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserDao(db)
			got, err := u.DeleteUsers(tt.args.ctx, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDaoGen_PageUsers(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		where Q
		page  Page
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    PageRet
		wantErr bool
	}{
		{
			name: "pageUsers",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:   context.TODO(),
				where: C().Col("age").Gt("27"),
				page: Page{
					Orders: []Order{
						{
							Col:  "age",
							Sort: "desc",
						},
					},
					Offset: 2,
					Size:   1,
				},
			},
			want: PageRet{
				PageNo:   3,
				PageSize: 1,
				Total:    3,
				HasNext:  false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUserDao(db)
			got, err := u.PageUsers(tt.args.ctx, tt.args.where, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var (
				items []domain.User
				ok    bool
			)
			if items, ok = got.Items.([]domain.User); !ok {
				t.Errorf("PageUsers() got = %v, want %v", got, tt.want)
			}
			if len(items) == 0 {
				t.Errorf("PageUsers() got = %v, want %v", got, tt.want)
			}
			if items[0].ID != 11 {
				t.Errorf("PageUsers() got = %v, want %v", got, tt.want)
			}
			if tt.want.HasNext != got.HasNext {
				t.Errorf("PageUsers() got = %v, want %v", got, tt.want)
			}
			if tt.want.PageNo != got.PageNo {
				t.Errorf("PageUsers() got = %v, want %v", got, tt.want)
			}
			if tt.want.PageSize != got.PageSize {
				t.Errorf("PageUsers() got = %v, want %v", got, tt.want)
			}
			if tt.want.Total != got.Total {
				t.Errorf("PageUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
