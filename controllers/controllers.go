package controllers

import (
	"github.com/aurora-pool/apollo/hub"
	"github.com/jinzhu/gorm"
)

type UserModel struct {
	Id      int `gorm:"primary_key"`
	Address string
}

type Controller struct {
	DB  *gorm.DB
	hub *hub.Hub
}

func (ctr *Controller) SetDB(db *gorm.DB) {
	ctr.DB = db
}

func (ctr *Controller) SetHub(hub *hub.Hub) {
	ctr.hub = hub
}
