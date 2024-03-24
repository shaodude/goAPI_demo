package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert" // add Testify package
)

func TestRegisterRoute(t *testing.T) {
	testCases := []struct {
		description  string
		route        string
		expectedCode int
	}{
		// first test case
		{
			description:  "get HTTP status 200",
			route:        "/api/register",
			expectedCode: 200,
		},
		// second test case
		{
			description:  "get HTTP status 404, when route does not exist",
			route:        "/api/reegiister",
			expectedCode: 404,
		},
	}

	app := fiber.New()

	for _, testCase := range testCases {
		// create  http request
		req := httptest.NewRequest("GET", testCase.route, nil)

		resp, _ := app.Test(req, 1)

		assert.Equalf(t, testCase.expectedCode, resp.StatusCode, testCase.description)
	}
}
