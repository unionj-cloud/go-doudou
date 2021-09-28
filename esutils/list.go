package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
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
		hits := make(chan *elastic.SearchHit)
		g, ctx := errgroup.WithContext(context.Background())
		g.Go(func() error {
			defer close(hits)
			scroll := es.client.Scroll().Index(es.esIndex).Type(es.esType).Query(boolQuery).Size(1000).KeepAlive("1m")
			for {
				results, err := scroll.Do(ctx)
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return errors.Wrap(err, "call Scroll() error")
				}
				for _, hit := range results.Hits.Hits {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						hits <- hit
					}
				}
			}
			return nil
		})

		c := make(chan interface{})
		for i := 0; i < 10; i++ {
			g.Go(func() error {
				for hit := range hits {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						var ret interface{}
						if callback == nil {
							var p map[string]interface{}
							json.Unmarshal(*hit.Source, &p)
							p["_id"] = hit.Id
							ret = p
						} else {
							if ret, err = callback(*hit.Source); err != nil {
								return errors.Wrap(err, "call callback() error")
							}
						}
						c <- ret
					}
				}
				return nil
			})
		}

		go func() {
			g.Wait()
			close(c)
		}()

		for s := range c {
			rets = append(rets, s)
		}

		if err := g.Wait(); err != nil {
			return nil, errors.Wrap(err, "call Wait() error")
		}
	} else {
		var searchResult *elastic.SearchResult
		ss := es.client.Search().Index(es.esIndex).Type(es.esType).Query(boolQuery)
		if paging.Sortby != nil && len(paging.Sortby) > 0 {
			for _, v := range paging.Sortby {
				ss = ss.Sort(v.Field, v.Ascending)
			}
		}
		if searchResult, err = ss.From(paging.Skip).Size(paging.Limit).Do(ctx); err != nil {
			return nil, errors.Wrap(err, "call Search() error")
		}
		for _, hit := range searchResult.Hits.Hits {
			var ret interface{}
			if callback == nil {
				var p map[string]interface{}
				json.Unmarshal(*hit.Source, &p)
				p["_id"] = hit.Id
				ret = p
			} else {
				if ret, err = callback(*hit.Source); err != nil {
					return nil, errors.Wrap(err, "call callback() error")
				}
			}

			rets = append(rets, ret)
		}
	}

	return rets, err
}
