package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
)

func (es *Es) List(ctx context.Context, paging *Paging, callback func(message json.RawMessage) (interface{}, error)) ([]interface{}, error) {
	var (
		err       error
		boolQuery *elastic.BoolQuery
		src       interface{}
		data      []byte
	)
	if paging == nil {
		paging = &Paging{
			Limit: -1,
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

	//_, err = client.Refresh().Index(esIndex).Do(ctx)
	//if err != nil {
	//	err = tracerr.Wrap(err)
	//	return nil, err
	//}
	//
	//_, err = client.Flush().Index(esIndex).Do(ctx)
	//if err != nil {
	//	err = tracerr.Wrap(err)
	//	return nil, err
	//}

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
					err = tracerr.Wrap(err)
					return err
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
							err := json.Unmarshal(*hit.Source, &p)
							if err != nil {
								err = tracerr.Wrap(err)
								return err
							}
							p["_id"] = hit.Id
							ret = p
						} else {
							if ret, err = callback(*hit.Source); err != nil {
								err = tracerr.Wrap(err)
								return err
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
			err = tracerr.Wrap(err)
			return nil, err
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
			err = tracerr.Wrap(err)
			return nil, err
		}
		for _, hit := range searchResult.Hits.Hits {
			var ret interface{}
			if callback == nil {
				var p map[string]interface{}
				err := json.Unmarshal(*hit.Source, &p)
				if err != nil {
					err = tracerr.Wrap(err)
					return nil, err
				}
				p["_id"] = hit.Id
				ret = p
			} else {
				if ret, err = callback(*hit.Source); err != nil {
					err = tracerr.Wrap(err)
					return nil, err
				}
			}

			rets = append(rets, ret)
		}
	}

	return rets, err
}
