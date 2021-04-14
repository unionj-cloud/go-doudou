package es

import (
	"context"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
)

func ClearIndex(index string) error {
	var (
		err error
		res *elastic.BulkIndexByScrollResponse
	)

	if res, err = G_EsClient.DeleteByQuery(index).Query(elastic.NewMatchAllQuery()).Do(context.Background()); err != nil {
		err = tracerr.Wrap(err)
		return err
	}
	if len(res.Failures) > 0 {
		err = tracerr.New("failed to clear index " + index)
		return err
	}
	return nil
}
