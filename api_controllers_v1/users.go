package api_controller_v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type UsersCtrl struct {
	Controller
}

func (ctrl *UsersCtrl) Show(ctx *gin.Context) {
	userModel := UserModel{}
	err := ctrl.DB.First(&userModel, ctx.Param("id")).Error

	if err != nil {
		ctx.JSON(404, map[string]string{})
		return
	}

	ctx.JSON(200, map[string]string{
		"message": "Coming soon",
		"Id":      fmt.Sprintf("%d", userModel.Id),
		"address": userModel.Address,
	})
}
