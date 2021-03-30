package es

import (
	"reflect"
	"testing"
	"time"
	"github.com/unionj-cloud/go-doudou/kit/constants"
)

func TestPage(t *testing.T) {
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
		paging  *Paging
		esIndex string
		esType  string
	}
	tests := []struct {
		name    string
		args    args
		want    PageResult
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				paging: &Paging{
					StartDate: data1,
					EndDate:   data3,
					DateField: "createAt",
					Skip:      0,
					Limit:     1,
					Sortby: []Sort{
						{
							Field:     "createAt",
							Ascending: false,
						},
					},
					QueryConds: []QueryCond{
						{
							Pair: map[string][]interface{}{
								"text": {"考生"},
							},
							QueryLogic: SHOULD,
							QueryType:  MATCHPHRASE,
						},
					},
				},
				esIndex: index,
				esType:  index,
			},
			want: PageResult{
				Page:     1,
				PageSize: 1,
				Total:    3,
				Docs: []interface{}{
					map[string]interface{}{
						"text":     "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
						"createAt": "2020-07-10T00:00:00Z",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Page(tt.args.paging, tt.args.esIndex, tt.args.esType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Page() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Page() = %+v, want %v", got, tt.want)
			}
		})
	}
}
