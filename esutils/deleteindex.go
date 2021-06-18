package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
)

func (es *Es) DeleteIndex(ctx context.Context) error {
	var (
		err error
		res *elastic.IndicesDeleteResponse
	)
	if res, err = es.client.DeleteIndex(es.esIndex).Do(ctx); err != nil {
		if elastic.IsNotFound(err) {
			return nil
		}
		err = tracerr.Wrap(err)
		return err
	}
	if !res.Acknowledged {
		err = tracerr.New("failed to delete index" + es.esIndex)
		return err
	}
	return nil
}
