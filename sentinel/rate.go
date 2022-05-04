package sentinel

import (
	"sync"

	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

var (
	highOnce    sync.Once
	highLimiter *rate.Limiter

	lowOnce    sync.Once
	lowLimiter *rate.Limiter
)

func GetHighLimiter() *rate.Limiter {
	highOnce.Do(func() {
		limit := viper.GetFloat64("service.rate.high.limit")
		capacity := viper.GetInt("service.rate.high.capacity")
		highLimiter = rate.NewLimiter(rate.Limit(limit), capacity)
	})
	return highLimiter
}

func GetLowLimiter() *rate.Limiter {
	highOnce.Do(func() {
		limit := viper.GetFloat64("service.rate.low.limit")
		capacity := viper.GetInt("service.rate.low.capacity")
		lowLimiter = rate.NewLimiter(rate.Limit(limit), capacity)
	})
	return lowLimiter
}
