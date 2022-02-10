package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// Count counts docs by paging
func (es *Es) Count(ctx context.Context, paging *Paging) (int64, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
	)
	if paging == nil {
		paging = &Paging{}
	}
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	es.client.Refresh().Index(es.esIndex).Do(ctx)
	es.client.Flush().Index(es.esIndex).Do(ctx)
	var total int64
	if total, err = es.client.Count().Index(es.esIndex).Type(es.esType).Query(boolQuery).Do(ctx); err != nil {
		return 0, errors.Wrap(err, "call Count() error")
	}
	return total, nil
}
