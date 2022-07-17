package codegen

import (
	"github.com/goccy/go-json"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenDomainGo(t *testing.T) {
	dir := pathutils.Abs("../testdata/testdomain")
	metaJSON := `{"Name":"User","Fields":[{"Name":"Id","Type":"int","Tag":"dd:\"pk;auto;type:int\"","Comments":null},{"Name":"Name","Type":"string","Tag":"dd:\"type:varchar(255);default:'jack';index:name_phone_idx,2,asc\"","Comments":null},{"Name":"Phone","Type":"string","Tag":"dd:\"type:varchar(255);default:'13552053960';extra:comment 'mobile phone';index:name_phone_idx,1,asc\"","Comments":null},{"Name":"Age","Type":"int","Tag":"dd:\"type:int;index:age_idx,1,asc\"","Comments":null},{"Name":"No","Type":"int","Tag":"dd:\"type:int;unique:no_idx,1,asc\"","Comments":null},{"Name":"School","Type":"*string","Tag":"dd:\"type:varchar(255);default:'harvard';extra:comment 'school'\"","Comments":null},{"Name":"IsStudent","Type":"bool","Tag":"dd:\"type:tinyint\"","Comments":null},{"Name":"CreateAt","Type":"*time.Time","Tag":"dd:\"type:datetime;default:CURRENT_TIMESTAMP\"","Comments":null},{"Name":"UpdateAt","Type":"*time.Time","Tag":"dd:\"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP\"","Comments":null},{"Name":"DeleteAt","Type":"*time.Time","Tag":"dd:\"type:datetime\"","Comments":null}],"Comments":null,"Methods":null}`
	var meta astutils.StructMeta
	json.Unmarshal([]byte(metaJSON), &meta)
	type args struct {
		domainpath string
		meta       astutils.StructMeta
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				domainpath: dir,
				meta:       meta,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenDomainGo(tt.args.domainpath, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("GenDomainGo() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer os.RemoveAll(dir)
			expect := `package domain

import (
	"time"
)

//dd:table
type User struct {
	Id        int        ` + "`" + `dd:"pk;auto;type:int"` + "`" + `
	Name      string     ` + "`" + `dd:"type:varchar(255);default:'jack';index:name_phone_idx,2,asc"` + "`" + `
	Phone     string     ` + "`" + `dd:"type:varchar(255);default:'13552053960';extra:comment 'mobile phone';index:name_phone_idx,1,asc"` + "`" + `
	Age       int        ` + "`" + `dd:"type:int;index:age_idx,1,asc"` + "`" + `
	No        int        ` + "`" + `dd:"type:int;unique:no_idx,1,asc"` + "`" + `
	School    *string    ` + "`" + `dd:"type:varchar(255);default:'harvard';extra:comment 'school'"` + "`" + `
	IsStudent bool       ` + "`" + `dd:"type:tinyint"` + "`" + `
	CreateAt  *time.Time ` + "`" + `dd:"type:datetime;default:CURRENT_TIMESTAMP"` + "`" + `
	UpdateAt  *time.Time ` + "`" + `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"` + "`" + `
	DeleteAt  *time.Time ` + "`" + `dd:"type:datetime"` + "`" + `
}
`
			domainfile := filepath.Join(dir, "user.go")
			f, err := os.Open(domainfile)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			if string(content) != expect {
				t.Errorf("want %s, got %s\n", expect, string(content))
			}
		})
	}
}
