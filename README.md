# GWI Assets Platform

## Overview

GWI Platform is a robust web application that manages user data, assets, and favorites.
It provides RESTful API endpoints for various operations including user management, asset handling, and user favorites.

## Features

- User management (add)
- Asset management (add, delete, modify, get)
- User favorites handling
- Database integration with MySQL
- Middleware for logging and user authentication
- Graceful server shutdown

## Project Structure

gwi-platform/
├── README.md
├── app
│   ├── Dockerfile
│   ├── assets
│   │   ├── asset.go
│   │   ├── asset_test.go
│   │   ├── audience.go
│   │   ├── audience_test.go
│   │   ├── chart.go
│   │   ├── chart_test.go
│   │   ├── insight.go
│   │   ├── insight_test.go
│   │   └── queries.go
│   ├── auth
│   │   └── auth.go
│   ├── database
│   │   └── database.go
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   ├── asset_handlers.go
│   │   ├── favorite_handlers.go
│   │   ├── handlers.go
│   │   ├── queries.go
│   │   └── users_handlers.go
│   ├── html
│   │   ├── api_docs.html
│   │   └── index.html
│   ├── logs
│   │   └── app.log
│   ├── main.go
│   ├── models
│   │   └── models.go
│   ├── server
│   │   ├── app.go
│   │   └── app_test.go
│   └── utils
│       └── utils.go
├── db
│   └── init.sql
├── docker-compose.yml
└── test
├── Dockerfile
├── go.mod
└── main.go


## Prerequisites

- Go 1.20 or higher
- MySQL 8.0 or higher

## Setup

1. Clone the repository:
- git clone https://github.com/your-username/gwi-platform.git
- cd gwi-platform

2. Set up the database:
- Create a MySQL database
- Update the database connection details in `database/database.go`

3. Install dependencies:
- cd app
- go mod tidy

4. Build the application:
- go build -o gwi-platform main.go

5. Run the application:
- ./gwi-platform

## API Endpoints

### User Management

- `POST /user/add`: Add a new user

### Asset Management

- `POST /assets/add`: Add a new asset
- `DELETE /assets/delete`: Delete an asset
- `PUT /assets/modify`: Modify an existing asset
- `GET /assets/get`: Retrieve an asset

### User Favorites

- `POST /user/favorite/add`: Add a favorite for a user
- `GET /user/favorites/`: Get detailed favorites for a user

### Other Endpoints

- `GET /`: Home endpoint
- `GET /docs`: Documentation endpoint
- `GET /ping`: Health check endpoint

## Authentication

The application implements a basic authentication system through the `UserStatusAuth` middleware,
located in the `auth` package. This middleware serves as a secondary authentication layer,
focusing on user status verification rather than primary authentication,
which should be another application eg a CAS.

## Logging

Request and response logging is implemented using the `LogHandler` middleware in the `utils` package.

## Database

The application uses MySQL for data persistence. Database operations are handled in the `database` package.

### Docker Compose

The `docker-compose.yml` file defines the services needed to run the application:

- `app`: The main application service
- `db`: MySQL database service
- `test`: Service for running integration tests

To start the services:
- docker-compose up -d

This will run the integration tests defined in `test/main.go` against the running application.

The integration tests cover:
- User creation
- Asset management (creation, retrieval, modification, deletion)
- User favorites operations

These tests help ensure that all parts of the system are working together as expected in a new environment.


## Future Work

### Enhance Test Coverage
1. **Unit Tests**
   - Implement comprehensive unit tests for all packages

2. **Integration Tests**
   - Expand current integration test suite to cover more complex scenarios
   - Introduce performance tests to ensure scalability
   - Rewrite the current integration test suite using a BDD (Behavior-Driven Development).
     eg
     ```gherkin
         Feature: User Favorites Management

         Scenario: User adds an asset to favorites
            Given a user {testUser} exists in the system
            And an asset {testAsset} of type "CHART" exists
            When the user adds the asset to their favorites
            Then the asset should appear in the user's list of favorites
            And the total count of user's favorites should increase by 1


### Performance Optimization
1. **Database Optimization**
   - Introduce caching mechanisms for frequently accessed data

### Documentation
3. **API Documentation**
   - Provide interactive API exploration tools (e.g., Swagger)
