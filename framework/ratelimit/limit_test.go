package ratelimit

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Limit
		wantErr bool
	}{
		{
			name: "",
			args: args{
				value: "0.0055-S-20",
			},
			want: Limit{
				Rate:   0.0055,
				Burst:  20,
				Period: time.Second,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				value: "1000-H",
			},
			want: Limit{
				Rate:   1000,
				Burst:  1,
				Period: time.Hour,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				value: "1000-H-1-1-1",
			},
			want:    Limit{},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				value: "1000-Y",
			},
			want:    Limit{},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				value: "a-H",
			},
			want:    Limit{},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				value: "1000-Y-100",
			},
			want:    Limit{},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				value: "a-H-100",
			},
			want:    Limit{},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				value: "10-H-abc",
			},
			want:    Limit{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerSecond(t *testing.T) {
	Convey("Should equal to 10 in 1s", t, func() {
		So(PerSecond(10), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Second,
			Burst:  1,
		})
	})

	Convey("Should equal to 10 in 1s with 100 burst", t, func() {
		So(PerSecondBurst(10, 100), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Second,
			Burst:  100,
		})
	})

	Convey("Should equal to 10 in 1m", t, func() {
		So(PerMinute(10), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Minute,
			Burst:  1,
		})
	})

	Convey("Should equal to 10 in 1m with 100 burst", t, func() {
		So(PerMinuteBurst(10, 100), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Minute,
			Burst:  100,
		})
	})

	Convey("Should equal to 10 in 1h", t, func() {
		So(PerHour(10), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Hour,
			Burst:  1,
		})
	})

	Convey("Should equal to 10 in 1h with 100 burst", t, func() {
		So(PerHourBurst(10, 100), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Hour,
			Burst:  100,
		})
	})

	Convey("Should equal to 10 in 1d", t, func() {
		So(PerDay(10), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Hour * 24,
			Burst:  1,
		})
	})

	Convey("Should equal to 10 in 1d with 100 burst", t, func() {
		So(PerDayBurst(10, 100), ShouldResemble, Limit{
			Rate:   10,
			Period: time.Hour * 24,
			Burst:  100,
		})
	})
}
