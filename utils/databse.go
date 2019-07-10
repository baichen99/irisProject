package utils

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"irisProject/models"
)

// ConnectDB connect to a sqlite database
func ConnectDB() (db *gorm.DB) {
	dataURL := "test.db"
	db, err := gorm.Open("sqlite3", dataURL)
	if err != nil{
		panic(err)
	}
	return
}

// InitDB initialize database
func InitDB(db *gorm.DB) *gorm.DB {
	db.DropTableIfExists(&models.User{})
	db.AutoMigrate(&models.User{})
	return db
}