package esutils

import (
	"context"
	"encoding/json"
	"github.com/ztrue/tracerr"
)

func (es *Es) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	if getResult, err := es.client.Get().Index(es.esIndex).Type(es.esType).Id(id).Do(ctx); err != nil {
		err = tracerr.Wrap(err)
		return nil, err
	} else {
		var p map[string]interface{}
		if err := json.Unmarshal(*getResult.Source, &p); err != nil {
			err = tracerr.Wrap(err)
			return nil, err
		}
		p["_id"] = getResult.Id
		return p, nil
	}
}
