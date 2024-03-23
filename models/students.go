package models

import "gorm.io/gorm"

type Students struct {
	ID    uint    `gorm:"primary key;autoIncrement" json:"id"`
	Email *string `json:"email"`
}

func MigrateStudents(db *gorm.DB) error {
	err := db.AutoMigrate(&Students{})
	return err
}
