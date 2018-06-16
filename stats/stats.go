package stats

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	RedisHost      = "localhost"
	RedisPort      = "6379"
	globalStatsUrl = "https://nimiq.mopsus.com/api/quick-stats"
)

type Stats struct {
	Closed chan bool
}

func NewStats() *Stats {
	return &Stats{Closed: make(chan bool)}
}

type Broadcastable interface {
	Send(data []byte)
}

func (st *Stats) Run(hub Broadcastable) {
	go poolStats(hub, st)
	go globalStats(hub, st)
}

func globalStats(hub Broadcastable, stats *Stats) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hub.Send(getGlobalStats())
		case <-stats.Closed:
			return
		}
	}
}

func poolStats(hub Broadcastable, stats *Stats) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			redisConn := RedisPool.Get()
			stats, _ := redis.Bytes(redisConn.Do("get", "aurora-pool:stats"))
			redisConn.Close()
			hub.Send(stats)
		case <-stats.Closed:
			return
		}
	}
}

func getGlobalStats() []byte {
	parsedURL, _ := url.Parse(globalStatsUrl)
	resp := fetchUrl(parsedURL)
	log.Println(resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	formattedGlobalStats := fmt.Sprintf(`{"type":"global:stats","payload":%s}`, body)

	if err != nil {
		log.Fatal(err)
	}

	return []byte(formattedGlobalStats)
}

func fetchUrl(url *url.URL) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url.String(), nil)
	resp, _ := client.Do(req)

	return resp
}

var RedisPool *redis.Pool

func createRedisPool() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 50 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				getRedisUrl(),
				redis.DialConnectTimeout(
					10*time.Second,
				),
				redis.DialDatabase(0),
			)
		},
	}

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")

	if err != nil {
		log.Printf("Could not connect to redis on %s", getRedisUrl())
		panic(err)
	}

	return pool
}

func getRedisUrl() string {
	if redisEnv := os.Getenv("REDIS_URL"); len(redisEnv) > 1 {
		return redisEnv
	}
	return RedisHost + ":" + RedisPort
}

func InitRedis() {
	RedisPool = createRedisPool()
}
