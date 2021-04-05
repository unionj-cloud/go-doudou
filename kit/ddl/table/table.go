package table

import (
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/kit/astutils"
	"github.com/unionj-cloud/go-doudou/kit/ddl/columnenum"
	"github.com/unionj-cloud/go-doudou/kit/ddl/sortenum"
	"github.com/unionj-cloud/go-doudou/kit/pathutils"
	"github.com/unionj-cloud/go-doudou/kit/stringutils"
	"github.com/unionj-cloud/go-doudou/kit/templateutils"
	"sort"
	"strconv"
	"strings"
)

type Extra string

const (
	update Extra = "on update CURRENT_TIMESTAMP"
)

const (
	now = "CURRENT_TIMESTAMP"
)

type IndexItems []IndexItem

type IndexItem struct {
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

type Column struct {
	Table         string
	Name          string
	Type          columnenum.ColumnType
	Default       interface{}
	Pk            bool
	Nullable      bool
	Unsigned      bool
	Autoincrement bool
	Extra         Extra
	Meta          astutils.FieldMeta
	AutoSet       bool
}

func (c *Column) ChangeColumnSql() (string, error) {
	return templateutils.StringBlock(pathutils.Abs("alter.tmpl"), "change", c)
}

func (c *Column) AddColumnSql() (string, error) {
	return templateutils.StringBlock(pathutils.Abs("alter.tmpl"), "add", c)
}

type Key string

const (
	pri Key = "PRI"
	uni Key = "UNI"
	mul Key = "MUL"
)

type Null string

const (
	yes Null = "YES"
	no  Null = "NO"
)

type DbColumn struct {
	Field   string  `db:"Field"`
	Type    string  `db:"Type"`
	Null    Null    `db:"Null"`
	Key     *Key    `db:"Key"`
	Default *string `db:"Default"`
	Extra   *string `db:"Extra"`
}

type DbIndex struct {
	Table string `db:"Table"`
	Non_unique `db:"Non_unique"`
	Key_name `db:"Field"`
	Seq_in_index `db:"Field"`
	Column_name `db:"Field"`
	Collation `db:"Field"`
	Cardinality `db:"Field"`
	Sub_part `db:"Field"`
	Packed `db:"Field"`
	Index_type `db:"Field"`
	Comment `db:"Field"`
	Index_comment `db:"Field"`
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
	table = strcase.ToSnake(structMeta.Name) + "s"
	for _, field := range structMeta.Fields {
		var (
			columnName    string
			columnType    columnenum.ColumnType
			columnDefault interface{}
			nullable      bool
			unsigned      bool
			autoincrement bool
			extra         Extra
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
							extra = Extra(value)
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

		if stringutils.IsEmpty(string(columnType)) {
			var trimmedType string
			if trimmedType = strings.TrimPrefix(field.Type, "*"); len(trimmedType) < len(field.Type) {
				nullable = true
			}
			columnType = toColumnType(trimmedType)
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

func (t *Table) CreateSql() (string, error) {
	return templateutils.String(pathutils.Abs("create.tmpl"), t)
}
