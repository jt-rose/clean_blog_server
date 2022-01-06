package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	rdb "github.com/jt-rose/clean_blog_server/database"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func InitRateLimiter() (gin.HandlerFunc, error) {

	// set a limit rate of 200 requests per minute
	rate, err := limiter.NewRateFromFormatted("200-M")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(rdb.RedisClient, limiter.StoreOptions{
		Prefix:   "rate-limiter",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// return the customized gin middleware
	return mgin.NewMiddleware(limiter.New(store, rate)), nil
}
