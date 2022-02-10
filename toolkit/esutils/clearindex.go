package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// ClearIndex remove all docs
func (es *Es) ClearIndex(ctx context.Context) error {
	var (
		err error
		res *elastic.BulkIndexByScrollResponse
	)

	if res, err = es.client.DeleteByQuery(es.esIndex).Query(elastic.NewMatchAllQuery()).Do(ctx); err != nil {
		return errors.Wrap(err, "call DeleteByQuery() error")
	}
	if len(res.Failures) > 0 {
		return errors.New("failed to clear index " + es.esIndex)
	}
	return nil
}
