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

# My solution
## How to use

Build the Dockerfile and run a container:
```
docker build --tag favourites-service .
docker run -p 8086:8086 favourites-service
```

`GET` the favourites for a user. By default this is an empty set. Note that this example request includes a jwt token valid until 19/11/2024:
```
$ curl -f 'http://localhost:8086/favourites?offset=0&pagesize=10' -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJmYXZvdXJpdGVzIiwidXNlcmlkIjoiNjA5ZGFjOWMtYWM3OS00ZGM4LWExZjUtZjJhZjdhNTUxOWNmIiwiaWF0IjoxNzI4MDM2ODYyLCJleHAiOjE3MzIwNTA0NjJ9.19NYScvE4f6FIHAojMcn0sv-xOgzrCE26jFsTqddoOM"
[]'
```

`POST` some favourites for this user. To post a favourite, include a uuid (random), description (whatever you like) and resource type (audience|chart|insight):
```
$  curl -f 'http://localhost:8086/favourites' -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJmYXZvdXJpdGVzIiwidXNlcmlkIjoiNjA5ZGFjOWMtYWM3OS00ZGM4LWExZjUtZjJhZjdhNTUxOWNNmIiwiaWF0IjoxNzI4MDM2ODYyLCJleHAiOjE3MzIwNTA0NjJ9.19NYScvE4f6FIHAojMcn0sv-xOgzrCE26jFsTqddoOM" -X POST -H 'Content-Type: application/json' -d '{"id": "5791fb0a-130d-4150-835a-d0433bf6eb73", "description": "my favourite chart", "resourceType": "chart"}'
{"Description":"my favourite chart","ResourceType":"chart","Id":"5791fb0a-130d-4150-835a-d0433bf6eb73"}'

$ curl -f  'http://localhost:8086/favourites' -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJmYXZvdXJpdGVzIiwidXNlcmlkIjoiNjA5ZGFjOWMtYWM3OS00ZGM4LWExZjUtZjJhZjdhNTUxOWNmIiwiaWF0IjoxNzI4MDM2ODYyLCJleHAiOjE3MzIwNTA0NjJ9.19NYScvE4f6FIHAojMcn0sv-xOgzrCE26jFsTqddoOM" -X POST -H 'Content-Type: application/json' -d '{"id": "19054c1a-5a10-43e1-aef1-14a37fddb351", "description": "my favourite insight", "resourceType": "insight"}'
{"Description":"my favourite insight","ResourceType":"insight","Id":"19054c1a-5a10-43e1-aef1-14a37fddb351"}'
```

`GET` the user's favourites (chart and insight data are generated randomly, but IDs are the ones supplied before):
```
$ curl -f 'http://localhost:8086/favourites?offset=0&pagesize=10' -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJmYXZvdXJpdGVzIiwidXNlcmlkIjoiNjA5ZGFjOWMtYWM3OS00ZGM4LWExZjUtZjJhZjdhNTUxOWNmIiwiaWF0IjoxNzI4MDM2ODYyLCJleHAiOjE3MzIwNTA0NjJ9.19NYScvE4f6FIHAojMcn0sv-xOgzrCE26jFsTqddoOM"
[
    {
        "Description": "my favourite chart",
        "Chart": {
            "Id": "5791fb0a-130d-4150-835a-d0433bf6eb73",
            "Title": "sales chart number 5",
            "XAxisTitle": "time",
            "YAxisTitle": "number of sales",
            "DataPoints": [
                {
                    "X": 0,
                    "Y": 0
                },
                {
                    "X": 1,
                    "Y": 2
                },
                {
                    "X": 2,
                    "Y": 4
                }
            ]
        }
    },
    {
        "Description": "my favourite insight",
        "Insight": {
            "Id": "19054c1a-5a10-43e1-aef1-14a37fddb351",
            "Description": "70% of people between 38 and 49 spend more than 5 hours per day online"
        }
    }
]'

For full API documentation, visit `http://localhost:8086/swagger/index.html` when the service is running.

```

## About the implementation
I chose to implement the web server without any data persistence. 
I added some tests against a few endpoints in ![main_test.go](cmd/main_test.go) and some unit tests for one of the repositories in ![repository_test.go](internal/repository/favourite/repository_test.go).

The `/favourites` endpoint supports paging to improve response time.

I imagine the data for audiences, charts and insights would come from unrelated datasources, for example audiences might come from some cache, while charts can be generated by an external service relatively slowly. When getting a user's favourites, I retrieve the data from these 3 data sources concurrently and stitch it together into a heterogenous list.

I added very rudimentary authentication in the form of a jwt token that must be present on the request. 
The signing key used is `[]byte("12345678123456781234567812345678")`.
The token should include a custom claim e.g. `"userid": "609dac9c-ac79-4dc8-a1f5-f2af7a5519cf"`.
This is used to identify the current user.

The token should be included as a bearer token in an authorization header, which is supplied in the `curl` examples above.

A production-ready version would include more tests, and timeouts when fetching data from external services, potentially also graceful degradation.