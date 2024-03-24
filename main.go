package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/shaodude/goAPI_demo/models"
	"github.com/shaodude/goAPI_demo/storage"
	"gorm.io/gorm"
)

// declare custom data types
type Repository struct {
	DB *gorm.DB
}

type registerRequest struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type SuspensionRequest struct {
	Student string `json:"student"`
}

type NotificationRequest struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}

type RetrieveResponse struct {
	Recipients []string `json:"recipients"`
}

// User story 1 POST API
func (r *Repository) Register(context *fiber.Ctx) error {

	// validate req
	var requestBody registerRequest
	if err := context.BodyParser(&requestBody); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// validate teacher
	var teacher models.Teacher
	if err := r.DB.Where("email = ?", requestBody.Teacher).First(&teacher).Error; err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": requestBody.Teacher + " not found"})
	}

	// validate student and create association
	for _, studentEmail := range requestBody.Students {
		var student models.Student
		if err := r.DB.Where("email = ?", studentEmail).First(&student).Error; err != nil {
			return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": studentEmail + " not found"})
		}
		if err := r.DB.Model(&teacher).Association("Students").Append(&student); err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to register student"})
		}
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{"message": "successfully registered student(s)"})
}

// User story 2 GET API
func (r *Repository) GetCommonStudents(context *fiber.Ctx) error {

	// get all occurrences of the "teacher" parameter
	teacherEmails := context.Context().QueryArgs().PeekMulti("teacher")
	var teacherList []string
	for _, emailBytes := range teacherEmails {
		teacherList = append(teacherList, string(emailBytes))
	}

	// validate teachers
	for _, teacherEmail := range teacherList {
		var teacher models.Teacher
		if err := r.DB.Where("email = ?", teacherEmail).First(&teacher).Error; err != nil {
			return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": teacherEmail + " not found"})
		}
	}

	// add students associated with teachers
	var teacherIDs []uint
	var studentIDs []uint
	// get teachers id based on email
	for _, teacherEmail := range teacherList {
		var teacher models.Teacher
		if err := r.DB.Where("email = ?", teacherEmail).First(&teacher).Error; err != nil {
			return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": teacherEmail + " not found"})
		}

		// handle duplicate teacher emails
		found := false
		for _, id := range teacherIDs {
			if id == teacher.ID {
				found = true
				break
			}
		}

		if !found {
			teacherIDs = append(teacherIDs, teacher.ID)
		}
	}

	// get student id associated with teacher id where occurrences = number of teachers (common among all)
	if err := r.DB.
		Table("teacher_students").
		Select("student_id").
		Where("teacher_id IN (?)", teacherIDs).
		Group("student_id").
		Having("COUNT(DISTINCT teacher_id) = ?", len(teacherIDs)).
		Pluck("student_id", &studentIDs).
		Error; err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to fetch student IDs"})
	}

	// get student emails from common student IDs
	var studentEmails []string
	if err := r.DB.
		Table("students").
		Select("email").
		Where("id IN (?)", studentIDs).
		Pluck("email", &studentEmails).
		Error; err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to fetch student emails"})
	}

	// return list of common students
	return context.Status(http.StatusOK).JSON(fiber.Map{"students": studentEmails})
}

// User story 3 POST API
func (r *Repository) SuspendStudent(context *fiber.Ctx) error {

	// validate req
	var requestBody SuspensionRequest
	if err := context.BodyParser(&requestBody); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// validate student
	var student models.Student
	if err := r.DB.Where("email = ?", requestBody.Student).First(&student).Error; err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "student not found"})
	}

	// suspend student
	student.Suspended = true
	if err := r.DB.Save(&student).Error; err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to suspend student"})
	}

	return context.Status(fiber.StatusNoContent).Send(nil)
}

// User story 4 POSTAPI
func (r *Repository) RetrieveForNotifications(context *fiber.Ctx) error {

	// validate req
	var req NotificationRequest
	if err := context.BodyParser(&req); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// separate the teacher and notification
	teacherEmail := req.Teacher
	notification := req.Notification

	// validate teacher
	var teacher models.Teacher
	if err := r.DB.Where("email = ?", teacherEmail).First(&teacher).Error; err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "teacher not found"})
	}

	// get students in the notification
	// use regex pattern to match email addresses starting with "@" and ending with ".com"
	pattern := `@[^\s]+\.com`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(notification, -1)

	// validate students in notification
	var validStudents []string
	var invalidStudents []string
	var count int64

	for _, match := range matches {
		match = strings.Replace(match, "@", "", 1)

		// count the number of rows with the given email
		if err := r.DB.Model(&models.Student{}).Where("Email = ?", match).Count(&count).Error; err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while validating students"})
		}

		if count == 0 {
			invalidStudents = append(invalidStudents, match)
		} else {
			validStudents = append(validStudents, match)
		}
	}

	// check if any invalid students found
	if len(invalidStudents) > 0 {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid students in notification", "invalid_students": invalidStudents})
	}

	// add all students tagged to the given teacher
	var teacherStudents []models.Student
	if err := r.DB.Table("students").
		Joins("JOIN teacher_students ON students.id = teacher_students.student_id").
		Joins("JOIN teachers ON teacher_students.teacher_id = teachers.id").
		Where("teachers.email = ?", teacherEmail).
		Find(&teacherStudents).Error; err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while finding students tagged to teacher"})
	}

	// merge students tagged with students in notifcation
	for _, student := range teacherStudents {
		found := false
		for _, validStudent := range validStudents {
			if validStudent == student.Email {
				found = true
				break
			}
		}
		if !found {
			validStudents = append(validStudents, student.Email)
		}
	}

	// remove students that are suspended
	var nonSuspendedStudents []string
	for _, studentEmail := range validStudents {
		var student models.Student
		if err := r.DB.Where("Email = ?", studentEmail).First(&student).Error; err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error while finding suspended students"})
		}

		if !student.Suspended {
			nonSuspendedStudents = append(nonSuspendedStudents, studentEmail)
		}
	}

	// return remaining students
	return context.Status(fiber.StatusOK).JSON(fiber.Map{"recipients": nonSuspendedStudents})
}

// routes
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", r.Register)
	api.Get("/commonstudents", r.GetCommonStudents)
	api.Post("/suspend", r.SuspendStudent)
	api.Post("/retrievefornotifications", r.RetrieveForNotifications)

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
