package es

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStat(t *testing.T) {

	const esindex = "weibo_articles_and_weiboers";
	const estype = "weibo_articles_and_weiboer";
	const aggrJson = `{
    "volume": {
      "date_histogram": {
        "field": "published_at",
        "format": "M",
        "interval": "month",
        "time_zone": "+08:00",
        "min_doc_count": 0,
        "extended_bounds": {
          "min": 1590940800000,
          "max": 1594310400000
        }
      },
      "aggs": {
        "groupBy": {
          "terms": {
            "field": "weiboer_id",
            "size": 10,
            "execution_hint": "map"
          },
          "aggs": {
            "document_count": {
              "sum": {
                "script": "1"
              }
            }
          }
        }
      }
    }
  }`
	var aggrMap map[string]interface{}
	if err := json.Unmarshal([]byte(aggrJson), &aggrMap); err != nil {
		panic(err)
	}

	type args struct {
		paging  *Paging
		esIndex string
		esType  string
		aggr    interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				paging: &Paging{
					StartDate: "2020-06-01",
					EndDate:   "2020-07-10",
					DateField: "published_at",
					Skip:      0,
					Limit:     1,
					Sortby: []Sort{
						{
							Field:     "published_at",
							Ascending: false,
						},
					},
					QueryConds: []QueryCond{
						{
							Pair: map[string][]interface{}{
								"content_full": {"美国暴乱"},
							},
							QueryLogic: SHOULD,
							QueryType:  MATCHPHRASE,
						},
					},
				},
				esIndex: esindex,
				esType:  estype,
				aggr:    aggrMap,
			},
			want:    `{"volume":{"buckets":[{"key_as_string":"6","key":1590940800000,"doc_count":33403,"groupBy":{"doc_count_error_upper_bound":98,"sum_other_doc_count":31131,"buckets":[{"key":"1650926392","doc_count":481,"document_count":{"value":481.0}},{"key":"2133442851","doc_count":364,"document_count":{"value":364.0}},{"key":"3365355414","doc_count":287,"document_count":{"value":287.0}},{"key":"1746159404","doc_count":193,"document_count":{"value":193.0}},{"key":"2587796532","doc_count":191,"document_count":{"value":191.0}},{"key":"1659893422","doc_count":162,"document_count":{"value":162.0}},{"key":"1770493244","doc_count":156,"document_count":{"value":156.0}},{"key":"1724256551","doc_count":151,"document_count":{"value":151.0}},{"key":"5855598209","doc_count":146,"document_count":{"value":146.0}},{"key":"2197546765","doc_count":141,"document_count":{"value":141.0}}]}},{"key_as_string":"7","key":1593532800000,"doc_count":563,"groupBy":{"doc_count_error_upper_bound":0,"sum_other_doc_count":410,"buckets":[{"key":"3365355414","doc_count":40,"document_count":{"value":40.0}},{"key":"1724256551","doc_count":28,"document_count":{"value":28.0}},{"key":"1393202693","doc_count":22,"document_count":{"value":22.0}},{"key":"2743112592","doc_count":13,"document_count":{"value":13.0}},{"key":"1650926392","doc_count":12,"document_count":{"value":12.0}},{"key":"1746159404","doc_count":11,"document_count":{"value":11.0}},{"key":"5394248433","doc_count":9,"document_count":{"value":9.0}},{"key":"1652966954","doc_count":7,"document_count":{"value":7.0}},{"key":"1751910164","doc_count":6,"document_count":{"value":6.0}},{"key":"1657967434","doc_count":5,"document_count":{"value":5.0}}]}}]}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err := Stat(tt.args.paging, tt.args.esIndex, tt.args.esType, tt.args.aggr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			bb, err := json.Marshal(ret)
			if err != nil {
				panic(err)
			}
			got := string(bb)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stat() = %v, want %v", got, tt.want)
			}
		})
	}
}
