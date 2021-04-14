package es

import (
	"context"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
)

func BulkDelete(esindex, estype string, ids []string) error {
	bulkRequest := G_EsClient.Bulk().Index(esindex).Type(estype)

	for _, id := range ids {
		bulkRequest.Add(elastic.NewBulkDeleteRequest().Index(esindex).Type(estype).Id(id))
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
