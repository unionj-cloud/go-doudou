package esutils

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// GetByID gets a doc by id
func (es *Es) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var (
		getResult *elastic.GetResult
		err       error
	)
	if getResult, err = es.client.Get().Index(es.esIndex).Type(es.esType).Id(id).Do(ctx); err != nil {
		return nil, errors.Wrap(err, "call Get() error")
	}
	var p map[string]interface{}
	json.Unmarshal(*getResult.Source, &p)
	p["_id"] = getResult.Id
	return p, nil
}
