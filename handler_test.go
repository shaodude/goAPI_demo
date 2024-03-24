package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	// Create a new repository instance
	repo := &Repository{} // Replace with your actual repository initialization

	// Create a sample request body
	requestBody := map[string]interface{}{
		"teacher": "teacherken@gmail.com",
		"students": []string{
			"studentjon@gmail.com",
			"studenthon@gmail.com",
		},
	}

	// Convert request body to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// Create a new Fiber app
	app := fiber.New()

	// Define a handler function for the `/api/register` endpoint
	app.Post("/api/register", func(ctx *fiber.Ctx) error {
		return repo.Register(ctx)
	})

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(requestBodyBytes))
	assert.NoError(t, err)

	// Perform the request using app.Test
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
