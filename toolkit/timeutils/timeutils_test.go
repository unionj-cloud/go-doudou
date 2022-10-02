package timeutils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type args struct {
		t          string
		defaultDur time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				t:          "9h1m30s",
				defaultDur: 15 * time.Second,
			},
			want:    32490000000000,
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				t:          "wrongdurationstr",
				defaultDur: 15 * time.Second,
			},
			want:    15000000000,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.t, tt.args.defaultDur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallWithCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := CallWithCtx(ctx, func() struct{} {
		time.Sleep(2 * time.Second)
		fmt.Println("Job Done")
		return struct{}{}
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
}
