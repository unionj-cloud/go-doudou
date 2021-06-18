package esutils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_query(t *testing.T) {
	type args struct {
		startDate  string
		endDate    string
		dateField  string
		queryConds []QueryCond
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				startDate: "2020-06-01",
				endDate:   "2020-07-10",
				dateField: "createAt",
				queryConds: []QueryCond{
					{
						Pair: map[string][]interface{}{
							"text":    {"考生"},
							"school":  {"西安理工"},
							"address": {"北京"},
							"company": {"清研"},
						},
						QueryLogic: SHOULD,
						QueryType:  MATCHPHRASE,
					},
					{
						Pair: map[string][]interface{}{
							"content":      {"北京"},
							"content_full": {"清研"},
						},
						QueryLogic: MUST,
						QueryType:  TERMS,
					},
				},
			},
			want: `{"bool":{"minimum_should_match":"1","must":[{"range":{"createAt":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-10"}}},{"terms":{"content_full":["清研"]}},{"terms":{"content":["北京"]}}],"should":[{"bool":{"should":{"match_phrase":{"address":{"query":"北京"}}}}},{"bool":{"should":{"match_phrase":{"company":{"query":"清研"}}}}},{"bool":{"should":{"match_phrase":{"text":{"query":"考生"}}}}},{"bool":{"should":{"match_phrase":{"school":{"query":"西安理工"}}}}}]}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := query(tt.args.startDate, tt.args.endDate, tt.args.dateField, tt.args.queryConds)
			var src interface{}
			var err error
			if src, err = bq.Source(); err != nil {
				panic(err)
			}
			want := make(map[string]interface{})
			json.Unmarshal([]byte(tt.want), &want)
			_src := src.(map[string]interface{})
			if !assert.ElementsMatch(t, _src["must"], want["must"]) {
				t.Errorf("query() = %v, want %v", _src["must"], want["must"])
			}
			if !assert.ElementsMatch(t, _src["should"], want["should"]) {
				t.Errorf("query() = %v, want %v", _src["should"], want["should"])
			}
		})
	}
}

func Test_range_query(t *testing.T) {
	param1 := make(map[string]interface{})
	param1["to"] = 0.4
	param1["includeUpper"] = true

	param2 := make(map[string]interface{})
	param2["from"] = 0.6
	param2["includeLower"] = true

	type args struct {
		startDate  string
		endDate    string
		dateField  string
		queryConds []QueryCond
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				startDate: "2020-06-01",
				endDate:   "2020-07-01",
				dateField: "acceptDate",
				queryConds: []QueryCond{
					{
						Pair: map[string][]interface{}{
							"senseResult":      {param1},
							"visitSenseResult": {param2},
						},
						QueryLogic: MUST,
						QueryType:  RANGE,
					},
					{
						Pair: map[string][]interface{}{
							"orderPhrase": {300},
						},
						QueryLogic: MUST,
						QueryType:  TERMS,
					},
				},
			},
			want: `{"bool":{"must":[{"range":{"acceptDate":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-01"}}},{"range":{"senseResult":{"from":null,"include_lower":true,"include_upper":true,"to":0.4}}},{"range":{"visitSenseResult":{"from":0.6,"include_lower":true,"include_upper":true,"to":null}}},{"terms":{"orderPhrase":[300]}}]}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := query(tt.args.startDate, tt.args.endDate, tt.args.dateField, tt.args.queryConds)
			var src interface{}
			var err error
			if src, err = bq.Source(); err != nil {
				panic(err)
			}
			want := make(map[string]interface{})
			json.Unmarshal([]byte(tt.want), &want)
			_src := src.(map[string]interface{})
			if !assert.ElementsMatch(t, _src["must"], want["must"]) {
				t.Errorf("query() = %v, want %v", _src["must"], want["must"])
			}
			if !assert.ElementsMatch(t, _src["should"], want["should"]) {
				t.Errorf("query() = %v, want %v", _src["should"], want["should"])
			}
		})
	}
}

func Test_exists_query(t *testing.T) {
	type args struct {
		startDate  string
		endDate    string
		dateField  string
		queryConds []QueryCond
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				startDate: "2020-06-01",
				endDate:   "2020-07-01",
				dateField: "acceptDate",
				queryConds: []QueryCond{
					{
						Pair: map[string][]interface{}{
							"delete_at": {},
							"flag":      {},
						},
						QueryLogic: MUSTNOT,
						QueryType:  EXISTS,
					},
					{
						Pair: map[string][]interface{}{
							"orderPhrase": {300},
						},
						QueryLogic: MUST,
						QueryType:  TERMS,
					},
				},
			},
			want: `{"bool":{"must":[{"range":{"acceptDate":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-01"}}},{"terms":{"orderPhrase":[300]}}],"must_not":[{"exists":{"field":"delete_at"}},{"exists":{"field":"flag"}}]}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := query(tt.args.startDate, tt.args.endDate, tt.args.dateField, tt.args.queryConds)
			var src interface{}
			var err error
			if src, err = bq.Source(); err != nil {
				panic(err)
			}
			want := make(map[string]interface{})
			json.Unmarshal([]byte(tt.want), &want)
			_src := src.(map[string]interface{})
			if !assert.ElementsMatch(t, _src["must"], want["must"]) {
				t.Errorf("query() = %v, want %v", _src["must"], want["must"])
			}
			if !assert.ElementsMatch(t, _src["should"], want["should"]) {
				t.Errorf("query() = %v, want %v", _src["should"], want["should"])
			}
		})
	}
}

func Test_children_query(t *testing.T) {
	type args struct {
		startDate  string
		endDate    string
		dateField  string
		queryConds []QueryCond
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				startDate: "2020-06-01",
				endDate:   "2020-07-01",
				dateField: "acceptDate",
				queryConds: []QueryCond{
					{
						Pair: map[string][]interface{}{
							"delete_at": {},
						},
						QueryLogic: MUSTNOT,
						QueryType:  EXISTS,
					},
					{
						Pair: map[string][]interface{}{
							"status": {100, 300},
						},
						QueryLogic: MUSTNOT,
						QueryType:  TERMS,
					},
					{
						QueryLogic: MUSTNOT,
						Children: []QueryCond{
							{
								Pair: map[string][]interface{}{
									"type": {"网络调查"},
								},
								QueryLogic: MUST,
								QueryType:  TERMS,
							},
							{
								Pair: map[string][]interface{}{
									"price": {0},
								},
								QueryLogic: MUST,
								QueryType:  TERMS,
							},
						},
					},
				},
			},
			want: `{"bool":{"must":[{"range":{"acceptDate":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-01"}}},{"terms":{"orderPhrase":[300]}}],"must_not":[{"exists":{"field":"delete_at"}},{"exists":{"field":"flag"}}]}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := query(tt.args.startDate, tt.args.endDate, tt.args.dateField, tt.args.queryConds)
			var src interface{}
			var err error
			if src, err = bq.Source(); err != nil {
				panic(err)
			}
			want := make(map[string]interface{})
			json.Unmarshal([]byte(tt.want), &want)
			_src := src.(map[string]interface{})
			if !assert.ElementsMatch(t, _src["must"], want["must"]) {
				t.Errorf("query() = %v, want %v", _src["must"], want["must"])
			}
			if !assert.ElementsMatch(t, _src["should"], want["should"]) {
				t.Errorf("query() = %v, want %v", _src["should"], want["should"])
			}
		})
	}
}
