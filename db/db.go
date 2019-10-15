package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func InitDB(dbType string,url string) {
	db, err := gorm.Open(dbType,url)
	if err != nil {
		panic(err)
	}
	DB = db
}
