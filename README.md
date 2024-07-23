This application is designed to manage user favorites for various types of assets using a RESTful API. It provides endpoints to add, retrieve, update, and delete favorite assets for users. The application uses PostgreSQL for persistent storage.

ENDPOINTS :

GET Request    :    Get a list of favorite assets for a user with optional filtering and pagination.
POST Request   :    Add a new favorite asset for a user.
POST Request   :    Add multiple new favorite assets for a user.
DELETE Request :    Remove Favorite Asset: Delete a favorite asset for a user.
PUT Request    :    Update the description of a favorite asset.

Sample requests for each endpoint:

GET Request without filtering or pagination -> http://localhost:8080/favorites/user1

--------------------------------------------------------------------------------------------------------------

GET Request with filtering and pagination -> http://localhost:8080/favorites/user1?type=Chart&page=1&pageSize=10

Default values for pagination : page = 1 , pageSize = 10
With page we specify which page we want to retrieve and with pageSize we specify how many rows each page has.
So with page = 1 and pageSize = 10 we basically want to retrieve the first 10 rows.

--------------------------------------------------------------------------------------------------------------

POST Request to add a single asset -> http://localhost:8080/favorites/user1

body->json raw:

{
    "id": "testInsight",
    "type": "Insight",
    "description": "A Sample text",
    "data": {
        "text": "only 15% of the people in Greece watch One Piece"
    }
}

--------------------------------------------------------------------------------------------------------------

POST Request to add multiple assets -> http://localhost:8080/multiple/favorites/user1

body->json raw :

[
    {
        "id": "insight4",
        "type": "Insight",
        "description": "A Sample text",
        "data": {
            "text": "only 15% of the people in Greece watch One Piece"
        }
    },
    {
        "id": "insight5",
        "type": "Insight",
        "description": "An insightful text",
        "data": {
            "text": "only 15% of the people in Greece watch One Piece"
        }
    }
]

--------------------------------------------------------------------------------------------------------------

DELETE Request to remove an asset for a user -> http://localhost:8080/favorites/user1/chart1

--------------------------------------------------------------------------------------------------------------

PUT Request to edit the description of an asset for a user -> http://localhost:8080/favorites/user1/insight1

body->json raw:

{
    "description": "Updated description for the asset"
}

--------------------------------------------------------------------------------------------------------------

HOW TO RUN :

In order to run the app we need to have docker installed.

Then all we have to do is to run this command -> docker-compose --env-file db.env up --build

Since our 2 containers are up, we can use Postman to send requests like the ones I specified above.

NOTE :

I created a init.sql script which will be executed when we first run the app and it basically creates

the tables we need with some dummy data in order to test our endpoints.

**If we already have a another postgres instance in our machine we should terminate it in order for the script to run!!

TESTS :

In order to run the tests simply execute this command -> go test ./api
