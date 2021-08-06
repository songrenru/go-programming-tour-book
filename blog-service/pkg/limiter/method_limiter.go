package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"strings"
)

type MethodLimiter struct {
	*Limiter
}

func NewMethodLimiter() MethodLimiter {
	return MethodLimiter{
		&Limiter{limiterBuckes: map[string]*ratelimit.Bucket{}},
	}
}

func (l MethodLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	idx := strings.Index(uri, "?")
	if idx == -1 {
		return uri
	}

	return uri[:idx]
}

func (l MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := l.limiterBuckes[key]
	return bucket, ok
}

func (l MethodLimiter) AddBuckets(rules ...LimiterBucketRule) LimiterI {
	for _, rule := range rules {
		if _, ok := l.limiterBuckes[rule.Key]; !ok {
			l.limiterBuckes[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}

	return l
}
