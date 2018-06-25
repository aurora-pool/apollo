package stats

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

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

func fetchUrl(url *url.URL) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url.String(), nil)
	resp, _ := client.Do(req)

	return resp
}
