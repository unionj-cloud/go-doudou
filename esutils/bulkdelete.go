package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// BulkDelete delete es docs specified by ids in bulk
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
		return errors.Wrap(err, "call Bulk() error")
	}
	if bulkRes.Errors {
		return errors.New("bulk partially failed")
	}

	es.client.Flush(es.esIndex).Do(ctx)

	return nil
}
