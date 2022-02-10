package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// List fetch docs by paging
func (es *Es) List(ctx context.Context, paging *Paging, callback func(message json.RawMessage) (interface{}, error)) ([]interface{}, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
	)
	if paging == nil {
		paging = &Paging{
			Limit: -1,
		}
	}
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	fsc := elastic.NewFetchSourceContext(true)
	if len(paging.Includes) > 0 {
		fsc = fsc.Include(paging.Includes...)
	}
	if len(paging.Excludes) > 0 {
		fsc = fsc.Exclude(paging.Excludes...)
	}
	var rets []interface{}
	if paging.Limit < 0 || paging.Limit > 10000 {
		if rets, err = es.fetchAll(fsc, boolQuery, callback); err != nil {
			return nil, errors.Wrap(err, "call es.fetchAll error")
		}
	} else {
		if rets, err = es.doPaging(ctx, fsc, paging, boolQuery, callback); err != nil {
			return nil, errors.Wrap(err, "call es.fetchAll error")
		}
	}
	return rets, nil
}
