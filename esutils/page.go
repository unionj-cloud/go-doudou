package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

type PageResult struct {
	Page        int           `json:"page"` // from 1
	PageSize    int           `json:"page_size"`
	Total       int           `json:"total"`
	Docs        []interface{} `json:"docs"`
	HasNextPage bool          `json:"has_next_page"`
}

func (es *Es) Page(ctx context.Context, paging *Paging) (PageResult, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
		pr        PageResult
	)
	if paging == nil {
		paging = &Paging{
			Limit: -1,
		}
	}
	if paging.Limit < 0 || paging.Limit > 10000 {
		docs, err := es.List(ctx, paging, nil)
		if err != nil {
			return pr, errors.Wrap(err, "call List() error")
		}
		pr.Total = len(docs)
		pr.Docs = docs
		return pr, nil
	}
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	var rets []interface{}

	var searchResult *elastic.SearchResult
	ss := es.client.Search().Index(es.esIndex).Type(es.esType).Query(boolQuery)
	if paging.Sortby != nil && len(paging.Sortby) > 0 {
		for _, v := range paging.Sortby {
			ss = ss.Sort(v.Field, v.Ascending)
		}
	}
	if searchResult, err = ss.From(paging.Skip).Size(paging.Limit).Do(ctx); err != nil {
		return pr, errors.Wrap(err, "call Search() error")
	}
	for _, hit := range searchResult.Hits.Hits {
		var p map[string]interface{}
		json.Unmarshal(*hit.Source, &p)
		p["_id"] = hit.Id
		rets = append(rets, p)
	}

	pr.Docs = rets
	pr.Total = int(searchResult.TotalHits())
	pr.PageSize = paging.Limit
	if paging.Limit > 0 {
		pr.Page = paging.Skip/paging.Limit + 1
	}

	var totalPage int
	if pr.PageSize > 0 {
		if pr.Total%pr.PageSize > 0 {
			totalPage = pr.Total/pr.PageSize + 1
		} else {
			totalPage = pr.Total / pr.PageSize
		}
	}

	if pr.Page < totalPage {
		pr.HasNextPage = true
	}
	return pr, err
}
