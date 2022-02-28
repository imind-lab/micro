package sentinel

import (
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"sync"
)

var (
	once    sync.Once
	limiter *rate.Limiter
)

func GetLimiter() *rate.Limiter {
	once.Do(func() {
		limit := viper.GetFloat64("service.concurrence.limit")
		capacity := viper.GetInt("service.concurrence.capacity")
		limiter = rate.NewLimiter(rate.Limit(limit), capacity)
	})
	return limiter
}
