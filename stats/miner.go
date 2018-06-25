package stats

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func minersStats(hub Broadcastable, stats *Stats) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			redisConn := RedisPool.Get()
			stats, _ := redis.Strings(redisConn.Do("zrange", "miners:active", 1, -1))
			redisConn.Close()

			urlsJson, err := json.Marshal(stats)
			if err != nil {
				log.Println("ERROR:minerStats: Cannot parse json!")
				continue
			}

			formattedStats := fmt.Sprintf(`{"type":"pool:stats:miners","payload":%s}`, string(urlsJson))

			hub.Send([]byte(formattedStats))
		case <-stats.Closed:
			return
		}
	}
}
