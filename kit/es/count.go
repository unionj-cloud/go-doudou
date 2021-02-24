package es

import (
	"context"
	"encoding/json"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func Count(paging *Paging, esIndex string, esType string) (int64, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
		src       interface{}
		data      []byte
	)
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	if src, err = boolQuery.Source(); err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}
	if data, err = json.Marshal(src); err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}
	log.Println(string(data))

	ctx := context.Background()

	_, err = G_EsClient.Refresh().Index(esIndex).Do(ctx)
	if err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}

	_, err = G_EsClient.Flush().Index(esIndex).Do(ctx)
	if err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}

	var total int64
	if total, err = G_EsClient.Count().Index(esIndex).Type(esType).Query(boolQuery).Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}

	return total, nil
}
