# GlobalWebIndex Engineering Challenge

This is an implementation of the [GWI go challenge](https://github.com/GlobalWebIndex/platform-go-challenge) by [Anastasopoulos Spyros](http://github.com/anastasop).

## Application Design

The challenge is given in a very general way, to allow implementations set their path. I decided to make an implementation that demostrates go usage. The code includes concurrency with contexts, an http server with graceful shutdown and a concurrent http client pool for fetching assets.

I split the application in 2 services. The asset server, that serves the assets, charts, insights and audiences, and the favorites server which is a CRUD service for favorites. You can start each service in a different machine or port, or both in one instance. The point is that the favorites server does not store assets, it must call the asset server to fetch them.

### Data Design

The data design is very simple. All assets are stored in static files. The _production_ one is `assets.json` which includes 3000 mock assets. A smaller example is `testdata/test-assets.json` with only 60 assets used for testing. The file is bundled with the binary using `go:embed`.

The favorites are also stored in a json file by default `./dbfav.json`. The favorites are preserved in the file between server restarts.

Everything, including assets and favorites, is stored in memory, there is no SQL database. This is a decision. The mapping of the entities, as described at least in the challenge is straightforward and maps directly to major ORMs, therefore i decided not to focus on data design.

### Security

It uses [JsonWebTokens](https://jwt.io). Users should login to get the token and use it to access the favorites and assets services. The favorites service propagates the user token to the assets service when needed. User names and passwords are stored in the text file `users.txt` and the JWT secret in `secret.txt`. Both are bundled with the binary using `go:embed`.

```
POST /login expects a form email=&passwd and return 200 OK with the token. Caller should store it and use it. Expires after 8 hours.
GET  /verify returns 200 OK if the token is accepted. For testing.
```

### API Design

The API design is straightforward. It is implemented using custom go code, structs and http, no OpenAPI or something similar.

Assets
```
GET /assets/charts/{id}     returns 200 OK with the chart or 404 NotFound
GET /assets/insights/{id}   returns 200 OK with the insight or 404 NotFound
GET /assets/audiences/{id}  returns 200 OK with the audience or 404 NotFound
```

Favorites
```
GET    /favorites      returns 200 OK with all the favorites including assets
POST   /favorites      creates a new favorite, returns 201 Created and includes a Location header with the url of the favorite
GET    /favorites/{id} get the favorite, returns 200 OK. Response does not include the asset
PUT    /favorites/{id} updates the favorite, returns 200 OK and the new entity
DELETE /favorites/{id} deletes the favorite, returns 204 NoContent
```

## Usage

This is a standalone go application. You can build with
```
go build
```

The usage is
```
% ./gwifav -h
Usage of ./gwifav:
  -a string
        assets service host:port (default "localhost:8080")
  -c string
        assets service host:port for the fetch pool. Use with different machines
  -d string
        favorites file (default "./dbfav.json")
  -f string
        favorites service host:port (default "localhost:8080")
  -w duration
        asset server delay per request. Use for testing.
```

can do a simple run with
```
./gwifav
```

this will launch an http server listening at `localhost:8080` that serves favorites and assets. You can use the API using curl. There are tools you can install with `. ./tools.sh`. A sample demo session is:

```
% go build
% ./gwifav & # or start it in another window to separate logs and output
% . ./tools.sh
% gwifav-login `head -1 users.txt` # stores token in ./token.jwt
% gwifav-add insights/1
/favorites/1
% gwifav-add charts/2
/favorites/2
% gwifav-add audiences/4
/favorites/3
% gwifav-get 2
{"id":2,"assetKey":"charts/2","description":"description of charts/2"}
% gwifav-getall | jq .
{
  "total": 3,
  "favorites": [
    {
      "id": 3,
      "assetKey": "audiences/4",
      "description": "description of audiences/4",
      "IncludesAsset": true,
      "Asset": {
        "gender": "Females",
        "ageFrom": 1,
        "ageTo": 57,
        "country": "TUV",
        "dailyHoursOnSocialMedia": 10,
        "numPurchasesLastMonth": -1
      }
    },
    {
      "id": 1,
      "assetKey": "insights/1",
      "description": "description of insights/1",
      "IncludesAsset": true,
      "Asset": {
        "text": "rc: Usage: whatis name ...: fd out of range or not open"
      }
    },
    {
      "id": 2,
      "assetKey": "charts/2",
      "description": "description of charts/2",
      "IncludesAsset": true,
      "Asset": {
        "title": "Hydrogen/Iodine with 3 points",
        "axes": [
          "Hydrogen",
          "Iodine"
        ],
        "data": [
          0.1356436,
          0.4782989,
          0.35305363
        ]
      }
    }
  ]
}
```

## Testing

There are unit tests for the important things with coverage larger than 90%. The functionality not covered with unit tests is

* `main.go` but this is only setup code
* graceful shutdown and signal handling.

You can test the later with the `-w` flag which sets a delay to the asset service handlers. The cases are:
- start the server with a big delay, call `gwifav-getall` and see that it timeouts.
- start the server with a big timeout, call `gwifav-getall` and after a while kill the client. In the server logs you will see that the cancellation is propagated to the fetch pool and stops.
- start the server with a big timeout, call `gwifav-getall` and after a while kill the server. In the server logs you will see that the cancellation is propagated to the fetch pool, the client get a partial response, the server shuts down gracefully and the favorites are preserved in the file.

## Caveats

The separation between the favorites and the assets services raised the question what happens on failures.

The favorites case are simple. When the service fails, we can restart it.

The assets service however is different. The interesting case is what happens if it fails while it is used by the asset fetch pool. In this case the pool can retry but there is no guarantee that the service will be up before the timeout. For this implementation i decided that the favorites server returns the favorite without the asset. This way the client always gets a response and can check the number of expected vs actual assets to try again. Note that for this to work best, the `/favorites` endpoint needs paging but this is not implemented.
