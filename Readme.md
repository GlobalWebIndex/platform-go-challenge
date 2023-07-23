# GWI Testing API

This API was designed for GWI project(testing)

## Running the API

To get the API up and running, follow these steps:

1. **Generate protobuf files**
   Start by generating the protobuf files. Use the following command:

   ```shell
   make generate
   ```

   This will generate the protobuf files and place them in the `./pkg` directory.

2. **Configure JWT**
   Open the `config.env` file located in the `./config/config.env` directory. You will need to add a key-value pair for `jwt.signedKey`. This key can be obtained from OAuth or Firebase, or generated directly using OpenSSL with the following command:

   ```shell
   openssl rand -base64 64
   ```

3. **API Documentation**
   You can explore all available API endpoints using the Swagger UI available at:
   ```
   http://localhost:8901/api/swagger-ui/
   ```
   Please note that the Swagger file is currently not complete, and therefore you should only use it to understand the API request style. After authorization, remember to add the `Authorization` header field with the `Bearer {your token}` to your requests. The token can be obtained during the sign-in or sign-up process.

## Disclaimer

Please note that this API was initially designed for testing purposes and its structure might undergo changes based on more detailed requirements.

For testing purposes, the current version of the API simulates a database in memory instead of using a real database setup. Depending on the finalized requirements, a decision can be made between using a relational or non-relational database in the future.
