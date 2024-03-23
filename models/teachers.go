package models

import "gorm.io/gorm"

type Teachers struct {
	ID    uint    `gorm:"primary key;autoIncrement" json:"id"`
	Email *string `json:"email"`
}

func MigrateTeachers(db *gorm.DB) error {
	err := db.AutoMigrate(&Teachers{})
	return err
}
