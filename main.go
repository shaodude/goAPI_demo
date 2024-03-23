package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/shaodude/goAPI_demo/models"
	"github.com/shaodude/goAPI_demo/storage"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Student struct {
	Email string `json:"email"`
}

type Teacher struct {
	Email string `json:"email"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Register(context *fiber.Ctx) error {
	// student := Student{}

	// err := context.BodyParser(&student)

	// if err != nil {
	// 	context.Status(http.StatusUnprocessableEntity).JSON(
	// 		&fiber.Map{"message": "request failed"})
	// 	return err
	// }

	// err = r.DB.Create(&student).Error
	// if err != nil {
	// 	context.Status(http.StatusBadRequest).JSON(
	// 		&fiber.Map{"message": "could not register student"})
	// 	return err
	// }

	// context.Status(http.StatusOK).JSON(&fiber.Map{
	// 	"message": "student has been registered"})
	// return nil
	// Parse request body
	var requestBody struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}
	if err := context.BodyParser(&requestBody); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Fetch teacher from database
	var teacher models.Teacher
	if err := r.DB.Where("email = ?", requestBody.Teacher).First(&teacher).Error; err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "teacher not found"})
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

		// Associate student with the teacher
		if err := r.DB.Model(&teacher).Association("Students").Append(&student); err != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to register student"})
		}
	}

	// Return success response
	return context.SendStatus(fiber.StatusNoContent)
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Post("/register", r.Register)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
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
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	err = models.MigrateStudents(db)
	if err != nil {
		log.Fatal("could not migrate students db")
	}
	err = models.MigrateTeachers(db)
	if err != nil {
		log.Fatal("could not migrate teachers db")
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
