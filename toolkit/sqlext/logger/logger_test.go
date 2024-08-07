package logger_test

import (
	"context"
	"os"
	"testing"

	"github.com/ascarter/requestid"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/unionj-cloud/go-doudou/v2/framework/tracing"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/sqlext/logger"
)

func TestMain(m *testing.M) {
	os.Setenv("GDD_SERVICE_NAME", "TestSqlLogger")
	os.Setenv("GDD_SQL_LOG_ENABLE", "true")

	tracer, closer := tracing.Init()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	m.Run()
}

func TestSqlLogger_Log(t *testing.T) {
	type fields struct {
		ctx func() context.Context
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
				ctx: func() context.Context {
					return requestid.NewContext(context.Background(), uuid.NewString())
				},
			},
			args: args{
				query: "select * from ddl_user where (`school` = 'Shanxi Univ.' and `delete_at` is null) order by `create_at` desc limit ?,?",
				args:  []interface{}{0, 2},
			},
		},
		{
			name: "",
			fields: fields{
				ctx: func() context.Context {
					_, ctx := opentracing.StartSpanFromContext(requestid.NewContext(context.Background(), uuid.NewString()), "TestSqlLogger")
					return ctx
				},
			},
			args: args{
				query: "select * from ddl_user where (`school` = 'Shanxi Univ.' and `delete_at` is null) order by `create_at` desc limit ?",
				args:  []interface{}{20},
			},
		},
		{
			name: "",
			fields: fields{
				ctx: func() context.Context {
					return opentracing.ContextWithSpan(context.Background(), mocktracer.New().StartSpan("TestSqlLogger"))
				},
			},
			args: args{
				query: "select * from ddl_user where (`school` = ? and `delete_at` is null) order by `create_at` desc",
				args:  []interface{}{"Shanxi Univ."},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := logger.NewSqlLogger()
			receiver.Log(tt.fields.ctx(), tt.args.query, tt.args.args...)
		})
	}
}
