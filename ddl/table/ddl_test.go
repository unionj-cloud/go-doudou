package table

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/ddl/config"
	"github.com/unionj-cloud/go-doudou/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/ddl/wrapper"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/test"
	"os"
	"testing"
)

func setup() (func(), *sqlx.DB, error) {
	logger := logrus.New()
	var terminateContainer func() // variable to store function to terminate container
	var host string
	var port int
	var err error
	terminateContainer, host, port, err = test.SetupMySQLContainer(logger, pathutils.Abs("../../test/sql"), "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to setup MySQL container")
	}
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", fmt.Sprint(port))
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWD", "1234")
	os.Setenv("DB_SCHEMA", "test")
	os.Setenv("DB_CHARSET", "utf8mb4")
	var conf config.DbConfig
	err = envconfig.Process("db", &conf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error processing env")
	}
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.User,
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Schema,
		conf.Charset)
	conn += `&loc=Asia%2FShanghai&parseTime=True`
	var db *sqlx.DB
	db, err = sqlx.Connect("mysql", conn)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	db.MapperFunc(strcase.ToSnake)
	return terminateContainer, db, nil
}

func ExampleCreateTable() {
	terminator, db, err := setup()
	if err != nil {
		panic(err)
	}
	defer terminator()
	defer db.Close()

	expectjson := `{"Name":"user_createtable","Columns":[{"Table":"user","Name":"id","Type":"INT","Default":null,"Pk":true,"Nullable":false,"Unsigned":false,"Autoincrement":true,"Extra":"","Meta":{"Name":"ID","Type":"int","Tag":"dd:\"pk;auto\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"name","Type":"VARCHAR(255)","Default":"'jack'","Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"Name","Type":"string","Tag":"dd:\"index:name_phone_idx,2;default:'jack'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"phone","Type":"VARCHAR(255)","Default":"'13552053960'","Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"comment '手机号'","Meta":{"Name":"Phone","Type":"string","Tag":"dd:\"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"age","Type":"INT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"Age","Type":"int","Tag":"dd:\"index\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"no","Type":"INT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"No","Type":"int","Tag":"dd:\"unique\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"school","Type":"VARCHAR(255)","Default":"'harvard'","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"comment '学校'","Meta":{"Name":"School","Type":"string","Tag":"dd:\"null;default:'harvard';extra:comment '学校'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"is_student","Type":"TINYINT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"IsStudent","Type":"bool","Tag":"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"delete_at","Type":"DATETIME","Default":null,"Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"DeleteAt","Type":"*time.Time","Tag":"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"create_at","Type":"DATETIME","Default":"CURRENT_TIMESTAMP","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"CreateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP\"","Comments":null},"AutoSet":true,"Indexes":null},{"Table":"user","Name":"update_at","Type":"DATETIME","Default":"CURRENT_TIMESTAMP","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"ON UPDATE CURRENT_TIMESTAMP","Meta":{"Name":"UpdateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP\"","Comments":null},"AutoSet":true,"Indexes":null}],"Pk":"id","Indexes":[{"Unique":false,"Name":"name_phone_idx","Items":[{"Unique":false,"Name":"","Column":"phone","Order":1,"Sort":"asc"},{"Unique":false,"Name":"","Column":"name","Order":2,"Sort":"asc"}]},{"Unique":false,"Name":"age_idx","Items":[{"Unique":false,"Name":"","Column":"age","Order":1,"Sort":"asc"}]},{"Unique":true,"Name":"no_idx","Items":[{"Unique":false,"Name":"","Column":"no","Order":1,"Sort":"asc"}]}],"Meta":{"Name":"User","Fields":[{"Name":"ID","Type":"int","Tag":"dd:\"pk;auto\"","Comments":null},{"Name":"Name","Type":"string","Tag":"dd:\"index:name_phone_idx,2;default:'jack'\"","Comments":null},{"Name":"Phone","Type":"string","Tag":"dd:\"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'\"","Comments":null},{"Name":"Age","Type":"int","Tag":"dd:\"index\"","Comments":null},{"Name":"No","Type":"int","Tag":"dd:\"unique\"","Comments":null},{"Name":"School","Type":"string","Tag":"dd:\"null;default:'harvard';extra:comment '学校'\"","Comments":null},{"Name":"IsStudent","Type":"bool","Tag":"","Comments":null},{"Name":"DeleteAt","Type":"*time.Time","Tag":"","Comments":null},{"Name":"CreateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP\"","Comments":null},{"Name":"UpdateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP\"","Comments":null}],"Comments":["dd:table"],"Methods":null}}`
	var table Table
	if err = json.Unmarshal([]byte(expectjson), &table); err != nil {
		panic(err)
	}
	if err := CreateTable(context.Background(), db, table); (err != nil) != false {
		panic(fmt.Sprintf("CreateTable() error = %v, wantErr %v", err, false))
	}

	// Output:

}

