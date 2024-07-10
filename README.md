# Table of contents

 1. [GlobalWebIndex Engineering Challenge](#globalwebindex-engineering-challenge)
 2. [Features](#features)
 3. [Prerequisities](#prerequisities)
 4. [Installation](#Installation)
 5. [API Endpoints](#api-endpoints)
 6. [Running Tests](#running-tests)
 7. [Useful Resources](#useful-resources)
 8. [Contact](#contact)


# GlobalWebIndex Engineering Challenge

## Introduction

This challenge is designed to give you the opportunity to demonstrate your abilities as a software engineer and specifically your knowledge of the Go language.

On the surface the challenge is trivial to solve, however you should choose to add features or capabilities which you feel demonstrate your skills and knowledge the best. For example, you could choose to optimise for performance and concurrency, you could choose to add a robust security layer or ensure your application is highly available. Or all of these.

Of course, usually we would choose to solve any given requirement with the simplest possible solution, however that is not the spirit of this challenge.

## Challenge

Let's say that in GWI platform all of our users have access to a huge list of assets. We want our users to have a peronal list of favourites, meaning assets that favourite or “star” so that they have them in their frontpage dashboard for quick access. An asset can be one the following
* Chart (that has a small title, axes titles and data)
* Insight (a small piece of text that provides some insight into a topic, e.g. "40% of millenials spend more than 3hours on social media daily")
* Audience (which is a series of characteristics, for that exercise lets focus on gender (Male, Female), birth country, age groups, hours spent daily on social media, number of purchases last month)
e.g. Males from 24-35 that spent more than 3 hours on social media daily.

Build a web server which has some endpoint to receive a user id and return a list of all the user’s favourites. Also we want endpoints that would add an asset to favourites, remove it, or edit its description. Assets obviously can share some common attributes (like their description) but they also have completely different structure and data. It’s up to you to decide the structure and we are not looking for something overly complex here (especially for the cases of audiences). There is no need to have/deploy/create an actual database although we would like to discuss about storage options and data representations.

Note that users have no limit on how many assets they want on their favourites so your service will need to provide a reasonable response time.

A working server application with functional API is required, along with a clear readme.md. Useful and passing tests would be also be viewed favourably

It is appreciated, though not required, if a Dockerfile is included.

## Submission

Just create a fork from the current repo and send it to us!

Good luck, potential colleague!


# Features
- **Asset Management**: Users can manage their favorite assets including adding, retrieving, editing, and deleting favorites.
- **Authentication**: Secure endpoints using JWT-based authentication.
- **Pagination**: Efficiently handle large sets of data with pagination.
- **Caching**: Improve performance with in-memory caching for frequently accessed endpoints.  -- not implemented
- **Data Management**: Currently, data is managed in-memory. Future enhancements may include MongoDB or AWS-based data persistence.
- **Tests**: Some tests were implemented althought a lot of cases are still missing
- **Docker File**: A docker file was created bases on tutorials.No testing or knowledge if it is correct...although it looks correct...sorry


# Prerequisities
- **Go 1.22.4**: [Install Go](https://golang.org/doc/install)
- **Git**: Make sure Git is installed on your system. [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Docker](https://www.docker.com/)  

### Go Modules
- `github.com/davecgh/go-spew v1.1.1`
- `github.com/golang-jwt/jwt v3.2.2+incompatible`
- `github.com/joho/godotenv v1.5.1`
- `github.com/labstack/gommon v0.4.2`
- `github.com/mattn/go-colorable v0.1.13`
- `github.com/mattn/go-isatty v0.0.20`
- `github.com/pmezard/go-difflib v1.0.0`
- `github.com/valyala/bytebufferpool v1.0.0`
- `github.com/valyala/fasttemplate v1.2.2`
- `golang.org/x/crypto v0.22.0`
- `golang.org/x/net v0.24.0`
- `golang.org/x/sys v0.19.0`
- `golang.org/x/text v0.14.0`
- `golang.org/x/time v0.5.0`
- `gopkg.in/yaml.v3 v3.0.1`
- `github.com/labstack/echo/v4 v4.6.1`

 - setup env file to the root, althought for the purpose of demonstarion the .env file is pushed to git repository


# Installation

``` bash
   #Clone the repository: 
   git clone https://github.com/yourusername/platform-go-challenge.git
   #browse to tjhe directory
   cd platform-go-challenge
   #install depentancies
   go mod tidy
   #setup env file to the root with the following, 
   #althought for of the purpose the .env file is pushed to git repository
   GOPROXY=https://proxy.golang.org,direct

   go run server.go
```

# API Endpoints

The server defines the following RESTful endpoints:
``` bash
POST /users/{id}/favorites: Adds a favorite to a user favotites list.
#e.g. curl comamnd 
curl -X POST "http://localhost:1323/users/1/favorites" \
-H "Authorization: gwi-token-12345" \
-H "Content-Type: application/json" \
-d '{"Type":0,"Description":"Interested in sports","Chart":{"Title":"How interested are you in Sports?","XAxis":"Interest","YAxis":"Amount",
"Data":{"Quite interested":32,"Very interested":31.9}}}'


GET /users/{id}/favorites: Retrieve all favorities of a user with a certain id.
#e.g curl comamnd 
curl GET "http://localhost:1323/users/1/favorites" -H "Authorization:gwi-token-12345"
#curl command to retreive all user favorites using pagination
curl -X GET "http://localhost:1323/users/1/favorites?page=1&limit=1" \
-H "Authorization: Bearer gwi-token-12345" \
-H "Content-Type: application/json"


PUT /users/{id}/favorites/{id}: Edits a user\'s favorite description 
curl -X PUT "http://localhost:1323/users/1/favorites/1a80e52a-f604-41e2-9aa2-4fd0f732e649" \
     -H "Authorization: gwi-token-12345" \
     -H "Content-Type: application/json" \
     -d '{"description": "Favorite Sport Trends"}'

DELETE /users/{id}/favorites/{id}: Delete a favorite from a user\'s list.
#e.g curl comamnd
curl -X DELETE "http://localhost:1323/users/1/favorites/1a80e52a-f604-41e2-9aa2-4fd0f732e649" \
-H "Authorization:gwi-token-12345"
```

![explorer_sU84BsvTO6](https://github.com/artemis13/platform-go-challenge/assets/4024511/73911c80-1bf6-4177-bced-d51afa102f79)

# Running Tests

To execute the automated tests for this system, run the following command in the project directory:
``` bash
go test -count=1 ./...
```
# Useful Resources
Some of the tutorials used to decide the framework and the modules that were used

1. [Exploring Golang REST API Frameworks](https://dev.to/xngwng/top-5-go-rest-api-frameworks-k0e)
2. Tutorial: [Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin)
3. Tutorial: [How to build an API using Go and Gorilla mux](https://dev.to/envitab/how-to-build-an-api-using-go-ffk)
4. Tutorial: [Build a Golang RESTful Stock API With the Echo Framework](https://betterprogramming.pub/intro-77f65f73f6d3)
5. [echo QuickStart](https://echo.labstack.com/docs/quick-start)
6. [echo CRUD](https://echo.labstack.com/docs/cookbook/crud)
7. [How to Dockerize a Golang Application?](https://www.geeksforgeeks.org/how-to-dockerize-a-golang-application/)

# Contact

Created by [Artemis Apostolou](https://github.com/artemis13) for the purposes of GWI Go Platrform Challenge
