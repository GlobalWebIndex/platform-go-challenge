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

# Favourites REST API (Go)

## Introduction

A simple in-memory REST API written in Go (version 1.18.2) that allows users to manage a list of favourite assets.
Supported asset types include:

- **Chart**
- **Insight**
- **Audience**

There is **no database**. All data lives in memory using the following data structure:

map[string]map[string]Asset{}

(This is a map of maps containing the lists of favourites for all users of the system. UserID -> map of (AssetID -> Asset))


---

## Features

- Get the list of all favourites for a user (GET)
- Add a chart, insight, or audience to favourites (POST)  
- Update (partially) an existing asset (PATCH)  
- Delete an existing asset (DELETE)  
- In-memory storage (no database)
- Table-driven unit tests
- Minimal Dockerfile included

---

## API Endpoints

### Get all favourites for a user

GET /users/{user_id}/favourites

### Add a favourite (chart, insight or audience)

POST /users/{user_id}/favourites/{asset_type}

Specifically:

POST /users/{user_id}/favourites/chart
POST /users/{user_id}/favourites/insight
POST /users/{user_id}/favourites/audience

### Update (partially) an existing asset

PATCH /users/{user_id}/favourites/{asset_id}

### Delete an asset

DELETE /users/{user_id}/favourites/{asset_id}

### Auxiliary endpoint. Create a new user (just for the sake of this user having a list of favourites)

---

## Data Model

### Chart

```json
{
  "description": "Example chart",
  "data": {
    "title": "Chart Title",
    "x_axis": "X Label",
    "y_axis": "Y Label",
    "data": [1, 2, 3]
  }
}
```

### Insight

```json
{
  "description": "Example insight",
  "data": {
    "text": "Insightful text"
  }
}
```

### Audience
```json
{
  "description": "Example audience",
  "data": {
    "gender": "Male",
    "birth_country": "GR",
    "age_groups": ["18-24"],
    "daily_hours_on_social_media": 3,
    "purchases_last_month": 2
  }
}
```

---

## Build

go build -o api .

## Run locally (listens on localhost:8080)

./api

## Run tests (verbose)

go test

---

## User missing behaviour
1. GET: If a user does not exist, return 404.
2. POST: If a user does not exist, return 404. If a user does not exist, the list of favourites we are trying to add an asset in does not exist as well.
3. PATCH: If a user does not exist, return 404. It does not make sense to try to partially update an asset from the favourites list of a non-existing user.
4. DELETE: If a user does not exist, return 404. If a user does not exist, it does not make sense to try to delete an asset from the favourites list of a non-existing user.

---

## Build Docker image

docker build -t favourites-api .

## Run container

docker run -p 8080:8080 favourites-api