package chap05

import (
	"context"
	"golang.org/x/time/rate"
	"sort"
	"time"
)

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration/time.Duration(eventCount))
}

// Dummy API with two endpoints

type APIConnection struct {
	networkLimit,
	diskLimit,
	apiLimit RateLimiter
}

func Open() *APIConnection {
	return &APIConnection{
		apiLimit: MultiLimiter(
			rate.NewLimiter(Per(2, time.Second), 1),
			rate.NewLimiter(Per(10, time.Minute), 10),
		),
		diskLimit: MultiLimiter(
			rate.NewLimiter(rate.Limit(1), 1),
		),
		networkLimit: MultiLimiter(
			rate.NewLimiter(Per(3, time.Second), 3),
		),
	}
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.diskLimit).Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.networkLimit).Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}


// 流量制限のまとめやく
type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}