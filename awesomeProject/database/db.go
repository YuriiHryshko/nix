package database

import (
	"awesomeProject/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Perform database migrations
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
