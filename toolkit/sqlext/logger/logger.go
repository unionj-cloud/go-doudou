package logger

import (
	"context"
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/uber/jaeger-client-go"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/reflectutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"os"
	"regexp"
	"strings"
)

type SqlLogger struct {
	Logger zerolog.Logger
}

func (receiver SqlLogger) Enable() bool {
	return cast.ToBoolOrDefault(os.Getenv("GDD_SQL_LOG_ENABLE"), false)
}

func (receiver SqlLogger) Log(ctx context.Context, query string, args ...interface{}) {
	if !receiver.Enable() {
		return
	}
	receiver.LogWithErr(ctx, nil, nil, query, args...)
}

type SqlLoggerOption func(logger *SqlLogger)

func WithLogger(logger zerolog.Logger) SqlLoggerOption {
	return func(sqlLogger *SqlLogger) {
		sqlLogger.Logger = logger
	}
}

func NewSqlLogger(opts ...SqlLoggerOption) SqlLogger {
	sqlLogger := SqlLogger{
		Logger: zlogger.Logger,
	}
	for _, item := range opts {
		item(&sqlLogger)
	}
	return sqlLogger
}

var limitre *regexp.Regexp

func init() {
	limitre = regexp.MustCompile(`limit '\d+'(,'\d+')?`)
}

func PopulatedSql(query string, args ...interface{}) string {
	query = strings.Join(strings.Fields(query), " ")
	copiedArgs := make([]interface{}, len(args))
	copy(copiedArgs, args)
	for i, arg := range copiedArgs {
		if arg == nil {
			continue
		}
		value := reflectutils.ValueOf(arg)
		if value.IsValid() {
			copiedArgs[i] = value.Interface()
		}
	}
	str := strings.ReplaceAll(fmt.Sprintf(strings.ReplaceAll(query, "?", "'%v'"), copiedArgs...), "'<nil>'", "null")
	if limitre.MatchString(str) {
		str = limitre.ReplaceAllStringFunc(str, func(s string) string {
			return strings.ReplaceAll(s, "'", "")
		})
	}
	return str
}

func (receiver SqlLogger) LogWithErr(ctx context.Context, err error, hit *bool, query string, args ...interface{}) {
	if !receiver.Enable() {
		return
	}
	var sb strings.Builder
	if reqId, ok := requestid.FromContext(ctx); ok && stringutils.IsNotEmpty(reqId) {
		sb.WriteString(fmt.Sprintf("RequestID: %s\t", reqId))
	}
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if jspan, ok := span.(*jaeger.Span); ok {
			sb.WriteString(fmt.Sprintf("TraceID: %s\t", jspan.SpanContext().TraceID().String()))
		} else {
			sb.WriteString(fmt.Sprintf("TraceID: %s\t", span))
		}
	}
	sb.WriteString(fmt.Sprintf("SQL: %s", PopulatedSql(query, args...)))
	if hit != nil {
		sb.WriteString(fmt.Sprintf("\tHIT: %t", *hit))
	}
	if err != nil {
		sb.WriteString(fmt.Sprintf("\tERR: %s", errors.Wrap(err, caller.NewCaller().String())))
	}
	receiver.Logger.Info().Msgf(sb.String())
}
