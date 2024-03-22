package configs

import (
	"fmt"
	"log"

	"mygram/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartMySQL(configDb models.ConfigMySQL) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		configDb.User, configDb.Password, configDb.Host, configDb.Port, configDb.DBName, configDb.DBCharset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database :", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})

	if err != nil {
		log.Panic("failed to migrate :", err)
	}

	DB = db
}

func StartPostgres(configDb models.ConfigPostgres) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		configDb.Host, configDb.User, configDb.Password, configDb.DBName, configDb.Port, configDb.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database :", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})

	if err != nil {
		log.Panic("failed to migrate :", err)
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
