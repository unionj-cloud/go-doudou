package es

import (
	"context"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
)

func DeleteIndex(index string) error {
	var (
		err error
		res *elastic.IndicesDeleteResponse
	)
	if res, err = G_EsClient.DeleteIndex(index).Do(context.Background()); err != nil {
		if elastic.IsNotFound(err) {
			return nil
		}
		err = tracerr.Wrap(err)
		return err
	}
	if !res.Acknowledged {
		err = tracerr.New("failed to delete index" + index)
		return err
	}
	return nil
}
