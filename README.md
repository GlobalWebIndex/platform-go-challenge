# GlobalWebIndex Engineering Challenge

## Introduction
This is my implementation for GWI Engineering Challenge.

## Logic
In this implementation I developed a dashboard feature that enables a user to add various assets to their personal dashboard of assets.  
The user can see, add and remove assets and can also add/edit a description for each asset in their dashboard.

## How to use
Run  
`docker-compose build`  
`docker-compose up`

This will build and deploy localy 3 containers.  
* app - main server for go
* sqlserver - main DB for the app
* sqlserver-test - secondary DB for integration tests

When the app is deployed a job will seed the database with some test data so you can use the API.  
To use the API you can find the Postman collection in the `docs` directory.

## API
### Get Token
This endpoind provides a unique token for the user to be used in the required endpoints that edit their dashboard.  
`GET '0.0.0.0:6060/user/{user_id}/token/get'`
### Get Dashboard
This endpoint fetches a paginated list of assets tha are added to the user's dashboard.  
Note that this fetches x assets from every asset type.  
`GET '0.0.0.0:6060/user/{user_id}/dashboard/list?page=1&per_page=10'`  
### Add Asset To DashBoard
This endpoint adds an asset to a user's dashboard.  
*Needs Token.    
`PUT '0.0.0.0:6060/user/{user_id}/dashboard/asset/add'`  
### Remove Asset From Dashboard
This endpoint removes an asset from a users dashboard.  
*Needs Token.  
`PUT '0.0.0.0:6060/user/{user_id}/dashboard/asset/remove'`  
### Edit Asset Description
This endpoint edits the description of an asset in a user's dashboard.  
*Needs Token.  
`PATCH '0.0.0.0:6060/user/{user_id}/dashboard/asset/edit'`
### List Assets
This endpoint lists all available assets.  
`GET '0.0.0.0:6060/asset/list'`

### How to use the token.
The token can be used in the headers part of the request using `Authorization` key and the value must be set like `Bearer <token-string>`

## Technical Info
The app was implemented in go version 1.18 .  
The mysql DB is in version 8.0.23 .  
Default host for app is 0.0.0.0 .   
Default port for app container is :6060 .



