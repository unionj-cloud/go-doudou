package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/copier"
)

// aggr only accept map[string]interface{} or elastic.Aggregation
func (es *Es) Stat(ctx context.Context, paging *Paging, aggr interface{}) (map[string]interface{}, error) {
	var (
		err          error
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
			_query, _ := src.Source()
			statQueryMap["query"] = _query
		}
		if aggr != nil {
			statQueryMap["aggs"] = aggr
		}
		if len(statQueryMap) == 0 {
			return nil, nil
		}
		if sr, err = searchService.Source(statQueryMap).Do(ctx); err != nil {
			return nil, errors.Wrap(err, "call Search() error")
		}
	case elastic.Aggregation:
		if src != nil {
			searchService = searchService.Query(src)
		}
		if sr, err = searchService.Aggregation("volume", raw).Do(ctx); err != nil {
			return nil, errors.Wrap(err, "call Search() error")
		}
	}
	var result map[string]interface{}
	copier.DeepCopy(sr.Aggregations, &result)
	return result, nil
}
