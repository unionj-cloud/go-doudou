package esutils

import (
	"context"
	"encoding/json"
	"github.com/Jeffail/gabs/v2"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEs_Stat(t *testing.T) {
	es := setupSubTest("test_stat")

	jaggr := `{
        "groupBy": {
            "terms": {
                "field": "type.keyword",
                "size": 9999,
                "execution_hint": "map",
                "min_doc_count": 1
            }
        }
    }`
	aggr := make(map[string]interface{})
	json.Unmarshal([]byte(jaggr), &aggr)
	ret, _ := es.Stat(context.Background(), nil, aggr)
	expectj := `[{"doc_count":1,"key":"culture"},{"doc_count":1,"key":"education"},{"doc_count":1,"key":"sport"}]`
	var expect interface{}
	json.Unmarshal([]byte(expectj), &expect)
	assert.ElementsMatch(t, expect, gabs.Wrap(ret).Path("groupBy.buckets").Data())
}

func TestEs_Stat2(t *testing.T) {
	es := setupSubTest("test_stat1")

	aggr := elastic.NewTermsAggregation().Field("type.keyword").Size(9999).ExecutionHint("map").MinDocCount(1)
	ret, _ := es.Stat(context.Background(), nil, aggr)
	expectj := `[{"doc_count":1,"key":"culture"},{"doc_count":1,"key":"education"},{"doc_count":1,"key":"sport"}]`
	var expect interface{}
	json.Unmarshal([]byte(expectj), &expect)
	assert.ElementsMatch(t, expect, gabs.Wrap(ret).Path("volume.buckets").Data())
}
