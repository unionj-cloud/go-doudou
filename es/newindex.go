package es

import (
	"context"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
)

func NewIndex(index string, mapping string) (exists bool, err error) {
	lock.Lock()
	defer lock.Unlock()
	var (
		res *elastic.IndicesCreateResult
	)
	if exists, err = G_EsClient.IndexExists(index).Do(context.Background()); err != nil {
		err = tracerr.Wrap(err)
		return
	}
	if !exists {
		// Create a new index.
		if res, err = G_EsClient.CreateIndex(index).BodyString(mapping).Do(context.Background()); err != nil {
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
