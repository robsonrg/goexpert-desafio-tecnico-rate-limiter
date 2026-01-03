package api_test

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/configs"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/ratelimit"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/ratelimit/limiter"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/web"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/web/webserver"
	"github.com/robsonrg/goexpert-desafio-tecnico-rate-limiter/internal/infra/web/webserver/middleware"
)

type APITestSuite struct {
	suite.Suite
	rdb                             *redis.Client
	httpClient                      *http.Client
	URL                             string
	RateLimitRequestPerSecondsIP    int64
	RateLimitRequestPerSecondsToken int64
	ConcurrentRequests              int
}

func (suite *APITestSuite) SetupTest() {
	configs, err := configs.LoadConfig(".")
	assert.NoError(suite.T(), err)
	suite.RateLimitRequestPerSecondsIP = configs.RateLimitRequestPerSecondsIP
	suite.RateLimitRequestPerSecondsToken = configs.RateLimitRequestPerSecondsToken

	suite.rdb = redis.NewClient(&redis.Options{
		Addr:     configs.RedisHost,
		Password: configs.RedisPassword,
		DB:       0,
	})

	webHandler := web.NewWebHandler()
	ratelimit := ratelimit.NewRateLimit(limiter.NewRedisLimiter(suite.rdb), suite.RateLimitRequestPerSecondsIP, suite.RateLimitRequestPerSecondsToken)
	ratelimitMiddleware := middleware.NewWebServerRateLimiter(ratelimit)

	webserver := webserver.NewWebServerTest()
	webserver.AddMiddleware("rate_limiter", ratelimitMiddleware.Handle)
	webserver.AddHandler("/sample", webHandler.Handle)

	ts := webserver.Start()
	suite.URL = ts.URL
	suite.httpClient = ts.Client()

	suite.ConcurrentRequests = 10
}

func (suite *APITestSuite) TearDownTest() {
	suite.rdb.Close()
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) TestRateLimitRequestPerIP() {
	numberOfRequests := 5000
	expectedAllowRequests := suite.RateLimitRequestPerSecondsIP
	expectedNotAllowRequests := int64(numberOfRequests) - suite.RateLimitRequestPerSecondsIP

	var allowRequests int64 = 0
	var notAllowRequests int64 = 0
	concurrentControl := make(chan struct{}, suite.ConcurrentRequests)
	var wg sync.WaitGroup
	wg.Add(numberOfRequests)

	for i := 0; i < numberOfRequests; i++ {
		concurrentControl <- struct{}{}
		go func() {
			defer func() {
				<-concurrentControl
				wg.Done()
			}()
			req, err := http.NewRequest("GET", suite.URL+"/sample", nil)
			req.Header.Add("X-Forwarded-For", "127.0.0.1")
			assert.NoError(suite.T(), err)

			resp, err := suite.httpClient.Do(req)
			assert.NoError(suite.T(), err)

			if resp.StatusCode == http.StatusOK {
				atomic.AddInt64(&allowRequests, 1)
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				atomic.AddInt64(&notAllowRequests, 1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(suite.T(), expectedAllowRequests, allowRequests)
	assert.Equal(suite.T(), expectedNotAllowRequests, notAllowRequests)
}

func (suite *APITestSuite) TestRateLimitRequestPerToken() {
	numberOfRequests := 5000
	expectedAllowRequests := suite.RateLimitRequestPerSecondsToken
	expectedNotAllowRequests := int64(numberOfRequests) - suite.RateLimitRequestPerSecondsToken

	concurrentControl := make(chan struct{}, suite.ConcurrentRequests)
	var allowRequests int64 = 0
	var notAllowRequests int64 = 0
	var wg sync.WaitGroup
	wg.Add(numberOfRequests)

	for i := 0; i < numberOfRequests; i++ {
		concurrentControl <- struct{}{}
		go func() {
			defer func() {
				<-concurrentControl
				wg.Done()
			}()
			req, err := http.NewRequest("GET", suite.URL+"/sample", nil)
			req.Header.Add("API_KEY", "some-token")
			req.Header.Add("X-Forwarded-For", "127.0.0.2")
			assert.NoError(suite.T(), err)

			resp, err := suite.httpClient.Do(req)
			assert.NoError(suite.T(), err)

			if resp.StatusCode == http.StatusOK {
				atomic.AddInt64(&allowRequests, 1)
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				atomic.AddInt64(&notAllowRequests, 1)
			}
		}()
	}
	wg.Wait()

	assert.Equal(suite.T(), expectedAllowRequests, allowRequests)
	assert.Equal(suite.T(), expectedNotAllowRequests, notAllowRequests)
}
