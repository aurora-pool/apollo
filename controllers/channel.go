package controllers

import (
	"github.com/gin-gonic/gin"
)

type ChannelCtrl struct {
	Controller
}

func (ctr ChannelCtrl) ChannelIndex(c *gin.Context) {
	c.JSON(200, map[string]string{"message": "Coming soon"})
}
