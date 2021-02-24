package es

import (
	"context"
	"github.com/Jeffail/gabs/v2"
	"github.com/ztrue/tracerr"
	"gopkg.in/olivere/elastic.v5"
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

func PutMapping(mp MappingPayload) error {
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
	if res, err = G_EsClient.PutMapping().Index(mp.Index).Type(mp.Type).BodyString(mapping.String()).Do(context.Background()); err != nil {
		err = tracerr.Wrap(err)
		return err
	}
	if !res.Acknowledged {
		err = tracerr.New("putmapping failed!!!")
		return err
	}
	return nil
}

func CheckTypeExists(esindex, estype string) (b bool, err error) {
	if b, err = G_EsClient.TypeExists().Index(esindex).Type(estype).Do(context.Background()); err != nil {
		err = tracerr.Wrap(err)
	}
	return
}
