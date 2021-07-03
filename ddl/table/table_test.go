package table

import (
	"encoding/json"
	"fmt"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/ddl/extraenum"
	"github.com/unionj-cloud/go-doudou/ddl/keyenum"
	"github.com/unionj-cloud/go-doudou/ddl/nullenum"
	"github.com/unionj-cloud/go-doudou/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewTableFromStruct(t *testing.T) {
	testDir := pathutils.Abs("../testfiles/domain")
	expectjson := `{"Name":"user","Columns":[{"Table":"user","Name":"id","Type":"INT","Default":null,"Pk":true,"Nullable":false,"Unsigned":false,"Autoincrement":true,"Extra":"","Meta":{"Name":"ID","Type":"int","Tag":"dd:\"pk;auto\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"name","Type":"VARCHAR(255)","Default":"'jack'","Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"Name","Type":"string","Tag":"dd:\"index:name_phone_idx,2;default:'jack'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"phone","Type":"VARCHAR(255)","Default":"'13552053960'","Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"comment '手机号'","Meta":{"Name":"Phone","Type":"string","Tag":"dd:\"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"age","Type":"INT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"Age","Type":"int","Tag":"dd:\"index\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"no","Type":"INT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"No","Type":"int","Tag":"dd:\"unique\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"school","Type":"VARCHAR(255)","Default":"'harvard'","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"comment '学校'","Meta":{"Name":"School","Type":"string","Tag":"dd:\"null;default:'harvard';extra:comment '学校'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"is_student","Type":"TINYINT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"IsStudent","Type":"bool","Tag":"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"delete_at","Type":"DATETIME","Default":null,"Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"DeleteAt","Type":"*time.Time","Tag":"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"create_at","Type":"DATETIME","Default":"CURRENT_TIMESTAMP","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"CreateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP\"","Comments":null},"AutoSet":true,"Indexes":null},{"Table":"user","Name":"update_at","Type":"DATETIME","Default":"CURRENT_TIMESTAMP","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"ON UPDATE CURRENT_TIMESTAMP","Meta":{"Name":"UpdateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP\"","Comments":null},"AutoSet":true,"Indexes":null}],"Pk":"id","Indexes":[{"Unique":false,"Name":"name_phone_idx","Items":[{"Unique":false,"Name":"","Column":"phone","Order":1,"Sort":"asc"},{"Unique":false,"Name":"","Column":"name","Order":2,"Sort":"asc"}]},{"Unique":false,"Name":"age_idx","Items":[{"Unique":false,"Name":"","Column":"age","Order":1,"Sort":"asc"}]},{"Unique":true,"Name":"no_idx","Items":[{"Unique":false,"Name":"","Column":"no","Order":1,"Sort":"asc"}]}],"Meta":{"Name":"User","Fields":[{"Name":"ID","Type":"int","Tag":"dd:\"pk;auto\"","Comments":null},{"Name":"Name","Type":"string","Tag":"dd:\"index:name_phone_idx,2;default:'jack'\"","Comments":null},{"Name":"Phone","Type":"string","Tag":"dd:\"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'\"","Comments":null},{"Name":"Age","Type":"int","Tag":"dd:\"index\"","Comments":null},{"Name":"No","Type":"int","Tag":"dd:\"unique\"","Comments":null},{"Name":"School","Type":"string","Tag":"dd:\"null;default:'harvard';extra:comment '学校'\"","Comments":null},{"Name":"IsStudent","Type":"bool","Tag":"","Comments":null},{"Name":"DeleteAt","Type":"*time.Time","Tag":"","Comments":null},{"Name":"CreateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP\"","Comments":null},{"Name":"UpdateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP\"","Comments":null}],"Comments":["dd:table"],"Methods":null}}`
	var expect Table
	json.Unmarshal([]byte(expectjson), &expect)
	var files []string
	var err error
	err = filepath.Walk(testDir, astutils.Visit(&files))
	if err != nil {
		panic(err)
	}
	sc := astutils.NewStructCollector(astutils.ExprString)
	for _, file := range files {
		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		ast.Walk(sc, root)
	}
	flattened := ddlast.FlatEmbed(sc.Structs)

	for _, sm := range flattened {
		table := NewTableFromStruct(sm)
		if len(table.Columns) != len(expect.Columns) {
			t.Errorf("NewTableFromStruct() got = %v, want %v", len(table.Columns), len(expect.Columns))
		}
	}
}

