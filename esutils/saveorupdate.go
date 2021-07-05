package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/stringutils"
)

func (es *Es) SaveOrUpdate(ctx context.Context, doc interface{}) (string, error) {
	var (
		indexRes *elastic.IndexResponse
		err      error
	)

	indexRequest := es.client.Index().Index(es.esIndex).Type(es.esType)

	id, err := getId(doc)
	if err != nil {
		return "", errors.Wrap(err, "method SaveOrUpdate() error")
	}
	if stringutils.IsNotEmpty(id) {
		indexRequest = indexRequest.Id(id)
	}

	if indexRes, err = indexRequest.BodyJson(&doc).Do(ctx); err != nil {
		return "", errors.Wrap(err, "call Index() error")
	}

	es.client.Flush(es.esIndex).Do(ctx)

	return indexRes.Id, nil
}
