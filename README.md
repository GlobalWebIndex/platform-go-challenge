# GlobalWebIndex Engineering Challenge

## Introduction

This challenge is designed to give you the opportunity to demonstrate your abilities as a software engineer and specifically your knowledge of the Go language.

On the surface the challenge is trivial to solve, however you should choose to add features or capabilities which you feel demonstrate your skills and knowledge the best. For example, you could choose to optimise for performance and concurrency, you could choose to add a robust security layer or ensure your application is highly available. Or all of these.

Of course, usually we would choose to solve any given requirement with the simplest possible solution, however that is not the spirit of this challenge.

## Challenge

Let's say that in GWI platform all of our users have access to a huge list of assets. We want our users to have a peronal list of favourites, meaning assets that favourite or "star" so that they have them in their frontpage dashboard for quick access. An asset can be one the following
* Chart (that has a small title, axes titles and data)
* Insight (a small piece of text that provides some insight into a topic, e.g. "40% of millenials spend more than 3hours on social media daily")
* Audience (which is a series of characteristics, for that exercise lets focus on gender (Male, Female), birth country, age groups, hours spent daily on social media, number of purchases last month)
e.g. Males from 24-35 that spent more than 3 hours on social media daily.

Build a web server which has some endpoint to receive a user id and return a list of all the user's favourites. Also we want endpoints that would add an asset to favourites, remove it, or edit its description. Assets obviously can share some common attributes (like their description) but they also have completely different structure and data. It's up to you to decide the structure and we are not looking for something overly complex here (especially for the cases of audiences). There is no need to have/deploy/create an actual database although we would like to discuss about storage options and data representations.

Note that users have no limit on how many assets they want on their favourites so your service will need to provide a reasonable response time.

A working server application with functional API is required, along with a clear readme.md. Useful and passing tests would be also be viewed favourably

It is appreciated, though not required, if a Dockerfile is included.

## Submission

Just create a fork from the current repo and send it to us!

Good luck, potential colleague!

---

## Implementation

Favourites API implemented in **Go** with **hexagonal architecture** (ports and adapters). The API is defined in **`swagger.yaml`** (OpenAPI 3.0); server types and routes are generated with **oapi-codegen**. Handlers live in **`http/`**, use cases in **`favourites/application`**, and adapters for persistence and the asset catalog in **`favourites/adapters`** and **`assets/adapters`**. Favourites are stored in memory; the asset catalog is in-memory with pre-seeded data. JWT (Bearer) is used for auth; list endpoint is paginated.

---

## Run locally

**Prerequisites:** Go 1.26+ (or use Docker below).

From the repo root:

```bash
go run ./cmd/favourites-api
```

Server listens on **http://localhost:8080**. API is under **`/v1`** (e.g. `POST /v1/auth/token`, `GET /v1/users/{userID}/favourites`). **Swagger UI** at **http://localhost:8080/swagger/**.

Optional: set **`APP_ENV`** to `dev` (default), `staging`, or `production` to load config from `config/<env>.yaml`. Set **`JWT_SECRET`** for token signing (or use RSA keys via **`JWT_PRIVATE_KEY_PATH`**); if unset, a dev default is used (not for production). Set **`LOG_LEVEL`** to `debug`, `info` (default), `warn`, or `error` to control log verbosity (uses `log/slog`).

**Tests:**

```bash
go test ./...
```

---

## Run with Docker

**Build:**

```bash
docker build -t favourites-api .
```

**Run:**

```bash
docker run -p 8080:8080 -e JWT_SECRET=your-secret-here favourites-api
```

**With Docker Compose:**

```bash
docker-compose up --build
```

Optional: set **`JWT_SECRET`** in the environment or in a `.env` file. The image uses a multi-stage build and runs as a non-root user.
