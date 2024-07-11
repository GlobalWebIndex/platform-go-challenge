## Stivaktakis Giorgos - GWI test

This is a test for the GWI position. The task details can be found in the [task](task.md) file.

### Pre-requisites

- To init DB and the server you need to create a `.env` file in the root directory 
it should be something like this:
```shell
FIBER_PORT=3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=gwi
DB_PASS=gwi
DB_NAME=gwi-db
```

### To run the server and the DB

```shell
docker-compose up -d
```

### Endpoints
base url: `http://127.0.0.1:3000`

- `GET /assets/:userId` - Get all the assets of a user (this asset struct contains the asset id, the user id, the type of the asset and the id of the object, also the description of the asset)
- `GET /assets/objects/:userId` - Get all assets as objects of a user (this returns a responseAsset struct you can find it in the [response.go](/domain/response.go) file)
- `DELETE /assets/:assetId` - Delete an asset by id
- `POST /assets/` - Add an asset (it requires a json body with the asset struct)
- `PUT /assets/` - Update an asset (it requires a json body with the asset struct)


### The project structure

- main.go - The entry point of the application
- server - The server setup and the endpoints also handles the recieving of the requests
- repositories - The database queries as well as the connection to the DB
- domain - The structs of the application
- services - The core logic of the application
- ports - The interfaces of the application

### Extras

- The project uses the Fiber framework for the server
- The project uses GORM for the ORM and postgres as the DB
- The project uses the godotenv for the env variables
- I didnt used any concarrency in the project because I didnt see the need for it (maybe when we get the Assets from the DB we could have a different go routine for each asset but I didnt see the need for it in this case)

Thank you for the opportunity!