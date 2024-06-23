# NOTES

## Design

For the requirement a backend application with will implemented using Go. This backend will expose a RESTful API.

## Public API:

- `POST /users` -> user sign-up
- `POST /tokens` -> user's login

## Protected API:

- `GET /assets` -> get all assets
- `GET /assets/:assetId` -> get asset details
- `PATCH /assets/:assetId` -> patch asset's description
- `GET /assets/audiences` -> get all audiences
- `GET /assets/charts` -> get all charts
- `GET /assets/insights` -> get all insights
- `GET /users/:userId/favorites` -> get all user's favorites
- `POST /users/:userId/favorites` -> set user's favorite (add to favorites)
- `DELETE /users/:userId/:favoriteId` -> delete user's favorite

## Models:

- User (passwords will be hashed using `bcrypt`)
- Asset (`polymorphic` has-one relationship with Audience, Chart, Insight)
- Audience (pseudo `enum`s should enforce properties like `gender` and `ageGroup`)
- Chart
- Insight
- Favorite

Asset's foreign key will be the combination of `RelatedID` with `RelatedType` -> `RelatedType` could either be `charts` or `insights` or `audiences`

## Various Decisions

- To support efficiently a large amount of data, pagination should be used for all the relevant entities. The pattern should rely on query parameters of `page` and `limit`

- Caching layer should be introduced between server and database improve read operations and lift some height from the DB. The most potentially demanding query is to get user's favorites. Favorites could be stored in cache (Redis) using a key like `user:<userId>:favorites`.  
  Additionally if we want to enhance this to support pagination then the cache key could become `user:<userId>:favorites:page:<pageNumber>:limit:<limitSize>`. However this is a bit problematic as if it is important to allow users to change frequently the page size the above key becomes useless. In this case we need to use `ZRANGE` (is we are storing sorted sets) or `LRANGE` (if we are storing lists) and the overall complexity of using the cache increases.
  Invalidation of cache should happen when a user adds/removes favorites

* Seeding script should populate records in database

### Backend (Golang with Gin and Gorm)

#### Folder Structure:

- cmd/api -> `main.go`, server start
- internals/controllers -> application's controllers
- internals/database -> connection to DB
- internals/middlewares -> mainly auth logic for protected endpoints
- internals/models -> application's models
- internals/server -> init Gin server and application's routes declaration
- internals/services -> application's logic used from controllers
- internals/serializers -> map external entities to internal structures
- internals/utils -> small functions for handling hashing, tokens
- internals/helpers -> small reusable functions that interact with data layer

### Backend (React app using Vite and Antd for the UI)

Simple frontend which provides a login page and a dashboard

### Docker

Two containers, one for the frontend and one for the backend will be provided. Furthermore, compose files should make the boot of the application a breeze :)

### Quick Start

From the root of the cloned repo, one should execute
`docker compose up` or `docker-compose up` depending on the version of `docker` installed on user's OS

## TODOs

- [x] implement dummy frontend
- [x] implement backend models
- [x] implement backend API
- [x] implement proper frontend
- [x] write OpenAPI spec
- [ ] introduce a caching layer for caching users' favorites. This will improve response times of the main query which is responsible to fetch all the favorites of a logged-in user based on `userId`.
- [ ] write unit tests for the backend to thoroughly test the behavior of controllers and services
- [ ] implement validations of inputs in both frontend and backend
