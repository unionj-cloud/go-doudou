package es

import (
	"context"
	"encoding/json"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

type PageResult struct {
	Page        int           `json:"page"` // from 1
	PageSize    int           `json:"page_size"`
	Total       int           `json:"total"`
	Docs        []interface{} `json:"docs"`
	HasNextPage bool          `json:"has_next_page"`
}

func Page(paging *Paging, esIndex string, esType string) (PageResult, error) {
	var pr PageResult
	if paging.Limit < 0 || paging.Limit > 10000 {
		docs, err := List(paging, esIndex, esType, nil)
		if err != nil {
			err = tracerr.Wrap(err)
			return pr, err
		}
		pr.Total = len(docs)
		pr.Docs = docs
		return pr, nil
	}

	var (
		err       error
		boolQuery *elastic.BoolQuery
		src       interface{}
		data      []byte
	)
	boolQuery = query(paging.StartDate, paging.EndDate, paging.DateField, paging.QueryConds)
	if src, err = boolQuery.Source(); err != nil {
		err = tracerr.Wrap(err)
		return pr, err
	}
	if data, err = json.Marshal(src); err != nil {
		err = tracerr.Wrap(err)
		return pr, err
	}
	log.Println(string(data))

	var rets []interface{}

	var searchResult *elastic.SearchResult
	ss := G_EsClient.Search().Index(esIndex).Type(esType).Query(boolQuery)
	if paging.Sortby != nil && len(paging.Sortby) > 0 {
		for _, v := range paging.Sortby {
			ss = ss.Sort(v.Field, v.Ascending)
		}
	}
	if searchResult, err = ss.From(paging.Skip).Size(paging.Limit).Do(context.Background()); err != nil {
		err = tracerr.Wrap(err)
		return pr, err
	}
	for _, hit := range searchResult.Hits.Hits {
		var p map[string]interface{}
		err := json.Unmarshal(*hit.Source, &p)
		if err != nil {
			err = tracerr.Wrap(err)
			return pr, err
		}
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
