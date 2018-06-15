package controllers

import (
	"log"
	"os"
	"time"

	"github.com/aurora-pool/apollo/hub"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

const (
	RedisHost = "localhost"
	RedisPort = "6379"
)

var RedisPool *redis.Pool

type ChannelCtrl struct {
	Controller
}

func (UserModel) TableName() string {
	return "user"
}

func (ctr ChannelCtrl) ChannelIndex(c *gin.Context) {
	c.JSON(200, map[string]string{"message": "Coming soon"})
}

func (ctr ChannelCtrl) WebSocket(c *gin.Context) {
	client := hub.ServeWs(ctr.hub, c.Writer, c.Request)

	go func(hub *hub.Hub, c *hub.Client) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				redisConn := RedisPool.Get()
				stats, _ := redis.Bytes(redisConn.Do("get", "aurora-pool:stats"))
				redisConn.Close()
				hub.Broadcast <- stats
			case <-c.Closed:
				return
			}
		}
	}(ctr.hub, client)
}

func InitRedis() {
	RedisPool = createRedisPool()
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

type User struct {
	Address            string  `json:"address"`
	OutStandingBalance float64 `json:"balance"`
	PaidBalance        float64 `json:"paid"`
	Hashrate           float64 `json:"hashrate"`
}
