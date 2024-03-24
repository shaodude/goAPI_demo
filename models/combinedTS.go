package models

import "gorm.io/gorm"

// student object
type Student struct {
	ID        uint      // Primary key
	Email     string    `gorm:"unique"`
	Teachers  []Teacher `gorm:"many2many:teacher_students"` // Many-to-many relationship with teachers
	Suspended bool
}

// teacher object
type Teacher struct {
	ID       uint      // Primary key
	Email    string    `gorm:"unique"`
	Students []Student `gorm:"many2many:teacher_students"` // Many-to-many relationship with students
}

// migrate tables
func MigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&Teacher{}, &Student{}); err != nil {
		return err
	}
	return nil
}
