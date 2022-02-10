package table

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/ddlast"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/extraenum"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/keyenum"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/nullenum"
	"github.com/unionj-cloud/go-doudou/cmd/internal/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"testing"
)

func ExampleNewTableFromStruct() {
	testDir := pathutils.Abs("../testdata/domain")
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
		tab := NewTableFromStruct(sm)
		fmt.Println(len(tab.Indexes))
		var statement string
		if statement, err = tab.CreateSql(); err != nil {
			panic(err)
		}
		fmt.Println(statement)
	}

	// Output:
	//0
	//CREATE TABLE `order` (
	//`id` INT NOT NULL AUTO_INCREMENT,
	//`amount` BIGINT NOT NULL,
	//`user_id` int NOT NULL,
	//`create_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
	//`delete_at` DATETIME NULL,
	//`update_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	//PRIMARY KEY (`id`),
	//CONSTRAINT `fk_ddl_user` FOREIGN KEY (`user_id`)
	//REFERENCES `ddl_user`(`id`)
	//ON DELETE CASCADE ON UPDATE NO ACTION)
	//5
	//CREATE TABLE `user` (
	//`id` INT NOT NULL AUTO_INCREMENT,
	//`name` VARCHAR(255) NOT NULL DEFAULT 'jack',
	//`phone` VARCHAR(255) NOT NULL DEFAULT '13552053960' comment '手机号',
	//`age` INT NOT NULL,
	//`no` int NOT NULL,
	//`unique_col` int NOT NULL,
	//`unique_col_2` int NOT NULL,
	//`school` VARCHAR(255) NULL DEFAULT 'harvard' comment '学校',
	//`is_student` TINYINT NOT NULL,
	//`rule` varchar(255) NOT NULL comment '链接匹配规则，匹配的链接采用该css规则来爬',
	//`rule_type` varchar(45) NOT NULL comment '链接匹配规则类型，支持prefix前缀匹配和regex正则匹配',
	//`arrive_at` datetime NULL comment '到货时间',
	//`status` tinyint(4) NOT NULL comment '0进行中
	//1完结
	//2取消',
	//`create_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
	//`delete_at` DATETIME NULL,
	//`update_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	//PRIMARY KEY (`id`),
	//INDEX `age_idx` (`age` asc),
	//INDEX `name_phone_idx` (`phone` asc,`name` asc),
	//UNIQUE INDEX `no_idx` (`no` asc),
	//UNIQUE INDEX `rule_idx` (`rule` asc),
	//UNIQUE INDEX `unique_col_idx` (`unique_col` asc,`unique_col_2` asc))
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
						Default:       "",
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
						Default:       "",
						Pk:            false,
						Nullable:      true,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "",
					},
					{
						Name:          "no",
						Type:          columnenum.IntType,
						Default:       "",
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
			want:    "CREATE TABLE `users` (\n`id` INT NOT NULL AUTO_INCREMENT,\n`name` VARCHAR(255) NULL DEFAULT 'wubin',\n`phone` VARCHAR(255) NULL DEFAULT '13552053960' comment '手机号',\n`age` INT NULL,\n`no` INT NOT NULL,\nPRIMARY KEY (`id`),\nINDEX `name_phone_idx` (`name` asc,`phone` desc),\nUNIQUE INDEX `uni_no` (`no` asc))",
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
		Default       string
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
		Default       string
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
		}, {
			name: "8",
			args: args{
				goType: "decimal.Decimal",
			},
			want: "decimal(6,2)",
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
					Default:       "harvard",
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
				Tag:      `dd:"type:VARCHAR(255);default:'harvard';extra:comment '学校'"`,
				Comments: nil,
			},
		}, {
			name: "2",
			args: args{
				col: Column{
					Table:         "users",
					Name:          "favourite",
					Type:          columnenum.VarcharType,
					Default:       "current_timestamp",
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
				Tag:      `dd:"auto;type:VARCHAR(255);default:current_timestamp;extra:comment '学校';index:my_index,1,asc"`,
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
	type args struct {
		defaultVal string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				defaultVal: "CURRENT_TIMESTAMP",
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

func TestNewIndexFromDbIndexes(t *testing.T) {
	type args struct {
		dbIndexes []DbIndex
	}
	tests := []struct {
		name string
		args args
		want Index
	}{
		{
			name: "",
			args: args{
				dbIndexes: []DbIndex{
					{
						Table:      "ddl_user",
						NonUnique:  false,
						KeyName:    "age_idx",
						SeqInIndex: 1,
						ColumnName: "age",
						Collation:  "A",
					},
				},
			},
			want: Index{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIndexFromDbIndexes(tt.args.dbIndexes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndexFromDbIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_DropIndexSql(t *testing.T) {
	type fields struct {
		Table  string
		Unique bool
		Name   string
		Items  []IndexItem
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
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
			want:    "ALTER TABLE `ddl_user` DROP INDEX `age_idx`;",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := &Index{
				Table:  tt.fields.Table,
				Unique: tt.fields.Unique,
				Name:   tt.fields.Name,
				Items:  tt.fields.Items,
			}
			got, err := idx.DropIndexSql()
			if (err != nil) != tt.wantErr {
				t.Errorf("DropIndexSql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DropIndexSql() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_AddIndexSql(t *testing.T) {
	type fields struct {
		Table  string
		Unique bool
		Name   string
		Items  []IndexItem
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
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
			want:    "ALTER TABLE `ddl_user` ADD UNIQUE INDEX `age_idx` (`age` asc);",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := &Index{
				Table:  tt.fields.Table,
				Unique: tt.fields.Unique,
				Name:   tt.fields.Name,
				Items:  tt.fields.Items,
			}
			got, err := idx.AddIndexSql()
			if (err != nil) != tt.wantErr {
				t.Errorf("AddIndexSql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddIndexSql() got = %v, want %v", got, tt.want)
			}
		})
	}
}
