package table

import (
	"cloud/unionj/papilio/kit/astutils"
	"cloud/unionj/papilio/kit/stringutils"
	"github.com/iancoleman/strcase"
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

type Index struct {
	Name  string
	Order int
	Sort  string
}

type ColumnType string

const (
	bitType        ColumnType = "BIT"
	textType       ColumnType = "TEXT"
	blobType       ColumnType = "BLOB"
	dateType       ColumnType = "DATE"
	datetimeType   ColumnType = "DATETIME"
	decimalType    ColumnType = "DECIMAL"
	doubleType     ColumnType = "DOUBLE"
	enumType       ColumnType = "ENUM"
	floatType      ColumnType = "FLOAT"
	geometryType   ColumnType = "GEOMETRY"
	mediumintType  ColumnType = "MEDIUMINT"
	jsonType       ColumnType = "JSON"
	intType        ColumnType = "INT"
	longtextType   ColumnType = "LONGTEXT"
	longblobType   ColumnType = "LONGBLOB"
	bigintType     ColumnType = "BIGINT"
	mediumtextType ColumnType = "MEDIUMTEXT"
	mediumblobType ColumnType = "MEDIUMBLOB"
	smallintType   ColumnType = "SMALLINT"
	tinyintType    ColumnType = "TINYINT"
	varcharType    ColumnType = "VARCHAR(255)"
)

func toColumnType(goType string) ColumnType {
	switch goType {
	case "int":
		return intType
	case "int64":
		return bigintType
	case "float32":
		return floatType
	case "float64":
		return doubleType
	case "string":
		return varcharType
	case "bool":
		return tinyintType
	case "time.Time":
		return datetimeType
	}
	panic("no available type")
}

type Column struct {
	Name          string
	Type          ColumnType
	Default       interface{}
	Pk            bool
	Nullable      bool
	Unsigned      bool
	Autoincrement bool
	Extra         Extra
}

type Table struct {
	Name          string
	Columns       []Column
	Pk            string
	UniqueIndexes []Index
	Indexes       []Index
}

func NewTableFromStruct(structMeta astutils.StructMeta) Table {
	var (
		columns       []Column
		uniqueindexes []Index
		indexes       []Index
		pkColumn      Column
	)
	for _, field := range structMeta.Fields {
		var (
			columnName    string
			columnType    ColumnType
			columnDefault interface{}
			nullable      bool
			unsigned      bool
			autoincrement bool
			extra         Extra
			uniqueindex   Index
			index         Index
			pk            bool
		)
		if stringutils.IsNotEmpty(field.Tag) {
			tags := strings.Split(field.Tag, " ")
			var papiTag string
			for _, tag := range tags {
				if trimedTag := strings.TrimPrefix(tag, "papi:"); len(trimedTag) < len(tag) {
					papiTag = strings.Trim(trimedTag, `"`)
					break
				}
			}
			if stringutils.IsNotEmpty(papiTag) {
				kvs := strings.Split(papiTag, ";")
				for _, kv := range kvs {
					pair := strings.Split(kv, ":")
					if len(pair) > 1 {
						key := pair[0]
						value := pair[1]
						switch key {
						case "type":
							columnType = ColumnType(value)
							break
						case "default":
							columnDefault = value
							break
						case "column":
							columnName = value
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
							sort := props[2]
							index = Index{
								Name:  indexName,
								Order: orderInt,
								Sort:  sort,
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
							sort := props[2]
							uniqueindex = Index{
								Name:  indexName,
								Order: orderInt,
								Sort:  sort,
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
								Name:  strcase.ToSnake(field.Name) + "_idx",
								Order: 1,
								Sort:  "asc",
							}
							break
						case "unique":
							uniqueindex = Index{
								Name:  strcase.ToSnake(field.Name) + "_idx",
								Order: 1,
								Sort:  "asc",
							}
							break
						}
					}
				}
			}
		}

		if stringutils.IsEmpty(columnName) {
			columnName = strcase.ToSnake(field.Name)
		}

		if stringutils.IsEmpty(string(columnType)) {
			var trimmedType string
			if trimmedType = strings.TrimPrefix(field.Type, "*"); len(trimmedType) < len(field.Type) {
				nullable = true
			}
			columnType = toColumnType(trimmedType)
		}

		columns = append(columns, Column{
			Name:          columnName,
			Type:          columnType,
			Default:       columnDefault,
			Nullable:      nullable,
			Unsigned:      unsigned,
			Autoincrement: autoincrement,
			Extra:         extra,
			Pk:            pk,
		})

		uniqueindexes = append(uniqueindexes, uniqueindex)
		indexes = append(indexes, index)
	}

	for _, column := range columns {
		if column.Pk {
			pkColumn = column
			break
		}
	}

	

	return Table{
		Name:          strcase.ToSnake(structMeta.Name) + "s",
		Columns:       columns,
		Pk:            pkColumn.Name,
		UniqueIndexes: nil,
		Indexes:       nil,
	}
}
