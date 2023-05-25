# GlobalWebIndex Engineering Challenge

## Requirements
The application requires docker and docker-compose to run.

You also need the `.env` file provided below:
```dotenv
DATASOURCE_URI=mongodb://mongodb:27017
DATASOURCE_TIMEOUT_IN_MILLISECONDS=5000
DATASOURCE_DATABASE=gwi
DATASOURCE_COLLECTION=user
```
## Commands
Use `make` to run the following available commands.

| **Command**      | **Usage**                   |
|------------------|-----------------------------|
| `start`          | Starts the application      |
| `stop`           | Stops the application       |
| `migration-up`   | Applies all up migrations   |
| `migration-down` | Applies all down migrations |

## Known Issues

Due to lack of time, the following issues are present:
* the application does not implement an inbound HTTP adapter responsible for receiving and handling HTTP calls.
* the `buildPatchData` method of the mongodb adapter has not been implemented properly.
* the application has no unit, integration or fuzz tests
* the application does not implement any authorization or authentication schemes
* the `Execute` function responsible for running the application should initialize and inject all the necessary services and adapters. Instead, it fetches the only document on the database and prints it to stdout.