package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
	"log"
)

// if paging is nil, randomly return 10 pcs of documents as default
func (es *Es) Random(ctx context.Context, paging *Paging) ([]map[string]interface{}, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
		src       interface{}
		data      []byte
		sr        *elastic.SearchResult
		rets      []map[string]interface{}
	)
	if paging == nil {
		paging = &Paging{
			Limit: 10,
		}
	}
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	if src, err = boolQuery.Source(); err != nil {
		err = tracerr.Wrap(err)
		return nil, err
	}
	if data, err = json.Marshal(src); err != nil {
		err = tracerr.Wrap(err)
		return nil, err
	}
	log.Println(string(data))

	fsq := elastic.NewFunctionScoreQuery().Query(boolQuery).AddScoreFunc(elastic.NewScriptFunction(elastic.NewScriptInline("Math.random()")))

	if sr, err = es.client.Search().Index(es.esIndex).Type(es.esType).Query(fsq).From(paging.Skip).Size(paging.Limit).Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return nil, err
	}

	for _, hit := range sr.Hits.Hits {
		var ret map[string]interface{}
		if err = json.Unmarshal(*hit.Source, &ret); err != nil {
			err = tracerr.Wrap(err)
			return nil, err
		}
		rets = append(rets, ret)
	}

	return rets, nil
}
