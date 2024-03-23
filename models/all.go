package models

import "gorm.io/gorm"

// Teacher struct representing a teacher
type Teacher struct {
	ID       uint      // Primary key (you may need to adjust this based on your database setup)
	Email    string    `gorm:"unique"`                     // Unique constraint for email field
	Students []Student `gorm:"many2many:teacher_students"` // Many-to-many relationship with students
}

// Student struct representing a student
type Student struct {
	ID    uint   // Primary key (you may need to adjust this based on your database setup)
	Email string `gorm:"unique"` // Unique constraint for email field
}

// TeacherStudent struct representing the join table for teacher-student relationship
type TeacherStudent struct {
	TeacherID uint // Foreign key for teacher
	StudentID uint // Foreign key for student
}

// MigrateTables function to migrate tables
func MigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&Teacher{}, &Student{}, &TeacherStudent{}); err != nil {
		return err
	}
	return nil
}
