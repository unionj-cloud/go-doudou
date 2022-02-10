package ratelimit

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type Limit struct {
	Rate   float64
	Burst  int
	Period time.Duration
}

type LimiterOption func(Limiter)

func PerSecond(rate float64) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Second,
		Burst:  1,
	}
}

func PerMinute(rate float64) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Minute,
		Burst:  1,
	}
}

func PerHour(rate float64) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Hour,
		Burst:  1,
	}
}

func PerDay(rate float64) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Hour * 24,
		Burst:  1,
	}
}

func PerSecondBurst(rate float64, burst int) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Second,
		Burst:  burst,
	}
}

func PerMinuteBurst(rate float64, burst int) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Minute,
		Burst:  burst,
	}
}

func PerHourBurst(rate float64, burst int) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Hour,
		Burst:  burst,
	}
}

func PerDayBurst(rate float64, burst int) Limit {
	return Limit{
		Rate:   rate,
		Period: time.Hour * 24,
		Burst:  burst,
	}
}

// Parse is borrowed from https://github.com/ulule/limiter Copyright (c) 2015-2018 Ulule
// You can use the simplified format "<limit>-<period>"(burst is 1 by default) or "<limit>-<period>-<burst>", with the given
// periods:
//
// * "S": second
// * "M": minute
// * "H": hour
// * "D": day
//
// Examples for "<limit>-<period>" format:
//
// * 5.5 reqs/second: "5.5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
//
// Examples:
//
// * 0.0055 reqs/second with burst 20: "0.0055-S-20" https://github.com/go-redis/redis_rate/issues/63
// * 10 reqs/minute with burst 200: "10-M-200"
// * 1000 reqs/hour with burst 500: "1000-H-500"
// * 2000 reqs/day with burst 1000: "2000-D-1000"
//
func Parse(value string) (Limit, error) {
	var l Limit

	splits := strings.Split(value, "-")
	if len(splits) != 2 && len(splits) != 3 {
		return l, errors.Errorf("incorrect format '%s'", value)
	}

	periods := map[string]time.Duration{
		"S": time.Second,    // Second
		"M": time.Minute,    // Minute
		"H": time.Hour,      // Hour
		"D": time.Hour * 24, // Day
	}

	if len(splits) == 2 {
		r, period := splits[0], strings.ToUpper(splits[1])

		p, ok := periods[period]
		if !ok {
			return l, errors.Errorf("incorrect period '%s'", period)
		}

		rate, err := strconv.ParseFloat(r, 64)
		if err != nil {
			return l, errors.Errorf("incorrect rate '%s'", r)
		}

		l = Limit{
			Rate:   rate,
			Burst:  1,
			Period: p,
		}

		return l, nil
	}

	r, period, b := splits[0], strings.ToUpper(splits[1]), splits[2]

	p, ok := periods[period]
	if !ok {
		return l, errors.Errorf("incorrect period '%s'", period)
	}

	rate, err := strconv.ParseFloat(r, 64)
	if err != nil {
		return l, errors.Errorf("incorrect rate '%s'", r)
	}

	burst, err := strconv.Atoi(b)
	if err != nil {
		return l, errors.Errorf("incorrect burst '%s'", b)
	}

	l = Limit{
		Rate:   rate,
		Burst:  burst,
		Period: p,
	}

	return l, nil
}
