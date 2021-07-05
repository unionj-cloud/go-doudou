package esutils

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
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
			name: "",
			args: args{
				startDate: "2020-06-01",
				endDate:   "2020-07-10",
				dateField: "createAt",
				queryConds: []QueryCond{
					{
						Pair: map[string][]interface{}{
							"text":    {"考生"},
							"school":  {"西安理工+西安交大"},
							"address": {"北京+-西安"},
							"company": {"-unionj"},
						},
						QueryLogic: SHOULD,
						QueryType:  MATCHPHRASE,
					},
					{
						Pair: map[string][]interface{}{
							"text": {"高考"},
						},
						QueryLogic: MUST,
						QueryType:  MATCHPHRASE,
					},
					{
						Pair: map[string][]interface{}{
							"text": {"北京高考"},
						},
						QueryLogic: MUSTNOT,
						QueryType:  MATCHPHRASE,
					},
					{
						Pair: map[string][]interface{}{
							"content":      {"北京"},
							"content_full": {"unionj"},
						},
						QueryLogic: MUST,
						QueryType:  TERMS,
					},
				},
			},
			want: `{"bool":{"minimum_should_match":"1","must":[{"range":{"createAt":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-10"}}},{"bool":{"should":{"match_phrase":{"text":{"query":"高考"}}}}},{"terms":{"content":["北京"]}},{"terms":{"content_full":["unionj"]}}],"must_not":{"bool":{"should":{"match_phrase":{"text":{"query":"北京高考"}}}}},"should":[{"bool":{"should":{"match_phrase":{"text":{"query":"考生"}}}}},{"bool":{"should":{"bool":{"must":[{"match_phrase":{"school":{"query":"西安理工"}}},{"match_phrase":{"school":{"query":"西安交大"}}}]}}}},{"bool":{"should":{"bool":{"must":{"match_phrase":{"address":{"query":"北京"}}},"must_not":{"match_phrase":{"address":{"query":"西安"}}}}}}},{"bool":{"should":{"bool":{"must_not":{"match_phrase":{"company":{"query":"unionj"}}}}}}}]}}`,
		},
		{
			name: "",
			args: args{
				startDate: "2020-06-01",
				endDate:   "2020-07-10",
				dateField: "createAt",
				queryConds: []QueryCond{
					{
						Pair: map[string][]interface{}{
							"type.keyword": {"education"},
							"status":       {float64(200)},
						},
						QueryLogic: MUST,
						QueryType:  TERMS,
					},
					{
						Pair: map[string][]interface{}{
							"dept.keyword": {"unionj*"},
						},
						QueryLogic: SHOULD,
						QueryType:  WILDCARD,
					},
					{
						Pair: map[string][]interface{}{
							"position.keyword": {"dev*"},
						},
						QueryLogic: MUST,
						QueryType:  WILDCARD,
					},
					{
						Pair: map[string][]interface{}{
							"city.keyword": {"四川*"},
						},
						QueryLogic: MUSTNOT,
						QueryType:  WILDCARD,
					},
					{
						Pair: map[string][]interface{}{
							"project.keyword": {"unionj"},
						},
						QueryLogic: SHOULD,
						QueryType:  PREFIX,
					},
					{
						Pair: map[string][]interface{}{
							"name.keyword": {"unionj"},
						},
						QueryLogic: MUSTNOT,
						QueryType:  PREFIX,
					},
					{
						Pair: map[string][]interface{}{
							"book.keyword": {"go"},
						},
						QueryLogic: MUST,
						QueryType:  PREFIX,
					},
				},
			},
			want: `{"bool":{"minimum_should_match":"1","must":[{"range":{"createAt":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-10"}}},{"terms":{"type.keyword":["education"]}},{"terms":{"status":[200]}},{"wildcard":{"position.keyword":{"wildcard":"dev*"}}},{"prefix":{"book.keyword":"go"}}],"must_not":[{"wildcard":{"city.keyword":{"wildcard":"四川*"}}},{"prefix":{"name.keyword":"unionj"}}],"should":[{"wildcard":{"dept.keyword":{"wildcard":"unionj*"}}},{"prefix":{"project.keyword":"unionj"}}]}}`,
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
			want, _ := gabs.ParseJSON([]byte(tt.want))
			_src := gabs.Wrap(src)
			fmt.Println(_src.String())
			if !assert.ElementsMatch(t, _src.Path("bool.must").Data(), want.Path("bool.must").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.must").Data(), want.Path("bool.must").Data())
			}
			if !assert.ElementsMatch(t, _src.Path("bool.should").Data(), want.Path("bool.should").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.should").Data(), want.Path("bool.should").Data())
			}
		})
	}
}

