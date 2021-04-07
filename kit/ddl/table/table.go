package table

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/kit/ddl/extraenum"
	"github.com/unionj-cloud/go-doudou/kit/ddl/keyenum"
	"github.com/unionj-cloud/go-doudou/kit/ddl/nullenum"
	"github.com/unionj-cloud/go-doudou/kit/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"github.com/unionj-cloud/go-doudou/kit/stringutils"
	"github.com/unionj-cloud/go-doudou/kit/templateutils"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	now = "CURRENT_TIMESTAMP"
)

type IndexItems []IndexItem

type IndexItem struct {
	Unique bool
	Name   string
	Column string
	Order  int
	Sort   sortenum.Sort
}

func (it IndexItems) Len() int {
	return len(it)
}
func (it IndexItems) Less(i, j int) bool {
	return it[i].Order < it[j].Order
}
func (it IndexItems) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
}

type Index struct {
	Unique bool
	Name   string
	Items  []IndexItem
}

func toColumnType(goType string) columnenum.ColumnType {
	switch goType {
	case "int":
		return columnenum.IntType
	case "int64":
		return columnenum.BigintType
	case "float32":
		return columnenum.FloatType
	case "float64":
		return columnenum.DoubleType
	case "string":
		return columnenum.VarcharType
	case "bool":
		return columnenum.TinyintType
	case "time.Time":
		return columnenum.DatetimeType
	}
	panic("no available type")
}

func toGoType(colType columnenum.ColumnType, nullable bool) string {
	var goType string
	if nullable {
		goType += "*"
	}
	if strings.HasPrefix(string(colType), strings.ToLower(string(columnenum.IntType))) {
		goType += "int"
	} else if strings.HasPrefix(string(colType), strings.ToLower(string(columnenum.BigintType))) {
		goType += "int64"
	} else if strings.HasPrefix(string(colType), strings.ToLower(string(columnenum.FloatType))) {
		goType += "float32"
	} else if strings.HasPrefix(string(colType), strings.ToLower(string(columnenum.DoubleType))) {
		goType += "float64"
	} else if strings.HasPrefix(string(colType), strings.ToLower(string(columnenum.VarcharType))) {
		goType += "string"
	} else if strings.HasPrefix(string(colType), strings.ToLower(string(columnenum.TinyintType))) {
		goType += "bool"
	} else if strings.HasPrefix(string(colType), strings.ToLower(string(columnenum.DatetimeType))) {
		goType += "time.Time"
	} else {
		panic("no available type")
	}
	return goType
}

func CheckPk(key keyenum.Key) bool {
	return key == keyenum.Pri
}

func CheckNull(null nullenum.Null) bool {
	return null == nullenum.Yes
}

func CheckUnsigned(dbColType string) bool {
	splits := strings.Split(dbColType, " ")
	if len(splits) == 1 {
		return false
	}
	return splits[1] == "unsigned"
}

func CheckAutoincrement(extra string) bool {
	return strings.Contains(extra, "auto_increment")
}

func CheckAutoSet(defaultVal *string) bool {
	return defaultVal != nil && *defaultVal == now
}

type Column struct {
	Table         string
	Name          string
	Type          columnenum.ColumnType
	Default       interface{}
	Pk            bool
	Nullable      bool
	Unsigned      bool
	Autoincrement bool
	Extra         extraenum.Extra
	Meta          astutils.FieldMeta
	AutoSet       bool
	Indexes       []IndexItem
}

func (c *Column) ChangeColumnSql() (string, error) {
	return templateutils.StringBlock(pathutils.Abs("alter.tmpl"), "change", c)
}

func (c *Column) AddColumnSql() (string, error) {
	return templateutils.StringBlock(pathutils.Abs("alter.tmpl"), "add", c)
}

type DbColumn struct {
	Field   string        `db:"Field"`
	Type    string        `db:"Type"`
	Null    nullenum.Null `db:"Null"`
	Key     keyenum.Key   `db:"Key"`
	Default *string       `db:"Default"`
	Extra   string        `db:"Extra"`
	Comment string        `db:"Comment"`
}

