package esutils

import (
	"context"
	"encoding/json"
	"testing"
)

func TestList(t *testing.T) {
	es := setupSubTest("test_list")
	type args struct {
		paging   *Paging
		esIndex  string
		esType   string
		callback func(message json.RawMessage) (interface{}, error)
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
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
				callback: nil,
			},
			want:    "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
			wantErr: false,
		},
		{
			args: args{
				paging:   nil,
				callback: nil,
			},
			want:    "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := es.List(context.Background(), tt.args.paging, tt.args.callback)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Error("got's length shouldn't be zero")
				return
			}
		})
	}
}
