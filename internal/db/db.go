package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"interview-backend/internal/conf"
	"interview-backend/internal/model"
)

var db *gorm.DB

func Init(d *gorm.DB) {
	db = d
	err := autoMigrate(new(model.User), new(model.Token))
	if err != nil {
		log.Fatalf("failed migrate database: %s", err.Error())
	}
}

func GetDB() *gorm.DB {
	return db
}

func autoMigrate(dst ...interface{}) error {
	var err error
	if conf.Conf.Database.Type == "mysql" {
		err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(dst...)
	} else {
		err = db.AutoMigrate(dst...)
	}
	return err
}

func Close() {
	log.Infoln("closing database")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed opening connection to database: %s", err.Error())
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("failed closing connection to database: %s", err.Error())
		return
	}
}
