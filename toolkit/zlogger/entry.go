package zlogger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/constants"
	"io"
	"os"
	"strconv"
)

var Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

type LoggerConfig struct {
	Dev    bool
	Caller bool
	Pid    bool
}

type LoggerConfigOption func(*LoggerConfig)

func WithDev(dev bool) LoggerConfigOption {
	return func(lc *LoggerConfig) {
		lc.Dev = dev
	}
}

func WithCaller(caller bool) LoggerConfigOption {
	return func(lc *LoggerConfig) {
		lc.Caller = caller
	}
}

func NewLoggerConfig(opts ...LoggerConfigOption) LoggerConfig {
	lc := LoggerConfig{}
	for _, item := range opts {
		item(&lc)
	}
	return lc
}

func InitEntry(levelStr string, lc LoggerConfig) {
	var output io.Writer
	if lc.Dev {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: constants.FORMAT}
	} else {
		output = os.Stdout
	}
	level, _ := zerolog.ParseLevel(levelStr)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zeroCtx := zerolog.New(output).Level(level).With().Timestamp().Stack()
	if lc.Caller {
		zeroCtx = zeroCtx.Caller()
	}
	if lc.Pid {
		zeroCtx = zeroCtx.Str("__pid", strconv.Itoa(os.Getpid()))
	}
	Logger = zeroCtx.Logger()
}

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) zerolog.Logger {
	return Logger.Output(w)
}

// With creates a child logger with the field added to its context.
func With() zerolog.Context {
	return Logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level zerolog.Level) zerolog.Logger {
	return Logger.Level(level)
}

// Sample returns a logger with the s sampler.
func Sample(s zerolog.Sampler) zerolog.Logger {
	return Logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func Hook(h zerolog.Hook) zerolog.Logger {
	return Logger.Hook(h)
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *zerolog.Event {
	return Logger.Err(err)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return Logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return Logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *zerolog.Event {
	return Logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return Logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
