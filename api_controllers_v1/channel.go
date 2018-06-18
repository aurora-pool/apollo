package api_controller_v1

import (
	"github.com/aurora-pool/apollo/hub"
	"github.com/gin-gonic/gin"
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
	hub.ServeWs(ctr.hub, c.Writer, c.Request)
}
