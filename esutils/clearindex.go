package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
)

func (es *Es) ClearIndex(ctx context.Context) error {
	var (
		err error
		res *elastic.BulkIndexByScrollResponse
	)

	if res, err = es.client.DeleteByQuery(es.esIndex).Query(elastic.NewMatchAllQuery()).Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return err
	}
	if len(res.Failures) > 0 {
		err = tracerr.New("failed to clear index " + es.esIndex)
		return err
	}
	return nil
}
