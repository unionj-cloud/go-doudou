package es

import (
	"context"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
)

func BulkSaveOrUpdate(esindex, estype string, docs []map[string]interface{}) error {
	bulkRequest := G_EsClient.Bulk().Index(esindex).Type(estype)

	for _, doc := range docs {
		bulkIndexRequest := elastic.NewBulkIndexRequest().Index(esindex).Type(estype)
		if id, ok := doc["id"]; ok {
			bulkIndexRequest = bulkIndexRequest.Id(id.(string))
		}
		bulkIndexRequest = bulkIndexRequest.Doc(doc)
		bulkRequest.Add(bulkIndexRequest)
	}

	var (
		bulkRes *elastic.BulkResponse
		err     error
	)

	if bulkRes, err = bulkRequest.Do(context.TODO()); err != nil {
		err = tracerr.Wrap(err)
		return err
	}
	if bulkRes.Errors {
		err = tracerr.New("bulk partially failed")
		return err
	}

	G_EsClient.Flush(esindex).Do(context.TODO())

	return nil
}
