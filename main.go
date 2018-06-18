package main

import (
	api_controllers "github.com/aurora-pool/apollo/api_controllers_v1"
	"github.com/aurora-pool/apollo/helpers"
	"github.com/aurora-pool/apollo/hub"
	"github.com/aurora-pool/apollo/stats"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	stats.InitRedis()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	db, err := gorm.Open("mysql", helpers.GetDBurl())
	defer db.Close()

	hub := hub.NewHub()
	go hub.Run()
	statsUpdater := stats.NewStats()
	go statsUpdater.Run(hub)

	if err != nil {
		panic(err)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]string{"message": "API: welcome to the jungle!"})
	})

	api := router.Group("/api/v1")
	channelCtrl := api_controllers.ChannelCtrl{}
	usersCtrl := api_controllers.UsersCtrl{}

	channelCtrl.SetDB(db)
	channelCtrl.SetHub(hub)

	usersCtrl.SetDB(db)

	api.GET("/channels", channelCtrl.ChannelIndex)
	api.GET("/ws", channelCtrl.WebSocket)
	api.GET("/users/:id", usersCtrl.Show)

	router.Run(":8442")
}
