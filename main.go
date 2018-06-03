package main

import (
	"github.com/aurora-pool/apollo/controllers"
	"github.com/gin-gonic/gin"
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]string{"message": "API: welcome to the jungle!"})
	})

	api := router.Group("/api/v1")
	channelCtrl := controllers.ChannelCtrl{}
	channelCtrl.SetDB(nil)

	api.GET("/channels", channelCtrl.ChannelIndex)

	router.Run(":8242")
}
