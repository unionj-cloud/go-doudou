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
	var rets []interface{}
	if paging.Limit < 0 || paging.Limit > 10000 {
		if rets, err = es.fetchAll(boolQuery, callback); err != nil {
			return nil, errors.Wrap(err, "call es.fetchAll error")
		}
	} else {
		if rets, err = es.doPaging(ctx, paging, boolQuery, callback); err != nil {
			return nil, errors.Wrap(err, "call es.fetchAll error")
		}
	}
	return rets, nil
}
