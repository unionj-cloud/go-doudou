package gorm

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/wubin1989/gorm"
	"github.com/wubin1989/postgres"
)

type Parameter struct {
	Page    int64  `json:"page" form:"page"`
	Size    int64  `json:"size" form:"size"`
	Sort    string `json:"sort" form:"sort"`
	Order   string `json:"order" form:"order"`
	Fields  string `json:"fields" form:"fields"`
	Filters string `json:"filters" form:"filters"`
}

func (receiver Parameter) GetPage() int64 {
	return receiver.Page
}

func (receiver Parameter) GetSize() int64 {
	return receiver.Size
}

func (receiver Parameter) GetSort() string {
	return receiver.Sort
}

func (receiver Parameter) GetOrder() string {
	return receiver.Order
}

func (receiver Parameter) GetFields() string {
	return receiver.Fields
}

func (receiver Parameter) GetFilters() interface{} {
	return receiver.Filters
}

func (receiver Parameter) IParameterInstance() {

}

func Test_resContext_BuildWhereClause(t *testing.T) {
	pg := New(&Config{
		FieldSelectorEnabled: true,
	})
	filters := make([][]interface{}, 0)
	filters = append(filters, []interface{}{
		"sys_role_id",
		"=",
		123,
	})
	filters = append(filters, []interface{}{
		"and",
	})
	filters = append(filters, []interface{}{
		"deleted_at",
		"is",
		"null",
	})
	filterContent, _ := json.Marshal(filters)
	mockDB, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	resCxt := pg.With(db).Request(Parameter{
		Page:    0,
		Size:    0,
		Sort:    "",
		Order:   "",
		Fields:  "",
		Filters: string(filterContent),
	})
	statement, args := resCxt.BuildWhereClause()
	if resCxt.Error() != nil {
		panic(resCxt.Error())
	}
	fmt.Println(statement)
	fmt.Println(args)
}
