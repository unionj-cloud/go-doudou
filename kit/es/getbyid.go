package es

import (
	"context"
	"encoding/json"
	"github.com/ztrue/tracerr"
)

func GetByID(id string, esIndex string, esType string) (map[string]interface{}, error) {
	if getResult, err := G_EsClient.Get().Index(esIndex).Type(esType).Id(id).Do(context.Background()); err != nil {
		err = tracerr.Wrap(err)
		return nil, err
	} else {
		var p map[string]interface{}
		if err := json.Unmarshal(*getResult.Source, &p); err != nil {
			err = tracerr.Wrap(err)
			return nil, err
		}
		return p, nil
	}
}
