# Gwi API

## Table Of Contents
- [Description](#description)
- [Mongo](#mongo)
    * [Schema](#schema)
    * [Mongo Image](#mongo-image)
    * [Mongo Initialization](#mongo-initialization)
- [Swagger](#swagger)
- [API](#api)
    * [Image](#image)
    * [Configuration](#configuration)
    * [Deployment](#deployment)
    * [Tests](#tests)

### Description 
The backend of the `GWI` platform, where `Dashboard` interacts on the following entities:

1. assets
2. favourites
3. users

### Mongo


### Schema 
The schema of the `Mongo` database can be found at [MongoSchema](docs/GwiDatabaseMongo.md)

#### Mongo Image
Mongo image used is declared in the [compose.yml](deployments/compose.yml#L26)

#### Mongo Initialization
The necessary collections as documented in the [MongoSchema](docs/GwiDatabaseMongo.md) are initialized via the [init script](deployments/mounts/mongo/init/0-init.sh)</br>
The `init script` is passed as a volume in the associated service via the [compose.yml](deployments/compose.yml#L37)

### Swagger
The [swagger](./deployments/mounts/swagger/swagger.yaml) is used to document the functionality of the endpoints as well as the information schema inputs and outputs.

`Swagger UI` is deployed as a dedicated `container service`. </br>
The service can be found in the [compose](deployments/compose.yml#L34). </br>
The aforementioned service depends on the [init service](deployments/compose.yml#L12), responsible for generating the `swagger.yaml`, containing the appropriate url as declared in the `.env` file via `PORT_API` environmental variable.</br>

Then `init` service uses the [init.sh](scripts/init.sh) to create the needed files for the swagger-UI and logging as it will be explained in a later section.

### API

#### Image
Image is built based on the [Dockerfile](./Dockerfile).

#### Configuration
Configuration needed is given via the following template and is passed via a config file which can be found [here](./deployments/mounts/config/config.yml)

The aforementioned file is responsible for configuring the `logger`, the `mongo database` and the `API` TLS.

#### Deployment

- Local

    From the `/deployments` directory the following command needs to be executed

    ```bash
    docker compose up --build -d
    ```

    The aforementioned command will also build the image based on the [Dockerfile](./Dockerfile) as documented in the [compose](./deployments/compose.yml#L42) file

    The [compose](./deployments/compose.yml) file deploys the following container services 

        1. api
        2. swaggerAPI
        3. mongoDB

#### Tests

Tests can be found under the `./tests` directory.</br>
This repo supports support testing for the first two endpoints.

The following command needs to be executed in the aforementioned directory:

```
go test -v
```










































