package models

import "gorm.io/gorm"

// Student struct representing a student
type Student struct {
	ID       uint      // Primary key (you may need to adjust this based on your database setup)
	Email    string    `gorm:"unique"`                     // Unique constraint for email field
	Teachers []Teacher `gorm:"many2many:teacher_students"` // Many-to-many relationship with teachers
}

// Teacher struct representing a teacher
type Teacher struct {
	ID       uint      // Primary key (you may need to adjust this based on your database setup)
	Email    string    `gorm:"unique"`                     // Unique constraint for email field
	Students []Student `gorm:"many2many:teacher_students"` // Many-to-many relationship with students
}

// MigrateTables function to migrate tables
func MigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&Teacher{}, &Student{}); err != nil {
		return err
	}
	return nil
}