// https://www.mysqltutorial.org/mysql-index/mysql-show-indexes/
type DbIndex struct {
	Table        string `db:"Table"`        // The name of the table
	Non_unique   bool   `db:"Non_unique"`   // 1 if the index can contain duplicates, 0 if it cannot.
	Key_name     string `db:"Key_name"`     // The name of the index. The primary key index always has the name of PRIMARY.
	Seq_in_index int    `db:"Seq_in_index"` // The column sequence number in the index. The first column sequence number starts from 1.
	Column_name  string `db:"Column_name"`  // The column name
	Collation    string `db:"Collation"`    // Collation represents how the column is sorted in the index. A means ascending, B means descending, or NULL means not sorted.
}

type Table struct {
	Name    string
	Columns []Column
	Pk      string
	Indexes []Index
	Meta    astutils.StructMeta
}

func NewTableFromStruct(structMeta astutils.StructMeta) Table {
	var (
		columns       []Column
		uniqueindexes []Index
		indexes       []Index
		pkColumn      Column
		table         string
	)
	table = strcase.ToSnake(structMeta.Name)
	for _, field := range structMeta.Fields {
		var (
			columnName    string
			columnType    columnenum.ColumnType
			columnDefault interface{}
			nullable      bool
			unsigned      bool
			autoincrement bool
			extra         extraenum.Extra
			uniqueindex   Index
			index         Index
			pk            bool
			autoSet       bool
		)
		columnName = strcase.ToSnake(field.Name)
		if stringutils.IsNotEmpty(field.Tag) {
			tags := strings.Split(field.Tag, `" `)
			var ddTag string
			for _, tag := range tags {
				if trimedTag := strings.TrimPrefix(tag, "dd:"); len(trimedTag) < len(tag) {
					ddTag = strings.Trim(trimedTag, `"`)
					break
				}
			}
			if stringutils.IsNotEmpty(ddTag) {
				kvs := strings.Split(ddTag, ";")
				for _, kv := range kvs {
					pair := strings.Split(kv, ":")
					if len(pair) > 1 {
						key := pair[0]
						value := pair[1]
						switch key {
						case "type":
							columnType = columnenum.ColumnType(value)
							break
						case "default":
							columnDefault = value
							if value == now {
								autoSet = true
							}
							break
						case "extra":
							extra = extraenum.Extra(value)
							break
						case "index":
							props := strings.Split(value, ",")
							indexName := props[0]
							order := props[1]
							orderInt, err := strconv.Atoi(order)
							if err != nil {
								panic(err)
							}
							var sort sortenum.Sort
							if len(props) < 3 || stringutils.IsEmpty(props[2]) {
								sort = sortenum.Asc
							} else {
								sort = sortenum.Sort(props[2])
							}
							index = Index{
								Name: indexName,
								Items: []IndexItem{
									{
										Order: orderInt,
										Sort:  sort,
									},
								},
							}
							break
						case "unique":
							props := strings.Split(value, ",")
							indexName := props[0]
							order := props[1]
							orderInt, err := strconv.Atoi(order)
							if err != nil {
								panic(err)
							}
							var sort sortenum.Sort
							if len(props) < 3 || stringutils.IsEmpty(props[2]) {
								sort = sortenum.Asc
							} else {
								sort = sortenum.Sort(props[2])
							}
							uniqueindex = Index{
								Name: indexName,
								Items: []IndexItem{
									{
										Order: orderInt,
										Sort:  sort,
									},
								},
							}
							break
						}
					} else {
						key := pair[0]
						switch key {
						case "pk":
							pk = true
							break
						case "null":
							nullable = true
							break
						case "unsigned":
							unsigned = true
							break
						case "auto":
							autoincrement = true
							break
						case "index":
							index = Index{
								Name: strcase.ToSnake(field.Name) + "_idx",
								Items: []IndexItem{
									{
										Order: 1,
										Sort:  sortenum.Asc,
									},
								},
							}
							break
						case "unique":
							uniqueindex = Index{
								Name: strcase.ToSnake(field.Name) + "_idx",
								Items: []IndexItem{
									{
										Order: 1,
										Sort:  sortenum.Asc,
									},
								},
							}
							break
						}
					}
				}
			}
		}

		if strings.HasPrefix(field.Type, "*") {
			nullable = true
		}

		if stringutils.IsEmpty(string(columnType)) {
			columnType = toColumnType(strings.TrimPrefix(field.Type, "*"))
		}

		if stringutils.IsNotEmpty(uniqueindex.Name) {
			uniqueindex.Items[0].Column = columnName
			uniqueindexes = append(uniqueindexes, uniqueindex)
		}

		if stringutils.IsNotEmpty(index.Name) {
			index.Items[0].Column = columnName
			indexes = append(indexes, index)
		}

		columns = append(columns, Column{
			Table:         table,
			Name:          columnName,
			Type:          columnType,
			Default:       columnDefault,
			Nullable:      nullable,
			Unsigned:      unsigned,
			Autoincrement: autoincrement,
			Extra:         extra,
			Pk:            pk,
			Meta:          field,
			AutoSet:       autoSet,
		})
	}

	for _, column := range columns {
		if column.Pk {
			pkColumn = column
			break
		}
	}

	uniqueMap := make(map[string][]IndexItem)
	indexMap := make(map[string][]IndexItem)

	for _, unique := range uniqueindexes {
		if items, exists := uniqueMap[unique.Name]; exists {
			items = append(items, unique.Items...)
			uniqueMap[unique.Name] = items
		} else {
			uniqueMap[unique.Name] = unique.Items
		}
	}

	for _, index := range indexes {
		if items, exists := indexMap[index.Name]; exists {
			items = append(items, index.Items...)
			indexMap[index.Name] = items
		} else {
			indexMap[index.Name] = index.Items
		}
	}

	var uniquesResult, indexesResult []Index

	for k, v := range uniqueMap {
		it := IndexItems(v)
		sort.Stable(it)
		uniquesResult = append(uniquesResult, Index{
			Unique: true,
			Name:   k,
			Items:  it,
		})
	}

	for k, v := range indexMap {
		it := IndexItems(v)
		sort.Stable(it)
		indexesResult = append(indexesResult, Index{
			Name:  k,
			Items: it,
		})
	}

	indexesResult = append(indexesResult, uniquesResult...)

	return Table{
		Name:    table,
		Columns: columns,
		Pk:      pkColumn.Name,
		Indexes: indexesResult,
		Meta:    structMeta,
	}
}

