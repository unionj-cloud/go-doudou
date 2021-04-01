package dao

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/kit/ddl/example/domain"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"os"
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

var dbConfig DbConfig

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
}

func TestUserDaoGen_UpsertUser(t *testing.T) {
	var (
		db  *sqlx.DB
		err error
	)
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

	type args struct {
		ctx  context.Context
		db   *sqlx.DB
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
				db:  db,
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
			u := UserDaoGen{}
			got, err := u.UpsertUser(tt.args.ctx, tt.args.db, tt.args.user)
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
	var (
		db  *sqlx.DB
		err error
	)
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

	type args struct {
		ctx  context.Context
		db   *sqlx.DB
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
				db:  db,
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
			u := UserDaoGen{}
			got, err := u.UpsertUser(tt.args.ctx, tt.args.db, tt.args.user)
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