func TestTable_CreateSql(t1 *testing.T) {
	type fields struct {
		Name          string
		Columns       []Column
		Pk            string
		UniqueIndexes []Index
		Indexes       []Index
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "users",
			fields: fields{
				Name: "users",
				Columns: []Column{
					{
						Name:          "id",
						Type:          columnenum.IntType,
						Default:       nil,
						Pk:            true,
						Nullable:      false,
						Unsigned:      false,
						Autoincrement: true,
						Extra:         "",
					},
					{
						Name:          "name",
						Type:          columnenum.VarcharType,
						Default:       "'wubin'",
						Pk:            false,
						Nullable:      true,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "",
					},
					{
						Name:          "phone",
						Type:          columnenum.VarcharType,
						Default:       "'13552053960'",
						Pk:            false,
						Nullable:      true,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "comment '手机号'",
					},
					{
						Name:          "age",
						Type:          columnenum.IntType,
						Default:       nil,
						Pk:            false,
						Nullable:      true,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "",
					},
					{
						Name:          "no",
						Type:          columnenum.IntType,
						Default:       nil,
						Pk:            false,
						Nullable:      false,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "",
					},
				},
				Pk: "id",
				Indexes: []Index{
					{
						Name: "name_phone_idx",
						Items: []IndexItem{
							{
								Column: "name",
								Order:  2,
								Sort:   "asc",
							},
							{
								Column: "phone",
								Order:  1,
								Sort:   "desc",
							},
						},
					},
					{
						Unique: true,
						Name:   "uni_no",
						Items: []IndexItem{
							{
								Column: "no",
								Order:  0,
								Sort:   "asc",
							},
						},
					},
				},
			},
			want:    "CREATE TABLE `users` (\n`id` INT NOT NULL AUTO_INCREMENT,\n`name` VARCHAR(255) NULL DEFAULT 'wubin',\n`phone` VARCHAR(255) NULL DEFAULT '13552053960' comment '手机号',\n`age` INT NULL,\n`no` INT NOT NULL,\nPRIMARY KEY (`id`),\nINDEX `name_phone_idx` (`name` asc,`phone` desc),\nUNIQUE INDEX `uni_no` (`no` asc));",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:    tt.fields.Name,
				Columns: tt.fields.Columns,
				Pk:      tt.fields.Pk,
				Indexes: tt.fields.Indexes,
			}
			got, err := t.CreateSql()
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CreateSql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("CreateSql() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumn_AlterColumnSql(t *testing.T) {
	type fields struct {
		Table         string
		Name          string
		Type          columnenum.ColumnType
		Default       interface{}
		Pk            bool
		Nullable      bool
		Unsigned      bool
		Autoincrement bool
		Extra         extraenum.Extra
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "column",
			fields: fields{
				Table:         "users",
				Name:          "phone",
				Type:          columnenum.VarcharType,
				Default:       "'13552053960'",
				Pk:            false,
				Nullable:      false,
				Unsigned:      false,
				Autoincrement: false,
				Extra:         "comment '手机号'",
			},
			want:    "ALTER TABLE `users`\nCHANGE COLUMN `phone` `phone` VARCHAR(255) NOT NULL DEFAULT '13552053960' comment '手机号';",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Column{
				Table:         tt.fields.Table,
				Name:          tt.fields.Name,
				Type:          tt.fields.Type,
				Default:       tt.fields.Default,
				Pk:            tt.fields.Pk,
				Nullable:      tt.fields.Nullable,
				Unsigned:      tt.fields.Unsigned,
				Autoincrement: tt.fields.Autoincrement,
				Extra:         tt.fields.Extra,
			}
			got, err := c.ChangeColumnSql()
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeColumnSql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChangeColumnSql() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumn_AddColumnSql(t *testing.T) {
	type fields struct {
		Table         string
		Name          string
		Type          columnenum.ColumnType
		Default       interface{}
		Pk            bool
		Nullable      bool
		Unsigned      bool
		Autoincrement bool
		Extra         extraenum.Extra
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "column",
			fields: fields{
				Table:         "users",
				Name:          "school",
				Type:          columnenum.VarcharType,
				Default:       "'harvard'",
				Pk:            false,
				Nullable:      false,
				Unsigned:      false,
				Autoincrement: false,
				Extra:         "comment '学校'",
			},
			want:    "ALTER TABLE `users`\nADD COLUMN `school` VARCHAR(255) NOT NULL DEFAULT 'harvard' comment '学校';",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Column{
				Table:         tt.fields.Table,
				Name:          tt.fields.Name,
				Type:          tt.fields.Type,
				Default:       tt.fields.Default,
				Pk:            tt.fields.Pk,
				Nullable:      tt.fields.Nullable,
				Unsigned:      tt.fields.Unsigned,
				Autoincrement: tt.fields.Autoincrement,
				Extra:         tt.fields.Extra,
			}
			got, err := c.AddColumnSql()
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddColumnSql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddColumnSql() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toColumnType(t *testing.T) {
	type args struct {
		goType string
	}
	tests := []struct {
		name string
		args args
		want columnenum.ColumnType
	}{
		{
			name: "1",
			args: args{
				goType: "float32",
			},
			want: columnenum.FloatType,
		}, {
			name: "2",
			args: args{
				goType: "int",
			},
			want: columnenum.IntType,
		}, {
			name: "3",
			args: args{
				goType: "bool",
			},
			want: columnenum.TinyintType,
		}, {
			name: "4",
			args: args{
				goType: "time.Time",
			},
			want: columnenum.DatetimeType,
		}, {
			name: "5",
			args: args{
				goType: "int64",
			},
			want: columnenum.BigintType,
		}, {
			name: "6",
			args: args{
				goType: "float64",
			},
			want: columnenum.DoubleType,
		}, {
			name: "7",
			args: args{
				goType: "string",
			},
			want: columnenum.VarcharType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toColumnType(tt.args.goType); got != tt.want {
				t.Errorf("toColumnType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toGoType(t *testing.T) {
	type args struct {
		colType  columnenum.ColumnType
		nullable bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				colType:  "int",
				nullable: false,
			},
			want: "int",
		},
		{
			name: "2",
			args: args{
				colType:  "bigint",
				nullable: false,
			},
			want: "int64",
		},
		{
			name: "3",
			args: args{
				colType:  "float",
				nullable: false,
			},
			want: "float32",
		},
		{
			name: "4",
			args: args{
				colType:  "double",
				nullable: false,
			},
			want: "float64",
		},
		{
			name: "5",
			args: args{
				colType:  "varchar",
				nullable: false,
			},
			want: "string",
		},
		{
			name: "6",
			args: args{
				colType:  "tinyint",
				nullable: false,
			},
			want: "int8",
		},
		{
			name: "7",
			args: args{
				colType:  "text",
				nullable: false,
			},
			want: "string",
		},
		{
			name: "8",
			args: args{
				colType:  "datetime",
				nullable: false,
			},
			want: "time.Time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toGoType(tt.args.colType, tt.args.nullable); got != tt.want {
				t.Errorf("toGoType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFieldFromColumn(t *testing.T) {
	type args struct {
		col Column
	}
	tests := []struct {
		name string
		args args
		want astutils.FieldMeta
	}{
		{
			name: "1",
			args: args{
				col: Column{
					Table:         "users",
					Name:          "school",
					Type:          columnenum.VarcharType,
					Default:       "'harvard'",
					Pk:            false,
					Nullable:      false,
					Unsigned:      false,
					Autoincrement: false,
					Extra:         "comment '学校'",
				},
			},
			want: astutils.FieldMeta{
				Name:     "School",
				Type:     "string",
				Tag:      `dd:"type:VARCHAR(255);extra:comment '学校'"`,
				Comments: nil,
			},
		}, {
			name: "2",
			args: args{
				col: Column{
					Table:         "users",
					Name:          "favourite",
					Type:          columnenum.VarcharType,
					Default:       "'football'",
					Pk:            false,
					Nullable:      true,
					Unsigned:      false,
					Autoincrement: true,
					Extra:         "comment '学校'",
					Indexes: IndexItems{
						{
							Unique: false,
							Name:   "my_index",
							Column: "favourite",
							Order:  1,
							Sort:   sortenum.Asc,
						},
					},
				},
			},
			want: astutils.FieldMeta{
				Name:     "Favourite",
				Type:     "*string",
				Tag:      `dd:"auto;type:VARCHAR(255);extra:comment '学校';index:my_index,1,asc"`,
				Comments: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFieldFromColumn(tt.args.col); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFieldFromColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPk(t *testing.T) {
	type args struct {
		key keyenum.Key
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				key: keyenum.Pri,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPk(tt.args.key); got != tt.want {
				t.Errorf("CheckPk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckNull(t *testing.T) {
	type args struct {
		null nullenum.Null
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				null: nullenum.Yes,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckNull(tt.args.null); got != tt.want {
				t.Errorf("CheckNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUnsigned(t *testing.T) {
	type args struct {
		dbColType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				dbColType: "int unsigned",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUnsigned(tt.args.dbColType); got != tt.want {
				t.Errorf("CheckUnsigned() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckAutoincrement(t *testing.T) {
	type args struct {
		extra string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				extra: "auto_increment",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckAutoincrement(tt.args.extra); got != tt.want {
				t.Errorf("CheckAutoincrement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckAutoSet(t *testing.T) {
	s := "CURRENT_TIMESTAMP"
	type args struct {
		defaultVal *string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				defaultVal: &s,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckAutoSet(tt.args.defaultVal); got != tt.want {
				t.Errorf("CheckAutoSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTableFromStruct1(t *testing.T) {
	file := pathutils.Abs("../testfiles/domain/purchase.go")
	expectjson := `{"Name":"user","Columns":[{"Table":"user","Name":"id","Type":"INT","Default":null,"Pk":true,"Nullable":false,"Unsigned":false,"Autoincrement":true,"Extra":"","Meta":{"Name":"ID","Type":"int","Tag":"dd:\"pk;auto\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"name","Type":"VARCHAR(255)","Default":"'jack'","Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"Name","Type":"string","Tag":"dd:\"index:name_phone_idx,2;default:'jack'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"phone","Type":"VARCHAR(255)","Default":"'13552053960'","Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"comment '手机号'","Meta":{"Name":"Phone","Type":"string","Tag":"dd:\"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"age","Type":"INT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"Age","Type":"int","Tag":"dd:\"index\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"no","Type":"INT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"No","Type":"int","Tag":"dd:\"unique\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"school","Type":"VARCHAR(255)","Default":"'harvard'","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"comment '学校'","Meta":{"Name":"School","Type":"string","Tag":"dd:\"null;default:'harvard';extra:comment '学校'\"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"is_student","Type":"TINYINT","Default":null,"Pk":false,"Nullable":false,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"IsStudent","Type":"bool","Tag":"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"delete_at","Type":"DATETIME","Default":null,"Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"DeleteAt","Type":"*time.Time","Tag":"","Comments":null},"AutoSet":false,"Indexes":null},{"Table":"user","Name":"create_at","Type":"DATETIME","Default":"CURRENT_TIMESTAMP","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"","Meta":{"Name":"CreateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP\"","Comments":null},"AutoSet":true,"Indexes":null},{"Table":"user","Name":"update_at","Type":"DATETIME","Default":"CURRENT_TIMESTAMP","Pk":false,"Nullable":true,"Unsigned":false,"Autoincrement":false,"Extra":"ON UPDATE CURRENT_TIMESTAMP","Meta":{"Name":"UpdateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP\"","Comments":null},"AutoSet":true,"Indexes":null}],"Pk":"id","Indexes":[{"Unique":false,"Name":"name_phone_idx","Items":[{"Unique":false,"Name":"","Column":"phone","Order":1,"Sort":"asc"},{"Unique":false,"Name":"","Column":"name","Order":2,"Sort":"asc"}]},{"Unique":false,"Name":"age_idx","Items":[{"Unique":false,"Name":"","Column":"age","Order":1,"Sort":"asc"}]},{"Unique":true,"Name":"no_idx","Items":[{"Unique":false,"Name":"","Column":"no","Order":1,"Sort":"asc"}]}],"Meta":{"Name":"User","Fields":[{"Name":"ID","Type":"int","Tag":"dd:\"pk;auto\"","Comments":null},{"Name":"Name","Type":"string","Tag":"dd:\"index:name_phone_idx,2;default:'jack'\"","Comments":null},{"Name":"Phone","Type":"string","Tag":"dd:\"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'\"","Comments":null},{"Name":"Age","Type":"int","Tag":"dd:\"index\"","Comments":null},{"Name":"No","Type":"int","Tag":"dd:\"unique\"","Comments":null},{"Name":"School","Type":"string","Tag":"dd:\"null;default:'harvard';extra:comment '学校'\"","Comments":null},{"Name":"IsStudent","Type":"bool","Tag":"","Comments":null},{"Name":"DeleteAt","Type":"*time.Time","Tag":"","Comments":null},{"Name":"CreateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP\"","Comments":null},{"Name":"UpdateAt","Type":"*time.Time","Tag":"dd:\"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP\"","Comments":null}],"Comments":["dd:table"],"Methods":null}}`
	var expect Table
	json.Unmarshal([]byte(expectjson), &expect)
	sc := astutils.NewStructCollector(astutils.ExprString)
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Walk(sc, root)
	flattened := ddlast.FlatEmbed(sc.Structs)

	for _, sm := range flattened {
		table := NewTableFromStruct(sm)
		if len(table.Columns) != len(expect.Columns) {
			t.Errorf("NewTableFromStruct() got = %v, want %v", len(table.Columns), len(expect.Columns))
		}
	}
}
