package controllers

import (
	"github.com/jinzhu/gorm"
)

type UserModel struct {
	Id      int `gorm:"primary_key"`
	Address string
}

type Controller struct {
	DB *gorm.DB
}

func (ctr *Controller) SetDB(db *gorm.DB) {
	ctr.DB = db
}