func NewFieldFromColumn(col Column) astutils.FieldMeta {
	tag := "dd:"
	var feats []string
	if col.Pk {
		feats = append(feats, "pk")
	}
	if col.Autoincrement {
		feats = append(feats, "auto")
	}
	goType := toGoType(col.Type, col.Nullable)
	if col.Nullable && !strings.HasPrefix(goType, "*") {
		feats = append(feats, "null")
	}
	if stringutils.IsNotEmpty(string(col.Type)) {
		feats = append(feats, fmt.Sprintf("type:%s", string(col.Type)))
	}
	if !reflect.ValueOf(col.Default).IsZero() {
		if ptr, ok := col.Default.(*string); ok {
			val := *ptr
			re := regexp.MustCompile(`^\(.+\)$`)
			var defaultClause string
			if val == "CURRENT_TIMESTAMP" || re.MatchString(val) {
				defaultClause = fmt.Sprintf("default:%s", val)
			} else {
				defaultClause = fmt.Sprintf("default:'%s'", val)
			}
			feats = append(feats, defaultClause)
		}
	}
	if stringutils.IsNotEmpty(string(col.Extra)) {
		feats = append(feats, fmt.Sprintf("extra:%s", string(col.Extra)))
	}
	for _, idx := range col.Indexes {
		var indexClause string
		if idx.Name == "PRIMARY" {
			continue
		}
		if idx.Unique {
			indexClause = "unique:"
		} else {
			indexClause = "index:"
		}
		indexClause += fmt.Sprintf("%s,%d,%s", idx.Name, idx.Order, string(idx.Sort))
		feats = append(feats, indexClause)
	}

	return astutils.FieldMeta{
		Name: strcase.ToCamel(col.Name),
		Type: goType,
		Tag:  fmt.Sprintf(`%s"%s"`, tag, strings.Join(feats, ";")),
	}
}

func (t *Table) CreateSql() (string, error) {
	return templateutils.String(pathutils.Abs("create.tmpl"), t)
}
