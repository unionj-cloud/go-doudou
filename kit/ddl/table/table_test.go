package table

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

const testDir = "/Users/wubin1989/workspace/cloud/go-doudou/kit/ddl/example/models"

func ExampleNewTableFromStruct() {
	var files []string
	var err error
	err = filepath.Walk(testDir, astutils.Visit(&files))
	if err != nil {
		panic(err)
	}
	var sc astutils.StructCollector
	for _, file := range files {
		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		ast.Walk(&sc, root)
	}
	flattened := sc.FlatEmbed()

	for _, sm := range flattened {
		table := NewTableFromStruct(sm)
		fmt.Println(table)
	}
	//Output:
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
						Type:          intType,
						Default:       nil,
						Pk:            true,
						Nullable:      false,
						Unsigned:      false,
						Autoincrement: true,
						Extra:         "",
					},
					{
						Name:          "name",
						Type:          varcharType,
						Default:       "wubin",
						Pk:            false,
						Nullable:      true,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "",
					},
					{
						Name:          "phone",
						Type:          varcharType,
						Default:       "13552053960",
						Pk:            false,
						Nullable:      true,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "",
					},
					{
						Name:          "age",
						Type:          intType,
						Default:       nil,
						Pk:            false,
						Nullable:      true,
						Unsigned:      false,
						Autoincrement: false,
						Extra:         "",
					},
					{
						Name:          "no",
						Type:          intType,
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
			want:    "CREATE TABLE `users` (\n`id` INT NOT NULL AUTO_INCREMENT,\n`name` VARCHAR(255) NULL ,\n`phone` VARCHAR(255) NULL ,\n`age` INT NULL ,\n`no` INT NOT NULL ,\nPRIMARY KEY (`id`),\nINDEX `name_phone_idx` (`name` asc,`phone` desc),\nUNIQUE INDEX `uni_no` (`no` asc));",
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
