### How to Build and Run the Container

Once you have this `Dockerfile` in the root of your project, you can build and run your container with these simple commands:

1.  **Build the Docker Image**: Run this command from your project's root directory. The `-t` flag gives your image a name and tag (e.g., `gwi-favourites:latest`).

    ```sh
    docker build -t gwi-favourites .
    ```

2.  **Run the Container**: This command starts a new container from your image. The `-p` flag maps the container's port 8080 to your host machine's port 8080, allowing you to access the API.

    ```sh
    docker run -p 8080:8080 gwi-favourites