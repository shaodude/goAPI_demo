package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/shaodude/goAPI_demo/models"
	"github.com/shaodude/goAPI_demo/storage"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Register(context *fiber.Ctx) error {
	// Parse request body
	var requestBody struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}

	if err := context.BodyParser(&requestBody); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Check if teacher exists
	var teacher models.Teacher
	if err := r.DB.Where("email = ?", requestBody.Teacher).First(&teacher).Error; err != nil {
		// If teacher does not exist, create a new teacher
		teacher = models.Teacher{Email: requestBody.Teacher}
		if err := r.DB.Create(&teacher).Error; err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create teacher"})
		}
	}

	// Fetch or create students and associate them with the teacher
	for _, studentEmail := range requestBody.Students {
		var student models.Student
		if err := r.DB.Where("email = ?", studentEmail).First(&student).Error; err != nil {
			// If student not found, create a new student
			student = models.Student{Email: studentEmail}
			if err := r.DB.Create(&student).Error; err != nil {
				return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create student"})
			}
		}
		if err := r.DB.Model(&teacher).Association("Students").Append(&student); err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to register student"})
		}
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "registered!"})
}

func (r *Repository) GetCommonStudents(context *fiber.Ctx) error {
	// Parse query parameters to extract the list of teachers
	teachers := context.Query("teacher")
	fmt.Println(teachers)

	// Split the list of teachers into individual email addresses
	teacherList := strings.Split(teachers, ",")
	fmt.Println(teacherList)
	// Fetch students common to all the given teachers
	var commonStudents []models.Student

	// Prepare the base query to join students and teachers
	db := r.DB.Table("students")

	// Join the teacher_students and teachers tables for each teacher in the list
	for idx, teacherEmail := range teacherList {
		db = db.Joins(fmt.Sprintf("JOIN teacher_students ts%d ON students.id = ts%d.student_id", idx, idx)).
			Joins(fmt.Sprintf("JOIN teachers t%d ON ts%d.teacher_id = t%d.id", idx, idx, idx)).
			Where(fmt.Sprintf("t%d.email = ?", idx), teacherEmail)
	}

	// Group the results by student ID and ensure each student is associated with all teachers
	groupConditions := make([]string, len(teacherList))
	for idx := range teacherList {
		groupConditions[idx] = fmt.Sprintf("COUNT(DISTINCT t%d.id) = 1", idx)
	}
	groupCondition := strings.Join(groupConditions, " AND ")
	db = db.Group("students.id").Having(groupCondition)

	// Execute the query
	if err := db.Find(&commonStudents).Error; err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve common students"})
	}

	// Extract student emails from the results
	var studentEmails []string
	for _, student := range commonStudents {
		studentEmails = append(studentEmails, student.Email)
	}

	// Return success response with list of common students
	return context.Status(http.StatusOK).JSON(fiber.Map{"students": studentEmails})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", r.Register)
	api.Get("/commonstudents", r.GetCommonStudents)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateTables(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
