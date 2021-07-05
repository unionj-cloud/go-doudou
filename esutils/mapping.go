package esutils

import (
	"context"
	"github.com/Jeffail/gabs/v2"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

type MappingPayload struct {
	Base
	Fields []Field `json:"fields"`
}

func NewMapping(mp MappingPayload) string {
	var (
		mapping    *gabs.Container
		properties *gabs.Container
	)

	mapping = gabs.New()
	mapping.SetP("60s", "settings.refresh_interval")
	mapping.SetP("1", "settings.number_of_replicas")
	mapping.SetP("15", "settings.number_of_shards")

	properties = gabs.New()
	for _, f := range mp.Fields {
		properties.Set(f.Type, f.Name, "type")
	}

	mapping.Set(properties, "mappings", mp.Type, "properties")

	return mapping.String()
}

func (es *Es) PutMapping(ctx context.Context, mp MappingPayload) error {
	var (
		mapping    *gabs.Container
		properties *gabs.Container
		res        *elastic.PutMappingResponse
		err        error
	)
	mapping = gabs.New()
	properties = gabs.New()
	for _, f := range mp.Fields {
		properties.Set(f.Type, f.Name, "type")
	}
	mapping.Set(properties, "properties")
	if res, err = es.client.PutMapping().Index(mp.Index).Type(mp.Type).BodyString(mapping.String()).Do(ctx); err != nil {
		return errors.Wrap(err, "call PutMapping() error")
	}
	if !res.Acknowledged {
		return errors.New("putmapping failed!!!")
	}
	return nil
}

func (es *Es) CheckTypeExists(ctx context.Context) (b bool, err error) {
	if b, err = es.client.TypeExists().Index(es.esIndex).Type(es.esType).Do(ctx); err != nil {
		return false, errors.Wrap(err, "call TypeExists() error")
	}
	return
}
