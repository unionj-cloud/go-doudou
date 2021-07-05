package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// if paging is nil, randomly return 10 pcs of documents as default
func (es *Es) Random(ctx context.Context, paging *Paging) ([]map[string]interface{}, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
		sr        *elastic.SearchResult
		rets      []map[string]interface{}
	)
	if paging == nil {
		paging = &Paging{
			Limit: 10,
		}
	}
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	fsq := elastic.NewFunctionScoreQuery().Query(boolQuery).AddScoreFunc(elastic.NewScriptFunction(elastic.NewScriptInline("Math.random()")))
	if sr, err = es.client.Search().Index(es.esIndex).Type(es.esType).Query(fsq).From(paging.Skip).Size(paging.Limit).Do(ctx); err != nil {
		return nil, errors.Wrap(err, "call Search() error")
	}
	for _, hit := range sr.Hits.Hits {
		var ret map[string]interface{}
		json.Unmarshal(*hit.Source, &ret)
		rets = append(rets, ret)
	}
	return rets, nil
}
