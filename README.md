# Platform Go Challenge

---

## Description

Service that provides a REST API offering CRUD operations for adding, removing, editing
and getting a User's favourite Assets (Chart, Insight, Audience)

---

## Run

`cd script && make start-app`

* This command will start the app with `localhost` address and `:8080` port (specified in build/Dockerfile.dev and .env)

`cd script && make migrate-data-small`

* This command will add to DB 10 Users and 3 Assets per User for testing reasons

`cd script && make migrate-data-large`

* This command will add to DB 10 Users and 10.000 Assets per User for testing reasons

Then you can add, edit, remove and get User's favourites like the examples in /examples  
directory. To generate the needed Bearer token, please call /token endpoint with username  
& password like in the example

---

## Makefile Commands

| Command                         | Usage                                                                  |
|---------------------------------|------------------------------------------------------------------------|
| start-app                       | `Start app`                                                            |
| kill-app                        | `Stop app`                                                             |
| rebuild-app                     | `Rebuild app`                                                          |
| tests-all                       | `Run both unit and integration tests`                                  |
| tests-benchmark                 | `Run benchmark tests`                                                  |
| tests-unit                      | `Run unit tests `                                                      |
| tests-file FILE={filePath}      | `Run specific file test`                                               |
| generate-mock FILE={filePath}   | `Generate mock for a specific file`                                    |
| run-linter                      | `Runs linter`                                                          |
| generate-swagger-files          | `Generates swagger.json definitions in Docs dir`                       |
| tests-package PACKAGE={package} | `Run specific package test`                                            |
| tests-all-with-coverage         | `Run both unit and integration tests via docker with coverage details` |
| migrate-data-small              | `Run migration of a small dataset Users & Assets in the DB`            |
| migrate-data-large              | `Run migration of a larger dataset Users & Assets in the DB`           |

* All these are executed through docker containers
* In order to execute makefile commands type **make** plus a command from the table above

  make {command}

---

## Notes

1. .env is pushed to Git only for the Assessment purpose. Config and .env files should never be tracked.
2. There are three Dockerfile files.
    1. Dockerfile is the normal, production one
    2. Dockerfile.dev is for setting up a remote debugger Delve
    3. Dockerfile.utilities is for building a docker for "utilities" like running tests,  
       linting etc
3. Asset types' fields and their validation was done very simply for the scope of this assessment
4. As described at the assignment, the User sees a list of Assets and adds some to their
   favorites (by starring them for example). So the Assets are existing. This is why
   the AddAsset() takes just IDs and just associates a User with an Asset, and not create it.

## Known Issues

1. Unit tests for adding, removing, editing assets are done only for Charts and only for happy paths.
2. JWT mechanism just requires a fake username and password to generate a JWT token and does NOT do
   actual login due to lack of time. Also no test created for it

## Performance

If Users favourites data were causing performance issues, I would investigate the following solutions:

1. Add a cache to the repository layer, over the DB and keep there the User's favourite Assets
   in order to have faster access to it
1. Limit Payload by implementing Pagination
1. Limit Payload by compressing the data with Accept-Encoding: gzip or compress
1. Add indexes in DB. For example we could add indexed to the Users-Assets tables e.g "users_charts" on user_id
   and chart_id

* All these solutions would need to be discussed with the Engineering team to find the most suitable as
  they all have their pros and cons

## Security

1. JWT mechanism added for Authentication and Authorization (incomplete - see Known Issues)