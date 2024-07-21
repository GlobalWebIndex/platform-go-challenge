# Asset API

This is just another API, or is it?

## Introduction

This API is designed to allow users to interact with Asset data (stored in a MongoDB).

### Asset Types

There are three asset types in this API at the moment;
Charts, Insights and Audiences.

Each of them has its own characteristics but the can all be favoured by any user,
and fetched under the same favourites list.

### Roles

The users in this concept have roles assigned to them.
The roles determined which endpoints they can use from the API, 
and the information displayed in each.

There are two roles currently existing; admin and user.

## Usage

### Running the API

There are two ways to run this application;

1. Docker Compose (includes a MongoDB)
2. Go Run (along with a MongoDB).

Specifically:

1. For both approaches we need an `.env` file in the root directory of the project
   (same level with this `README.md` and `main.go` file). Please don't commit this file.
    ```ini
    MONGODB=mongodb://localhost:27017
    DB_NAME=favourites
    JWT_SECRET_KEY=SOOOOOO_SECRET_KEY
   ```
   where:
    * `MONGODB` is the URI of the MongoDB
    * `DB_NAME` is the name of the database to be used by the application
    * `JWT_SECRET_KEY` is the Secret Key used to sign the Json Web Token
      used during authentication and authorization
2. A Postman installation (see [official instructions here](https://www.postman.com/downloads/))
    1. Import the [GoFavourites.postman_environment.json](GoFavourites.postman_environment.json).
        Once imported update the `JWT_SECRET_KEY` to match the `.env` one.
    2. import [GoFavourites.postman_collection.json](GoFavourites.postman_collection.json)
3. Then we can either:
   1. (Recommended) Run with Docker Compose
      1. Have Docker and Docker Compose installed
         (see [official instructions here](https://docs.docker.com/compose/install/)).
      2. Then open a terminal, navigate in this directory and
      3. run `docker compose up`
         4. Depending on your system and version may be `docker-compose up`
   2. Run using Go Run for which we need to
      1. have set `GOROOT`, `GOPATH`
      2. Have a working MongoDB running (and the URI in `.env` file)
      3. Then open a terminal, navigate in this directory
      4. Run `go run main.go`


##  Documentation
### Postman Collection Usage
Since swagger turned out to be not that practical for displaying purposes.
A set of environment and postman collection along with some docs should do the trick.

The requests play around with environment variables for ease of using the entire collection.

#### Creating Assets via Admin user
First we need an admin user, 
to insert all the asset data 
available to the rest of the users to favourite,
thus we should create one.

1. Create an admin user: `User/Admin/Admin User SignUp` POST request
2. Login with this user: `User/Login` POST request
3. To add assets copy values from [INSIGHT_MOCK_DATA.json](INSIGHT_MOCK_DATA.json),
[CHARTS_MOCK_DATA.json](CHARTS_MOCK_DATA.json) and
[AUDIENCE_MOCK_DATA.json](AUDIENCE_MOCK_DATA.json) and paste it in requests:
   1. `User/Admin/Add Insights (bulk)` POST request
   2. `User/Admin/Add Charts (bulk)` POST request 
   3. `User/Admin/Add Audiences (bulk)` POST request 
4. Admin can also bulk add users and see all users:
   1. `User/Admin/Add Users (bulk)` POST request
   2. `User/Admin/List All Users` GET request

#### Viewing public assets 
After their creation assets are available to all, 
i.e. no authentication is required to view them.
1. To view them:
    1. `Public/List All Charts` GET request
    2. `Public/List All Insights` GET request
    3. `Public/List All Audiences` GET request
2. To view a specific Asset by its ID:
   1. `Public/Chart by ID` GET request
   1. `Public/Insight by ID` GET request
   1. `Public/Audience by ID` GET request

#### Creating a user and adding to favourites
1. Create a user: `User/User SignUp` POST request
2. Login with this user: `User/Login` POST request
3. Now user can add to Favourites: 
   1. `User/Add Favourite Chart` POST request
      2. if the flow has been followed as in this file, 
      a random ID from the Add requests should be in place.
   2. `User/Add Favourite Insight` POST request
   2. `User/Add Favourite Audience` POST request
4. Now that user has favourites (so that we don't get an empty list):
   1. we can list them `User/List User's Favourites` and then
   1. Fetch one by its ID `User/Favourite by ID`
   2. We can also remove it `User/Remove Favourite`.
      However this will empty the variable of the ID so try to add more
        and list them to populate it again.
5. Lastly the users can logout by `User/Logout` request

Have fun!

## Other Info

1. Based on [Instructions of the challenge](Instructions.md)
2. Let me know your thoughts!
