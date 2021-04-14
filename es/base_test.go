package es

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"
	"github.com/unionj-cloud/go-doudou/constants"
	"github.com/unionj-cloud/go-doudou/testutils"
)

func Test_query(t *testing.T) {
	testutils.SkipCI(t)
	const index = "team3_voice_analysis_wb"

	teardownSubTest := SetupSubTest(index, t)
	defer teardownSubTest(t)

	data1 := "2020-06-01"
	data2 := "2020-06-20"
	data3 := "2020-07-10"

	createAt1, _ := time.Parse(constants.FORMAT2, data1)
	createAt2, _ := time.Parse(constants.FORMAT2, data2)
	createAt3, _ := time.Parse(constants.FORMAT2, data3)

	err := BulkSaveOrUpdate(index, index, []map[string]interface{}{
		{
			"createAt": createAt1.UTC().Format(constants.FORMATES),
			"text":     "2020年7月8日11时25分，高考文科综合/理科综合科目考试将要结束时，平顶山市一中考点一考生突然情绪失控，先后抓其右边、后边考生答题卡，造成两位考生答题卡损毁。",
		},
		{
			"createAt": createAt2.UTC().Format(constants.FORMATES),
			"text":     "考场两位监考教师及时制止，并稳定了考场秩序，市一中考点按程序启用备用答题卡，按规定补足答题卡被损毁的两位考生耽误的考试时间，两位考生将损毁卡的内容誊写在新答题卡上。",
		},
		{
			"createAt": createAt3.UTC().Format(constants.FORMATES),
			"text":     "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
		},
	})
	if err != nil {
		panic(err)
	}

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
				startDate: data1,
				endDate:   data3,
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
			var got string
			if src, err := bq.Source(); err != nil {
				panic(err)
			} else {
				if data, err := json.Marshal(src); err != nil {
					panic(err)
				} else {
					log.Println(string(data))
					got = string(data)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_range_query(t *testing.T) {
	const index = "team3_voice_analysis_dev"

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
			var got string
			if src, err := bq.Source(); err != nil {
				panic(err)
			} else {
				if data, err := json.Marshal(src); err != nil {
					panic(err)
				} else {
					log.Println(string(data))
					got = string(data)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}
