package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

var RedisPool *redis.Pool

const (
	RedisHost = "localhost"
	RedisPort = "6379"
)

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
	wshandler(c.Writer, c.Request)
}

func (ctr ChannelCtrl) InitRedis() {
	InitRedis()
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header["Origin"]

		if len(origin) == 0 {
			return true
		}

		u, err := url.Parse(origin[0])

		if err != nil {
			return false
		}

		fmt.Println(u, r.URL)
		rUrl, _ := url.Parse(r.RequestURI)
		return isEqDomain(u, rUrl)
	},
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

func wshandler(w http.ResponseWriter, r *http.Request) {
	socketConnection, _ := wsupgrader.Upgrade(w, r, nil)
	redisConnection := RedisPool.Get()
	clientClosed := make(chan bool, 1)
	poolStats, _ := redis.Bytes(redisConnection.Do("get", "aurora-pool:stats"))

	go func(socketConnection *websocket.Conn, clientClosed chan bool) {
		for {
			_, _, err := socketConnection.ReadMessage()
			if err != nil {
				// We are done here
				clientClosed <- true
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v, user-agent: %v", err, r.Header.Get("User-Agent"))
				}

				socketConnection.Close()
			}
		}
	}(socketConnection, clientClosed)

	go func(socketConnection *websocket.Conn, clientClosed chan bool) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		socketConnection.WriteMessage(websocket.TextMessage, poolStats)
		redisConnection.Close()

		for {
			select {
			case <-ticker.C:
				socketConnection.WriteMessage(websocket.TextMessage, poolStats)
				redisConnection.Close()
			case <-clientClosed:
				return
			}
		}
	}(socketConnection, clientClosed)
}

func InitRedis() {
	RedisPool = createRedisPool()
}

type User struct {
	Address            string  `json:"address"`
	OutStandingBalance float64 `json:"balance"`
	PaidBalance        float64 `json:"paid"`
	Hashrate           float64 `json:"hashrate"`
}

func isEqDomain(lhs *url.URL, rhs *url.URL) bool {
	lhsHost := lhs.Host
	rhsHost := rhs.Host

	if partsLhs := strings.Split(lhsHost, "."); len(partsLhs) > 2 {
		lhsHost = strings.Join(partsLhs[1:], ".")
	}

	if partsLhs := strings.Split(rhsHost, "."); len(partsLhs) > 2 {
		rhsHost = strings.Join(partsLhs[1:], ".")
	}

	return rhsHost == lhsHost
}
