package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

type LimiterI interface {
	Key(c *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBuckets(rules ...LimiterBucketRule) LimiterI
}

type Limiter struct {
	limiterBuckes map[string]*ratelimit.Bucket
}

type LimiterBucketRule struct {
	Key string
	FillInterval time.Duration
	Capacity int64
	Quantum int64
}


