package template

const SvcImplImport = `
	"context"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/errorx"
	"{{.ConfigPackage}}"
	"{{.DtoPackage}}"
	"{{.ModelPackage}}"
	"{{.QueryPackage}}"
	"github.com/jinzhu/copier"
	"github.com/unionj-cloud/go-doudou/v2/framework/database"
	paginate "github.com/unionj-cloud/go-doudou/v2/toolkit/pagination/gorm"
`
