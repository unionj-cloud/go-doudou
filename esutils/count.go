package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
	"log"
)

func (es *Es) Count(ctx context.Context, paging *Paging) (int64, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
		src       interface{}
		data      []byte
	)
	if paging == nil {
		paging = &Paging{}
	}
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

	_, err = es.client.Refresh().Index(es.esIndex).Do(ctx)
	if err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}

	_, err = es.client.Flush().Index(es.esIndex).Do(ctx)
	if err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}

	var total int64
	if total, err = es.client.Count().Index(es.esIndex).Type(es.esType).Query(boolQuery).Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return 0, err
	}

	return total, nil
}
