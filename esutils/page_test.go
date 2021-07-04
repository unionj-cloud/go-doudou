package esutils

import (
	"context"
	"reflect"
	"testing"
)

func TestPage(t *testing.T) {
	es := setupSubTest("test_page")

	type args struct {
		paging *Paging
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
					StartDate: "2020-06-01",
					EndDate:   "2020-07-10",
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
			},
			want: PageResult{
				Page:     1,
				PageSize: 1,
				Total:    3,
				Docs: []interface{}{
					map[string]interface{}{
						"_id":      "9seTXHoBNx091WJ2QCh7",
						"id":       "9seTXHoBNx091WJ2QCh7",
						"text":     "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
						"type":     "culture",
						"createAt": "2020-07-10T00:00:00Z",
					},
				},
				HasNextPage: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := es.Page(context.Background(), tt.args.paging)
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
