

![gwi-logo](https://upload.wikimedia.org/wikipedia/commons/7/7c/GWI_Logo.webp)
# GlobalWebIndex Engineering Challenge

  

## Challenge

  

Let's say that in GWI platform all of our users have access to a huge list of assets. We want our users to have a peronal list of favourites, meaning assets that favourite or “star” so that they have them in their frontpage dashboard for quick access. An asset can be one the following

  

* Chart (that has a small title, axes titles and data)

* Insight (a small piece of text that provides some insight into a topic, e.g. "40% of millenials spend more than 3hours on social media daily")

* Audience (which is a series of characteristics, for that exercise lets focus on gender (Male, Female), birth country, age groups, hours spent daily on social media, number of purchases last month) e.g. Males from 24-35 that spent more than 3 hours on social media daily.

  

Build a web server which has some endpoint to receive a user id and return a list of all the user’s favourites. Also we want endpoints that would add an asset to favourites, remove it, or edit its description. Assets obviously can share some common attributes (like their description) but they also have completely different structure and data. It’s up to you to decide the structure and we are not looking for something overly complex here (especially for the cases of audiences). There is no need to have/deploy/create an actual database although we would like to discuss about storage options and data representations.

Note that users have no limit on how many assets they want on their favourites so your service will need to provide a reasonable response time.

A working server application with functional API is required, along with a clear readme.md. Useful and passing tests would be also be viewed favourably

It is appreciated, though not required, if a Dockerfile is included.

  

## Implementation

  

#### Features
Below you can find a comprehensive list of features implemented for this webserver.

#### Configuration
Configuration options for the webserver are passed using a YAML file file. Current available options are:

| Option | Description |
|:---|:---|
| Address: "0.0.0.0" |  String used set the interface to listen on. |
| Port: 8000 |  Port to listen on. |
| TokenSecret: "mytempsecretkey" |  String used to hash the stored used password. |
| AdminUser: "adminUser"  |  Initialized privileged user username. |
| AdminPass: "adminPass"  |  Initialized privileged user password.|
| Storage: "memory" |  Storage resource selection (currently only `memory` is implemented). |
| Cache: "memory" |  Cache resource selection (currently only `memory` is implemented). |
| Profiler: true |  Enable/Disable the attachment of a profiler to the webserver. |
| Metrics: true |  Enable/Disable the collection and exposure of metrics. |

#### API
The implemented endpoints for this API are organized in the following categories, based on their functionalities.
Additionally, the OpenAPI 3.0 specification is provided in the respective folder.
##### Miscellaneous

`GET /metrics`	*Collect the webserver's metrics (if enabled)*

`GET /debug/pprof`	*Access profiling information and record profiles/traces during operation (if enabled)*

`GET /version`	*Collect the binary's build information*

##### User authentication

`POST /signup`	*Register a new user*

`POST /signin`	*Sign in an existing user*

`GET /signout`	*Sign out of the current signed-in user*

##### Service

`GET /users/{user_id}/favourites`	*Get user's list of favourite assets*

`POST /users/{user_id}/favourites`	*A a new favourite asset to the user*

`PUT /users/{user_id}/favourites/{asset_id}`	*Update the description of a user's favourite asset*

`DELETE /users/{user_id}/favourites/{asset_id}`	*Delete a user's favourite assets*

`POST /users/{user_id}/create`	*Create a user manually (**requires privileged user**)*

`DELETE /users/{user_id}/delete`	*Delete a user manually (**requires privileged user**)*

  
#### Cache

This implementation can support any type of cache service by utilizing the created `CacheHandler` interface and implementing its methods.
The results of a *GET* request to fetch the favourite assets list of a user are cached for faster future retrieval. When this list is updated (either by adding or removing assets), the cache entry is invalidated and the next time the same *GET* request is triggered, the updated resuls will be cached anew. Cache utilizes a map using the client cookies as keys and the *GET* response as the value.
__At this time a resizing function for these maps is not implemented. To prevent a ever-growing memory allocation (since maps don't shrink) we can either create a fixed capacity map and handle the scaling (up and down when needed) manually, or create a map without specifying capacity and on events (intervals or triggers) scale the map down by recreating from scratch and adding the existing keys to it while setting the previous map to *nil* and allowing the GC to clean it up.__

Currently, `memory` is the only implemented option (in-memory cache).

  

#### Storage

This implementation can support any type of storage by utilizing the created `StorageHandler` interface and implementing its methods. Storage utilizes two maps  for fast access times, one for the users favourite assets lists and another for the users authentication information. 
__At this time a resizing function for these maps is not implemented. To prevent a ever-growing memory allocation (since maps don't shrink) we can either create a fixed capacity map and handle the scaling (up and down when needed) manually, or create a map without specifying capacity and on events (intervals or triggers) scale the map down by recreating from scratch and adding the existing keys to it while setting the previous map to *nil* and allowing the GC to clean it up.__

Currently, `memory` is the only implemented option (in-memory storage).

  

#### Metrics

If metric collection is enabled using the configuration (env) file, GO runtime metrics and metrics pertaining to the API endpoints are collected and exposed through the corresponding Prometheus scraping endpoint (*/metrics*). 

#### Tests

Execute all the tests with the following command:

```bash

go test  -v  ./...

```

Additionaly, you can extract the tests coverage profile using the commands below:

```bash

go test  -v  -coverprofile=cover.out  ./...

go tool  cover  -html  cover.out  -o  cover.html

```

Open `cover.html` with a browser.

  
  

#### Deployment

##### Native Execution

Build the webserver binary:

```bash

go build  .

```

Additionaly, build version information can be injected into the binary using the following command:

```bash

go build  -ldflags="-X 'challenge/webserver.Version=v1.0.0' -X 'challenge/webserver.BuildTime=$(date -u +"%Y-%m-%d %H:%M:%S")' -X 'challenge/webserver.CommitHash=$(git rev-parse HEAD 2>/dev/null)' -X 'challenge/webserver.BuildUser=$(id -u -n)'"  .

```

Execute the binary:

```bash

chmod +x  challenge

./challenge

```

##### Container

Build Docker image

```bash

docker build -t challenge .

```

Run a container with the build image, exposing the service port:

```bash

docker run  -d  -p  8000:8000  --name  challenge  challenge

```
