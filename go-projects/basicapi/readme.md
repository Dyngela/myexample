# basic api

This is a basic rest api that uses Gin router and gorm for database operations. It is based on postgres database.
But can be easily modified to use any other database thanks to gorm.

Project contains :

- CRUD operations for an entity
- API documentation using Swagger
- A generic repository pattern for database operations
- A config file for environment variables
- A logger for logging
- A simple dockerfile for containerization associated with a docker-compose file, so you do not even have to create a database neither tables nor insert data. Some dummy data are already inserted in the database.

## How to run

You can modify the .env file to your liking. Default value are already set for the database connection. Server will run on port 8080 by default.

You can access the swagger documentation at http://localhost:8080/swagger/index.html (if you are running the server on your local machine with default port and in dev mode)

```bash

### Using docker-compose

Make sure you have docker installed and running on your machine

```bash
docker-compose up
```

### Using go

```bash
go run main.go
```

## TODO

- [ ] Add redis caching
- [x] Add middleware for simple api-key check
- [x] Add middleware for rate limiting
- [x] Add middleware for logging
- [ ] Add unit tests
- [ ] Add integration tests