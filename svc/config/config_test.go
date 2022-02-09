package config

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLogLevel_Decode(t *testing.T) {
	var ll1 LogLevel
	var ll2 LogLevel
	var ll3 LogLevel
	var ll4 LogLevel
	var ll5 LogLevel
	var ll6 LogLevel
	var ll7 LogLevel
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		ll      *LogLevel
		args    args
		wantErr bool
		want    logrus.Level
	}{
		{
			name: "1",
			ll:   &ll1,
			args: args{
				value: "panic",
			},
			wantErr: false,
			want:    logrus.PanicLevel,
		},
		{
			name: "2",
			ll:   &ll2,
			args: args{
				value: "debug",
			},
			wantErr: false,
			want:    logrus.DebugLevel,
		},
		{
			name: "3",
			ll:   &ll3,
			args: args{
				value: "error",
			},
			wantErr: false,
			want:    logrus.ErrorLevel,
		},
		{
			name: "4",
			ll:   &ll4,
			args: args{
				value: "fatal",
			},
			wantErr: false,
			want:    logrus.FatalLevel,
		},
		{
			name: "5",
			ll:   &ll5,
			args: args{
				value: "info",
			},
			wantErr: false,
			want:    logrus.InfoLevel,
		},
		{
			name: "6",
			ll:   &ll6,
			args: args{
				value: "warn",
			},
			wantErr: false,
			want:    logrus.WarnLevel,
		},
		{
			name: "7",
			ll:   &ll7,
			args: args{
				value: "trace",
			},
			wantErr: false,
			want:    logrus.TraceLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ll.Decode(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if logrus.Level(*tt.ll) != tt.want {
				t.Errorf("Decode() want %v, got %v", tt.want, logrus.Level(*tt.ll))
			}
		})
	}
}

func TestInitEnv(t *testing.T) {
	assert.NotPanics(t, func() {
		InitEnv()
	})
}

func TestInitEnv1(t *testing.T) {
	os.Setenv("GDD_ENV", "test")
	assert.NotPanics(t, func() {
		InitEnv()
	})
}

func Test_envVariable_Load(t *testing.T) {
	os.Setenv("GDD_BANNER", "on")
	tests := []struct {
		name     string
		receiver envVariable
		want     string
	}{
		{
			name:     "",
			receiver: GddBanner,
			want:     "on",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.receiver.Load(); got != tt.want {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envVariable_Write(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name     string
		receiver envVariable
		args     args
		wantErr  bool
	}{
		{
			name:     "",
			receiver: GddBanner,
			args: args{
				value: "on",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.receiver.Write(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.receiver.Load() != "on" {
				t.Errorf("got = %v, want = %v", tt.receiver.Load(), "on")
			}
		})
	}
}

func Test_envVariable_String(t *testing.T) {
	GddBanner.Write("on")
	tests := []struct {
		name     string
		receiver envVariable
		want     string
	}{
		{
			name:     "",
			receiver: GddBanner,
			want:     "on",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.receiver.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