func TestChangeColumn(t *testing.T) {
	terminator, db, err := setup()
	if err != nil {
		panic(err)
	}
	defer terminator()
	defer db.Close()

	type args struct {
		db  *sqlx.DB
		col Column
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errmsg  string
	}{
		{
			name: "1",
			args: args{
				db: db,
				col: Column{
					Table:   "ddl_user",
					Name:    "school",
					Type:    "varchar(45)",
					Default: "'Beijing Univ.'",
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				db: db,
				col: Column{
					Table:   "ddl_user",
					Name:    "school",
					Type:    columnenum.TextType,
					Default: "'Beijing Univ.'",
				},
			},
			wantErr: true,
			errmsg:  `Error 1101: BLOB, TEXT, GEOMETRY or JSON column 'school' can't have a default value`,
		},
		{
			name: "3",
			args: args{
				db: db,
				col: Column{
					Table:   "ddl_user",
					Name:    "school",
					Type:    "varchar(45)",
					Default: "Beijing Univ.",
				},
			},
			wantErr: true,
			errmsg:  `Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'Beijing Univ.' at line 2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = ChangeColumn(context.Background(), tt.args.db, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("ChangeColumn() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if err.Error() != tt.errmsg {
					t.Errorf("want %s, got %s", tt.errmsg, err.Error())
				}
			}
		})
	}
}

func TestAddColumn(t *testing.T) {
	terminator, db, err := setup()
	if err != nil {
		panic(err)
	}
	defer terminator()
	defer db.Close()

	type args struct {
		db  *sqlx.DB
		col Column
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errmsg  string
	}{
		{
			name: "1",
			args: args{
				db: db,
				col: Column{
					Table:   "ddl_user",
					Name:    "favourite",
					Type:    "varchar(45)",
					Default: "'football'",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = AddColumn(context.Background(), tt.args.db, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("ChangeColumn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDropIndex(t *testing.T) {
	terminator, db, err := setup()
	if err != nil {
		panic(err)
	}
	defer terminator()
	defer db.Close()

	type args struct {
		ctx context.Context
		db  wrapper.Querier
		idx Index
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				db:  db,
				idx: Index{
					Table:  "ddl_user",
					Unique: true,
					Name:   "age_idx",
					Items: []IndexItem{
						{
							Column: "age",
							Order:  1,
							Sort:   sortenum.Asc,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DropIndex(tt.args.ctx, tt.args.db, tt.args.idx); (err != nil) != tt.wantErr {
				t.Errorf("DropIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddIndex(t *testing.T) {
	terminator, db, err := setup()
	if err != nil {
		panic(err)
	}
	defer terminator()
	defer db.Close()

	type args struct {
		ctx context.Context
		db  wrapper.Querier
		idx Index
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				db:  db,
				idx: Index{
					Table:  "ddl_user",
					Unique: true,
					Name:   "school_idx",
					Items: []IndexItem{
						{
							Column: "school",
							Order:  1,
							Sort:   sortenum.Asc,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddIndex(tt.args.ctx, tt.args.db, tt.args.idx); (err != nil) != tt.wantErr {
				t.Errorf("AddIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDropAddIndex(t *testing.T) {
	terminator, db, err := setup()
	if err != nil {
		panic(err)
	}
	defer terminator()
	defer db.Close()

	type args struct {
		ctx context.Context
		db  wrapper.Querier
		idx Index
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				db:  db,
				idx: Index{
					Table:  "ddl_user",
					Unique: true,
					Name:   "age_idx",
					Items: []IndexItem{
						{
							Column: "age",
							Order:  1,
							Sort:   sortenum.Asc,
						},
						{
							Column: "school",
							Order:  2,
							Sort:   sortenum.Asc,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DropAddIndex(tt.args.ctx, tt.args.db, tt.args.idx); (err != nil) != tt.wantErr {
				t.Errorf("DropAddIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
