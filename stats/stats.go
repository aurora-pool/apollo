package stats

import (
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	RedisHost      = "localhost"
	RedisPort      = "6379"
	globalStatsUrl = "https://nimiq.mopsus.com/api/quick-stats"
)

var RedisPool *redis.Pool

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
	go minersStats(hub, st)
}

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
