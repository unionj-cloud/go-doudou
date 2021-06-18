package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
	"log"
)

func (es *Es) Stat(ctx context.Context, paging *Paging, aggr interface{}) (interface{}, error) {
	var (
		err          error
		data         []byte
		sr           *elastic.SearchResult
		statQueryMap map[string]interface{}
		src          elastic.Query
	)

	if paging != nil {
		src = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	}

	searchService := es.client.Search().Index(es.esIndex).Type(es.esType)

	switch raw := aggr.(type) {
	case map[string]interface{}:
		statQueryMap = make(map[string]interface{})
		if src != nil {
			var query interface{}
			query, err = src.Source()
			if err != nil {
				return nil, err
			}
			statQueryMap["query"] = query
		}
		if aggr != nil {
			statQueryMap["aggs"] = aggr
		}
		if data, err = json.Marshal(statQueryMap); err != nil {
			err = tracerr.Wrap(err)
			return nil, err
		}
		log.Println(string(data))

		if len(statQueryMap) == 0 {
			return nil, nil
		}

		if sr, err = searchService.Source(statQueryMap).Do(ctx); err != nil {
			err = tracerr.Wrap(err)
			return nil, err
		}
	case elastic.Aggregation:
		if src != nil {
			searchService = searchService.Query(src)
		}
		if sr, err = searchService.Aggregation("volume", raw).Do(ctx); err != nil {
			err = tracerr.Wrap(err)
			return nil, err
		}
	}

	return sr.Aggregations, nil
}
