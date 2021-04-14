package es

import (
	"context"
	"encoding/json"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func Random(paging *Paging, esIndex string, esType string) ([]interface{}, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
		src       interface{}
		data      []byte
		sr        *elastic.SearchResult
		rets      []interface{}
	)
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

	ctx := context.Background()

	//_, err = G_EsClient.Refresh().Index(esIndex).Do(ctx)
	//if err != nil {
	//	err = tracerr.Wrap(err)
	//	return 0, err
	//}
	//
	//_, err = G_EsClient.Flush().Index(esIndex).Do(ctx)
	//if err != nil {
	//	err = tracerr.Wrap(err)
	//	return 0, err
	//}

	fsq := elastic.NewFunctionScoreQuery().Query(boolQuery).AddScoreFunc(elastic.NewScriptFunction(elastic.NewScriptInline("Math.random()")))

	if sr, err = G_EsClient.Search().Index(esIndex).Type(esType).Query(fsq).From(paging.Skip).Size(paging.Limit).Do(ctx); err != nil {
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
