package main

import (
	"github.com/aurora-pool/apollo/controllers"
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

	db, err := gorm.Open("mysql", "AuroraRoot:Enterprise32@tcp(us-east.cvvg11be4uiw.us-east-2.rds.amazonaws.com:3306)/pool?charset=utf8&parseTime=True&loc=Local")
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
	channelCtrl := controllers.ChannelCtrl{}
	usersCtrl := controllers.UsersCtrl{}

	channelCtrl.SetDB(db)
	channelCtrl.SetHub(hub)

	usersCtrl.SetDB(db)

	api.GET("/channels", channelCtrl.ChannelIndex)
	api.GET("/ws", channelCtrl.WebSocket)
	api.GET("/users/:id", usersCtrl.Show)

	router.Run(":8442")
}
