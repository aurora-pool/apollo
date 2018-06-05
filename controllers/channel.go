package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChannelCtrl struct {
	Controller
}

func (ctr ChannelCtrl) ChannelIndex(c *gin.Context) {
	c.JSON(200, map[string]string{"message": "Coming soon"})
}

func (ctr ChannelCtrl) WebSocket(c *gin.Context) {
	wshandler(c.Writer, c.Request)
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	connection, _ := wsupgrader.Upgrade(w, r, nil)

	go func(connection *websocket.Conn) {
		for {
			_, _, err := connection.ReadMessage()
			if err != nil {
				fmt.Println("Failed to set websocket upgrade: %+v", err)
				connection.Close()
			}
		}
	}(connection)

	go func(connection *websocket.Conn) {
		ch := time.Tick(5 * time.Second)

		for range ch {
			// we need to get these values from db.
			connection.WriteJSON(User{
				Address:            "QX87 QX87 QX87 QX87 QX87 QX87 QX87 QX87",
				OutStandingBalance: 220.23,
				PaidBalance:        1245.96,
				Hashrate:           100.12,
			})
		}
	}(connection)
}

type User struct {
	Address            string  `json:"address"`
	OutStandingBalance float64 `json:"balance"`
	PaidBalance        float64 `json:"paid"`
	Hashrate           float64 `json:"hashrate"`
}
