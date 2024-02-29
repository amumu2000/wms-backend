package models

import "time"

type Log struct {
	ID      int64  `gorm:"column:id"`
	Title   string `gorm:"column:title"`
	Message string `gorm:"column:message"`
	Time    int64  `gorm:"column:time"`
	UserID  int64  `gorm:"column:user_id"`
}

func GetLogByID(id int64) *Log {
	var count int64
	log := Log{}
	db.Table("logs").Where("id = ?", id).Count(&count).First(&log)

	if count == 0 {
		return nil
	} else {
		return &log
	}
}

func GetLogs() []Log {
	var logs []Log
	db.Table("logs").Find(&logs)

	return logs
}

func AddLog(title, message string, userID int64) int64 {
	log := Log{
		Title:   title,
		Message: message,
		UserID:  userID,
		Time:    time.Now().Unix(),
	}

	db.Table("logs").Create(&log)

	return log.ID
}
