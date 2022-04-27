package logger

import (
	"context"
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/reflectutils"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"os"
	"regexp"
	"strings"
)

type ISqlLogger interface {
	Log(ctx context.Context, query string, args ...interface{})
	Enable() bool
}

type SqlLogger struct {
	logger logrus.StdLogger
}

func (receiver *SqlLogger) Enable() bool {
	return cast.ToBoolOrDefault(os.Getenv("GDD_SQL_LOG_ENABLE"), false)
}

func (receiver *SqlLogger) SetLogger(logger logrus.FieldLogger) {
	receiver.logger = logger
}

func NewSqlLogger(logger logrus.StdLogger) ISqlLogger {
	return &SqlLogger{logger: logger}
}

func (receiver *SqlLogger) Log(ctx context.Context, query string, args ...interface{}) {
	if receiver.Enable() {
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
		re := regexp.MustCompile(`limit '\d+'(,'\d+')?`)
		if re.MatchString(str) {
			str = re.ReplaceAllStringFunc(str, func(s string) string {
				return strings.ReplaceAll(s, "'", "")
			})
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
		sb.WriteString(fmt.Sprintf("SQL: %s", str))
		receiver.logger.Println(sb.String())
	}
}
