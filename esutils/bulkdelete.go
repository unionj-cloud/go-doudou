package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
)

func (es *Es) BulkDelete(ctx context.Context, ids []string) error {
	bulkRequest := es.client.Bulk().Index(es.esIndex).Type(es.esType)

	for _, id := range ids {
		bulkRequest.Add(elastic.NewBulkDeleteRequest().Index(es.esIndex).Type(es.esType).Id(id))
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
		err = tracerr.New("bulk partially failed")
		return err
	}

	es.client.Flush(es.esIndex).Do(ctx)

	return nil
}
