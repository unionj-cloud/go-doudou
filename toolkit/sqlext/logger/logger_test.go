package logger

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
)

func TestSqlLogger_Log(t *testing.T) {
	type fields struct {
		logger logrus.StdLogger
	}
	type args struct {
		query string
		args  []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				logger: log.New(os.Stderr, "", log.LstdFlags),
			},
			args: args{
				query: "select * from ddl_user where (`school` = 'Shanxi Univ.' and `delete_at` is null) order by `create_at` desc limit ?,?",
				args:  []interface{}{0, 2},
			},
		},
		{
			name: "",
			fields: fields{
				logger: log.New(os.Stderr, "", log.LstdFlags),
			},
			args: args{
				query: "select * from ddl_user where (`school` = 'Shanxi Univ.' and `delete_at` is null) order by `create_at` desc limit ?",
				args:  []interface{}{20},
			},
		},
		{
			name: "",
			fields: fields{
				logger: log.New(os.Stderr, "", log.LstdFlags),
			},
			args: args{
				query: "select * from ddl_user where (`school` = ? and `delete_at` is null) order by `create_at` desc",
				args:  []interface{}{"Shanxi Univ."},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := &SqlLogger{
				logger: tt.fields.logger,
			}
			receiver.Log(tt.args.query, tt.args.args...)
		})
	}
}
