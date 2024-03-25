# Golang API Assessment by Shao Yi Goh
A backend application developed using Golang for backend code and PostgreSQL for database. The application is a part of a system for administrative functions by teachers for their students.

# Abstract:
While my working exposure in backend development primarily involves PHP using the Laravel Framework, I've ventured into learning Golang specifically for this assessment. Given the time constraints, the project's capabilities are currently limited. However, potential enhancements are outlined below in the "Future Enhancements" section.

# Future Enhancements:
  - Refactor code and restructure project directories to improve code readbility, maintainability, and scalability
  - Rectify unit tests and expand test coverage with additional test cases.
  - Explore the utilization of complex queries in API logic to optimize runtime and memory usage. Further testing will be necessary to validate these improvements.

# Prerequisites:
  - Go (developed on version 1.22.1)
  - PostgreSQL database (developed on version 16.2)

# Setup:

## Install dependencies

```sh
go mod tidy
```

## Set up the environment variables

Create a .env file in the root directory and add the following environment variables:

```sh
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
DB_SSLMODE=disable
```

Replace the values with your database credentials.

Note that you may need to create a new DB if there is not one available. This can be done using a PostgreSQL management tool like PGadmin4.

## Run Project

```sh
go run main.go
```
