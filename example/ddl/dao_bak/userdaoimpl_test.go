package dao_bak

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"example/ddl/domain_bak"
	. "github.com/unionj-cloud/go-doudou/kit/ddl/query"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"reflect"
	"testing"
)

type DbConfig struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Charset string
}

var (
	db *sqlx.DB
)

func init() {
	err := godotenv.Load(pathutils.Abs(".env"))
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	var dbConfig DbConfig
	err = envconfig.Process("db", &dbConfig)
	if err != nil {
		log.Fatal("Error processing env", err)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbConfig.User,
		dbConfig.Passwd,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Schema,
		dbConfig.Charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	db, err = sqlx.Connect("mysql", conn)
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
				where: C().Col("name").Eq(Literal("wu")).Or(C().Col("school").Eq(Literal("havard"))).And(C().Col("age").Gte(Literal("27"))),
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
				where: C().Col("age").Gt(Literal("27")),
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

func TestUserDaoImpl_UpdateUsers(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		user  domain.User
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
			name: "UpdateUsers",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				user: domain.User{
					Name: "jones",
					Age:  60,
					No:   1998,
				},
				where: C().Col("no").Eq(Literal(1998)),
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := NewUserDao(tt.fields.db)
			got, err := receiver.UpdateUsers(tt.args.ctx, tt.args.user, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDaoImpl_UpdateUsersNoneZero(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		user  domain.User
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
			name: "UpdateUsersNoneZero",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				user: domain.User{
					Name: "jones",
					Age:  60,
					No:   1990,
				},
				where: C().Col("no").Eq(Literal(1990)),
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := UserDaoImpl{
				db: tt.fields.db,
			}
			got, err := receiver.UpdateUsersNoneZero(tt.args.ctx, tt.args.user, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUsersNoneZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateUsersNoneZero() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDaoImpl_InsertUser(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			fields: fields{
				db,
			},
			name: "InsertUser",
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					Name:   "Lemon",
					Phone:  "18811386894",
					Age:    39,
					No:     1978,
					School: "waseda",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := UserDaoImpl{
				db: tt.fields.db,
			}
			got, err := receiver.InsertUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InsertUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDaoImpl_UpdateUser(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "UpdateUser",
			fields: fields{
				db,
			},
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					ID:     15,
					Age:    50,
					School: "Beijing Univ.",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := UserDaoImpl{
				db: tt.fields.db,
			}
			got, err := receiver.UpdateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDaoImpl_UpdateUserNoneZero(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "UpdateUserNoneZero",
			fields: fields{
				db,
			},
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					ID:     13,
					Age:    50,
					School: "Beijing Univ.",
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := UserDaoImpl{
				db: tt.fields.db,
			}
			got, err := receiver.UpdateUserNoneZero(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserNoneZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateUserNoneZero() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDaoImpl_UpsertUserNoneZero(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			fields: fields{
				db,
			},
			name: "UpsertUserNoneZero",
			args: args{
				ctx: context.TODO(),
				user: &domain.User{
					No:     1990,
					School: "havard",
				},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := UserDaoImpl{
				db: tt.fields.db,
			}
			got, err := receiver.UpsertUserNoneZero(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertUserNoneZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpsertUserNoneZero() got = %v, want %v", got, tt.want)
			}
		})
	}
}
