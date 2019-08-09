package utils

import (
	"irisProject/models"

	"irisProject/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ConnectDB connect to a psql database
func ConnectDB() (db *gorm.DB) {
	DatabaseURI := "host=" + config.Conf.Postgres.Host + " port=" +
		config.Conf.Postgres.Port + " user=" + config.Conf.Postgres.User + " dbname=" +
		config.Conf.Postgres.Database + " password=" + config.Conf.Postgres.Password + " sslmode=disable"
	db, err := gorm.Open("postgres", DatabaseURI)
	if err != nil {
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

// InitAdmin add admin
func InitAdmin(db *gorm.DB) {
	password, _ := HashPassword("password")
	user1 := models.User{
		Username: "admin",
		Password: password,
		Role:     "admin",
	}
	db.Create(&user1)
}
