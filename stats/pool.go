package stats

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

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
