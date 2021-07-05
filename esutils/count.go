package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

func (es *Es) Count(ctx context.Context, paging *Paging) (int64, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
	)
	if paging == nil {
		paging = &Paging{}
	}
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	if _, err = es.client.Refresh().Index(es.esIndex).Do(ctx); err != nil {
		return 0, errors.Wrap(err, "call Refresh() error")
	}
	if _, err = es.client.Flush().Index(es.esIndex).Do(ctx); err != nil {
		return 0, errors.Wrap(err, "call Flush() error")
	}
	var total int64
	if total, err = es.client.Count().Index(es.esIndex).Type(es.esType).Query(boolQuery).Do(ctx); err != nil {
		return 0, errors.Wrap(err, "call Count() error")
	}
	return total, nil
}
