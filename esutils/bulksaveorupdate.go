package esutils

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"reflect"
)

func getId(doc interface{}) (string, error) {
	docVal := reflect.ValueOf(doc)
	var idVal reflect.Value
	if docVal.Kind() == reflect.Map {
		idVal = docVal.MapIndex(reflect.ValueOf("id"))
	} else if docVal.Kind() == reflect.Struct {
		idVal = docVal.FieldByName("Id")
	} else {
		return "", errors.New("method getId() error: single document must be map or struct type")
	}
	if idVal.IsValid() {
		if idVal.Kind() == reflect.String {
			return idVal.String(), nil
		}
		return fmt.Sprintf("%v", idVal.Interface()), nil
	}
	return "", nil
}

// BulkSaveOrUpdate save or update docs in bulk
func (es *Es) BulkSaveOrUpdate(ctx context.Context, docs []interface{}) error {
	bulkRequest := es.client.Bulk().Index(es.esIndex).Type(es.esType)

	for _, doc := range docs {
		id, err := getId(doc)
		if err != nil {
			return errors.Wrap(err, "method BulkSaveOrUpdate() error")
		}
		bulkIndexRequest := elastic.NewBulkIndexRequest().Index(es.esIndex).Type(es.esType)
		if stringutils.IsNotEmpty(id) {
			bulkIndexRequest = bulkIndexRequest.Id(id)
		}
		bulkIndexRequest = bulkIndexRequest.Doc(doc)
		bulkRequest.Add(bulkIndexRequest)
	}

	var (
		bulkRes *elastic.BulkResponse
		err     error
	)

	if bulkRes, err = bulkRequest.Do(ctx); err != nil {
		return errors.Wrap(err, "call Bulk() error")
	}
	if bulkRes.Errors {
		for _, item := range bulkRes.Items {
			if item["index"].Error != nil {
				return errors.New(item["index"].Error.Reason)
			}
		}
	}

	es.client.Flush(es.esIndex).Do(ctx)

	return nil
}
