package esutils

import (
	"context"
	"testing"
)

func TestRandom(t *testing.T) {
	es := setupSubTest("test_random")
	type args struct {
		paging *Paging
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
					DateField: "createAt",
					Skip:      0,
					Limit:     1,
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
			want:    "目前，我办已将损毁其他考生答题卡的考生违规情况上报河南省招生办公室，将依规对该考生进行处理。平顶山市招生考试委员会办公室",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := es.Random(context.Background(), tt.args.paging)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Error("got's length shouldn't be zero")
				return
			}
			data := got[0]

			origin := data["text"]

			var current interface{}

			for i := 0; i < 10; i++ {
				got, err := es.Random(context.Background(), tt.args.paging)
				if (err != nil) != tt.wantErr {
					t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) == 0 {
					t.Error("got's length shouldn't be zero")
					return
				}
				data := got[0]

				current = data["text"]

				if current != origin {
					break
				}
			}

			if current == origin {
				t.Errorf("Random() = %v, want %v", data["text"], tt.want)
			}
		})
	}
}
