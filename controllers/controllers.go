package controllers

import (
	"github.com/jinzhu/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func (ctr Controller) SetDB(db *gorm.DB) {
	ctr.DB = db
}