func Test_range_query(t *testing.T) {
	param1 := make(map[string]interface{})
	param1["to"] = 0.4
	param1["include_upper"] = true

	param2 := make(map[string]interface{})
	param2["from"] = 0.6
	param2["include_lower"] = true

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
							"senseResult": {param1},
						},
						QueryLogic: MUST,
						QueryType:  RANGE,
					},
					{
						Pair: map[string][]interface{}{
							"visitSenseResult": {param2},
						},
						QueryLogic: SHOULD,
						QueryType:  RANGE,
					},
					{
						Pair: map[string][]interface{}{
							"commonSenseResult": {param2},
						},
						QueryLogic: MUSTNOT,
						QueryType:  RANGE,
					},
					{
						Pair: map[string][]interface{}{
							"orderPhrase": {float64(300)},
						},
						QueryLogic: MUST,
						QueryType:  TERMS,
					},
				},
			},
			want: `{"bool":{"minimum_should_match":"1","must":[{"range":{"acceptDate":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-01"}}},{"range":{"senseResult":{"from":null,"include_lower":true,"include_upper":true,"to":0.4}}},{"terms":{"orderPhrase":[300]}}],"must_not":{"range":{"commonSenseResult":{"from":0.6,"include_lower":true,"include_upper":true,"to":null}}},"should":{"range":{"visitSenseResult":{"from":0.6,"include_lower":true,"include_upper":true,"to":null}}}}}`,
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
			want, _ := gabs.ParseJSON([]byte(tt.want))
			_src := gabs.Wrap(src)
			fmt.Println(_src.String())
			if !assert.ElementsMatch(t, _src.Path("bool.must").Data(), want.Path("bool.must").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.must").Data(), want.Path("bool.must").Data())
			}
			if !assert.Equal(t, _src.Path("bool.should").Data(), want.Path("bool.should").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.should").Data(), want.Path("bool.should").Data())
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
						},
						QueryLogic: MUSTNOT,
						QueryType:  EXISTS,
					},
					{
						Pair: map[string][]interface{}{
							"flag": {},
						},
						QueryLogic: MUST,
						QueryType:  EXISTS,
					},
					{
						Pair: map[string][]interface{}{
							"status": {},
						},
						QueryLogic: SHOULD,
						QueryType:  EXISTS,
					},
					{
						Pair: map[string][]interface{}{
							"orderPhrase": {float64(300)},
						},
						QueryLogic: MUST,
						QueryType:  TERMS,
					},
				},
			},
			want: `{"bool":{"minimum_should_match":"1","must":[{"range":{"acceptDate":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-01"}}},{"exists":{"field":"flag"}},{"terms":{"orderPhrase":[300]}}],"must_not":{"exists":{"field":"delete_at"}},"should":{"exists":{"field":"status"}}}}`,
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
			want, _ := gabs.ParseJSON([]byte(tt.want))
			_src := gabs.Wrap(src)
			fmt.Println(_src.String())
			if !assert.ElementsMatch(t, _src.Path("bool.must").Data(), want.Path("bool.must").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.must").Data(), want.Path("bool.must").Data())
			}
			if !assert.Equal(t, _src.Path("bool.should").Data(), want.Path("bool.should").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.should").Data(), want.Path("bool.should").Data())
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
							"status": {float64(100), float64(300)},
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
									"price": {float64(0)},
								},
								QueryLogic: MUST,
								QueryType:  TERMS,
							},
						},
					},
				},
			},
			want: `{"bool":{"must":{"range":{"acceptDate":{"format":"yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis","from":"2020-06-01","include_lower":true,"include_upper":true,"time_zone":"Asia/Shanghai","to":"2020-07-01"}}},"must_not":[{"exists":{"field":"delete_at"}},{"terms":{"status":[100,300]}},{"bool":{"must":[{"terms":{"type":["网络调查"]}},{"terms":{"price":[0]}}]}}]}}`,
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
			want, _ := gabs.ParseJSON([]byte(tt.want))
			_src := gabs.Wrap(src)
			fmt.Println(_src.String())
			if !assert.Equal(t, _src.Path("bool.must").Data(), want.Path("bool.must").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.must").Data(), want.Path("bool.must").Data())
			}
			if !assert.ElementsMatch(t, _src.Path("bool.must_not").Data(), want.Path("bool.must_not").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.must_not").Data(), want.Path("bool.must_not").Data())
			}
			if !assert.ElementsMatch(t, _src.Path("bool.should").Data(), want.Path("bool.should").Data()) {
				t.Errorf("query() = %v, want %v", _src.Path("bool.should").Data(), want.Path("bool.should").Data())
			}
		})
	}
}

func TestNewEs(t *testing.T) {
	url := "http://test.com"
	username := "unionj"
	password := "unionj"
	client, err := elastic.NewSimpleClient(
		elastic.SetErrorLog(logrus.New()),
		elastic.SetURL(url),
		elastic.SetBasicAuth(username, password),
		elastic.SetGzip(true),
	)
	if err != nil {
		panic(fmt.Errorf("NewSimpleClient() error: %+v\n", err))
	}
	type args struct {
		esIndex string
		esType  string
		opts    []EsOption
	}
	tests := []struct {
		name string
		args args
		want *Es
	}{
		{
			name: "",
			args: args{
				esIndex: "test1",
				opts: []EsOption{
					WithUsername(username),
					WithPassword(password),
					WithUrls([]string{url}),
				},
			},
			want: nil,
		},
		{
			name: "",
			args: args{
				esIndex: "test2",
				opts: []EsOption{
					WithClient(client),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		assert.NotPanics(t, func() {
			got := NewEs(tt.args.esIndex, tt.args.esType, tt.args.opts...)
			got.SetType(got.esIndex)
		})
	}
}

func TestEs_newDefaultClient(t *testing.T) {
	assert.Panics(t, func() {
		NewEs("test", "test", WithUrls([]string{"wrongurl"}))
	})
}
