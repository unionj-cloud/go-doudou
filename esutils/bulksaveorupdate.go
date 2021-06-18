package esutils

import (
	"context"
	"errors"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
)

func (es *Es) BulkSaveOrUpdate(ctx context.Context, docs []map[string]interface{}) error {
	bulkRequest := es.client.Bulk().Index(es.esIndex).Type(es.esType)

	for _, doc := range docs {
		bulkIndexRequest := elastic.NewBulkIndexRequest().Index(es.esIndex).Type(es.esType)
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

	if bulkRes, err = bulkRequest.Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return err
	}
	if bulkRes.Errors {
		for _, item := range bulkRes.Items {
			if item["index"].Error != nil {
				err = tracerr.Wrap(errors.New(item["index"].Error.Reason))
				return err
			}
		}
	}

	es.client.Flush(es.esIndex).Do(ctx)

	return nil
}
