package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Getdb() *gorm.DB {
	dsn := "nadeer:nadeer@123@tcp(localhost:3306)/bookstore?charset=utf8&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return d
}
