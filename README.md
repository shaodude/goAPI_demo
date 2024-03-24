# Golang API Assessment by Shao Yi Goh

## Abstract
While my working exposure in backend development primarily involves PHP using the Laravel Framework, I've ventured into learning Golang specifically for this assessment. Given the time constraints, the project's capabilities are currently limited. However, potential enhancements are outlined below in the "Future Enhancements" section.

# Future Enhancements
1) Refactor code and restructure project directories to improve code readbility, maintainability, and scalability
2) Rectify unit tests and expand test coverage with additional test cases.
3) Explore the utilization of complex queries in API logic to optimize runtime and memory usage. Further testing will be necessary to validate these improvements.

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

## Run Project

```sh
go run main.go
```
