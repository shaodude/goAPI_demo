# Golang API Assessment
A backend application developed using Golang for backend code and PostgreSQL for database. The application is a part of a system for administrative functions by teachers for their students.

# Future Enhancements:
  - Refactor code and restructure project directories to improve code readbility, maintainability, and scalability
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
