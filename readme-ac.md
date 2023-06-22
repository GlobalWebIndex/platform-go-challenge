# GlobalWebIndex - homework from Artur

## **GO** lovers 

can simply

```sh
go run .
```

to start quickly in root folder of the project.

That's all for impatient colleagues ;)

It will run, with **defaults**, the **gRPC** :9090 & **REST GW** :9080 APIs in mode **dev** and with using in-process-memory for "storage" (a transient chance to see behavior of a panic recovery middleware). 

btw Config is based on env variables. It doesn't use globals nor flags because you never know in what env it will be run. So keep it simple.

There is added automation for passing current host ip to container in dev and test modes (using hostname -I). It is required to connect simple separate containers, app/services to dbs, without additional pod or another kube/kind env. So app/services can use the same db instance for go local dev & debug (address db connection as localhost) as well as one started in container (which can connect with dbs by host ip while it's localhost is inside it's container). Automation is conditional for simplicity reason and always overlayed by dedicated env variables if passed.


## **DATA** lovers

have **more** options starting with [ArangoDB](https://www.arangodb.com/blog/) and [local AQL UI](http://http://localhost:8529/) once db is activated

## and if you are by chance a **container** lover and **makefile** lover

do both by

```sh
# with podman
# 1
make container-1-run-db-aql-once
# 2
make container-4-build-app
# 3
make container-5-run-app--rm
```

```sh
# with docker
# 1
container=docker make container-1-run-db-aql-once
# 2
container=docker make container-4-build-app
# 3
container=docker make container-5-run-app--rm
```
[Makefile](Makefile) contains more commands to explore and adapt


## **vscode** lovers

can benefit from few extensions:

- [Task Explorer](https://marketplace.visualstudio.com/items?itemName=spmeesseman.vscode-taskexplorer) - to easier use all predefined commands in the **Makefile**
- [gRPC Clicker](https://marketplace.visualstudio.com/items?itemName=Dancheg97.grpc-clicker) - to use APIs with grpc reflection support (any api enhancement will be immediately working, with no additional implementation requirement; similar for openapi)
- [Swagger Viewer](https://marketplace.visualstudio.com/items?itemName=Arjun.swagger-viewer), [Redocly OpenAPI](https://marketplace.visualstudio.com/items?itemName=Redocly.openapi-vs-code), [OpenAPI (Swagger) Editor](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi) or other - to preview or work with [OpenAPI file](proto/apidocs.swagger.json) generated from proto files in this project


## **Integration Tests** lovers

can heavy load db (ArangoDB so far) with [gofakeit](https://pkg.go.dev/github.com/brianvoe/gofakeit/v6) data by running [Example_gRPC_Client_loading_fake_data()](test/fake/fake_test.go) (*it is an ad hoc study in dev; code will be splitted and spreaded asap*). Any big load(s) will be also adored by data lovers exploring amazing doc and graph functionalities in [local AQL UI](http://http://localhost:8529/)

```sh
# A - start db
# make container-1-run-db-aql-once
# or start if was stopped
make container-3-start-db-aql

# B - start app/server
# go run .
# or run in container
make container-5-run-app--rm

# C - load gofakeit data to db
# requires go env for now; no container yet
make container-x-go-test-Example_gRPC_Client_loading_fake_data
```


## Please **enjoy :)**

as I do, in time of adding new functionalities to this inspiring task

More to come soon ...
