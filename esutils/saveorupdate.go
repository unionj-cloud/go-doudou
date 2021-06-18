package esutils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/ztrue/tracerr"
)

func (es *Es) SaveOrUpdate(ctx context.Context, doc map[string]interface{}) (string, error) {
	var (
		indexRes *elastic.IndexResponse
		err      error
	)

	indexRequest := es.client.Index().Index(es.esIndex).Type(es.esType)
	if id, exists := doc["id"]; exists {
		if idstr, ok := id.(string); ok {
			indexRequest = indexRequest.Id(idstr)
		}
	}

	if indexRes, err = indexRequest.BodyJson(&doc).Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return "", err
	}

	es.client.Flush(es.esIndex).Do(ctx)

	return indexRes.Id, nil
}
