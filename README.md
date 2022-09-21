# GWI Exercise Server

## Table of Contents
- [Introduction](#introduction)
- [List of services](#list-of-services)
- [List of ports used](#list-of-ports-used)
- [How to deploy everything](#how-to-deploy-everything)
- [How to test everything](#how-to-test-everything)
- [How to close everything](#how-to-close-everything)
- [How to run indivindually](#how-to-run-indivindually)
- [Bugs](#bugs)


## Introduction
This is a simple server for admins to add/update/delete assets (audience, chart and insights) and users to favour the assets.

## List of services
- MariaDB for storing user and CV information.
- Adminer for looking at the data in MariaDB.
- GWI Server for the REST-API.

## List of ports used

The system will take these ports: 3306, 8080, 8000 </br>
Before starting the services, <span style="color:red">make sure you don't have any application that use any of these ports.</span>

## How to deploy everything
First, make sure you have docker-compose installed. </br>
After installing docker-compose, you have to build the docker images.
```shell
docker-compose -f docker-compose.yml -f docker-compose.server.yml build
```
Start all the services
```shell
docker-compose -f docker-compose.yml -f docker-compose.server.yml up -d
```

## How to test everything
There are two options to play with the system.</br>

The first option is to try it out from [Swagger](http://localhost:8000/swagger/index.html), but on authorization("Authorize" button) of the token don't forget to add "Bearer" front of the token.</br>

The second option is to try it out from Postman by importing the file 'GWI-Exercise-Go.postman_collection.json'</br>
Both contain some information on how to create a user, login and CRUD assets </br>

After you create a user and an asset, you can check the SQL data from [Adminer](http://localhost:8080/)</br>
For Adminer, the username and password is 'user'</br>

## How to close everything
Because the system is based on docker-compose, just call this command line to close everything down.
```shell
docker-compose -f docker-compose.yml -f docker-compose.server.yml down
```
The system uses only ephemeral volumes, so you don't have to worry on looking for any volume to delete.


## How to run indivindually
Run docker-compose to start the DB
```shell
docker-compose up
```

Compile and run the binary
```shell
go build

./platform-go-challenge
```
OR

Run the docker image
```shell
docker build -t gwiserver .

docker run --env-file .env.docker -p 8000:8000 --network platform-go-challenge_default -it gwiserver
```

## Bugs
Known bugs
- The isDesc query option for getting any list, does not seem to work properly