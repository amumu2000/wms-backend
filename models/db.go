package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
)

var (
	db *gorm.DB
)

func Init(username, password, host string, port int, dbname string) {
	dsn := username + ":" + password + "@" + "tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbname

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect database error:%s\n", err.Error())
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("get db error:%s\n", err.Error())
	}

	sqlDb.SetMaxOpenConns(100)
}
