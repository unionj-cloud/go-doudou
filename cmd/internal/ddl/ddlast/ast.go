package ddlast

import (
	"sort"
	"strings"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sliceutils"
)

type sortableFieldMeta []astutils.FieldMeta

// Len return length of sortableFieldMeta
func (it sortableFieldMeta) Len() int {
	return len(it)
}

// Less define asc or desc order
func (it sortableFieldMeta) Less(i, j int) bool {
	return it[i].Name < it[j].Name
}

// Swap change position of elements at i and j
func (it sortableFieldMeta) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
}

// FlatEmbed flat embed struct
func FlatEmbed(structs []astutils.StructMeta) []astutils.StructMeta {
	structMap := make(map[string]astutils.StructMeta)
	for _, structMeta := range structs {
		if _, exists := structMap[structMeta.Name]; !exists {
			structMap[structMeta.Name] = structMeta
		}
	}
	var result []astutils.StructMeta
	for _, structMeta := range structs {
		if sliceutils.IsEmpty(structMeta.Comments) {
			continue
		}
		if !strings.Contains(structMeta.Comments[0], "dd:table") {
			continue
		}
		_structMeta := astutils.StructMeta{
			Name:     structMeta.Name,
			Fields:   make([]astutils.FieldMeta, 0),
			Comments: make([]string, len(structMeta.Comments)),
		}
		copy(_structMeta.Comments, structMeta.Comments)

		fieldMap := make(map[string]astutils.FieldMeta)
		embedFieldMap := make(map[string]astutils.FieldMeta)
		for _, fieldMeta := range structMeta.Fields {
			if strings.HasPrefix(fieldMeta.Type, "embed") {
				if embedded, exists := structMap[fieldMeta.Name]; exists {
					for _, field := range embedded.Fields {
						embedFieldMap[field.Name] = field
					}
				}
			} else {
				_structMeta.Fields = append(_structMeta.Fields, fieldMeta)
				fieldMap[fieldMeta.Name] = fieldMeta
			}
		}

		var embedFields []astutils.FieldMeta
		for key, field := range embedFieldMap {
			if _, exists := fieldMap[key]; !exists {
				embedFields = append(embedFields, field)
			}
		}

		sort.Stable(sortableFieldMeta(embedFields))
		_structMeta.Fields = append(_structMeta.Fields, embedFields...)
		result = append(result, _structMeta)
	}

	return result
}
