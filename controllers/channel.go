package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	connection, _ := wsupgrader.Upgrade(w, r, nil)
	clientClosed := make(chan bool, 1)

	go func(connection *websocket.Conn, clientClosed chan bool) {
		for {
			_, _, err := connection.ReadMessage()
			if err != nil {
				// We are done here
				clientClosed <- true
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v, user-agent: %v", err, r.Header.Get("User-Agent"))
				}

				connection.Close()
			}
		}
	}(connection, clientClosed)

	go func(connection *websocket.Conn, clientClosed chan bool) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		connection.WriteJSON(User{
			Address:            "QX87 QX87 QX87 QX87 QX87 QX87 QX87 QX87",
			OutStandingBalance: 220.23,
			PaidBalance:        1245.96,
			Hashrate:           100.12,
		})

		for {
			select {
			case <-ticker.C:
				// we need to get these values from db.
				connection.WriteJSON(User{
					Address:            "QX87 QX87 QX87 QX87 QX87 QX87 QX87 QX87",
					OutStandingBalance: 220.23,
					PaidBalance:        1245.96,
					Hashrate:           100.12,
				})
			case <-clientClosed:
				return
			}
		}
	}(connection, clientClosed)
}

type User struct {
	Address            string  `json:"address"`
	OutStandingBalance float64 `json:"balance"`
	PaidBalance        float64 `json:"paid"`
	Hashrate           float64 `json:"hashrate"`
}
