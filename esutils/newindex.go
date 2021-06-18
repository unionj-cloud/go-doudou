package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
)

func (es *Es) NewIndex(ctx context.Context, mapping string) (exists bool, err error) {
	lock.Lock()
	defer lock.Unlock()
	var (
		res *elastic.IndicesCreateResult
	)
	if exists, err = es.client.IndexExists(es.esIndex).Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return
	}
	if !exists {
		// Create a new index.
		if res, err = es.client.CreateIndex(es.esIndex).BodyString(mapping).Do(ctx); err != nil {
			err = tracerr.Wrap(err)
			return
		}
		if !res.Acknowledged {
			err = tracerr.New("create index failed!!!")
			return
		}
	}
	return
}
