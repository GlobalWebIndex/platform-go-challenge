# GlobalWebIndex Engineering Challenge

Resolution of the proposed technical test.

Created using go version `1.17`.

## How to run it?

Clone the project, get into the correct folder and execute the command:

```bash
go run cmd/main.go
```

If you want to use the component as a docker service, run:

```bash
# build the image
docker build -t gwi_test .

# run it
docker run -it -p 4567:4567 gwi_test
```

To run a script that seeds some basic data, run:
```bash
chmod +x populate.sh

./populate.sh
```


## How to use it?

### Add favourite

To add a favourite, send a POST request using the URL `http://localhost:4567/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/add` and a valid payload:

A valid uuid v4 value must be used to provide a correct userid path.

```json
{
  "type": "audience",
  "audience": {
    "id": "a85b8a66-15f4-49be-a0a0-90df6ce2b79f",
    "gender": "updated",
    "born_country": "foo",
    "age_group": "bar",
    "daily_hours_social_media": "foo",
    "purchases_last_month": "bar"
  }
}
```

```bash
# Example of curl request to add a new audience favourite to user `a04c6d0e-0115-4931-9911-35d5fe61983e`

curl -s  --request POST  'http://localhost:4567/user/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/add' \
-d '{
  "type": "audience",
  "audience": {
    "id": "a85b8a66-15f4-49be-a0a0-90df6ce2b79f",
    "gender": "updated",
    "born_country": "foo",
    "age_group": "bar",
    "daily_hours_social_media": "foo",
    "purchases_last_month": "bar"
  }
}'
```

You will get a response with this format if the process succeed:

```json
{
  "status":200,
  "success":true
}
```

### Edit a favourite

To update a favourite, send a PUT request using the URL `http://localhost:4567/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/edit` and a valid payload:

A valid uuid v4 value must be used to provide a correct userid path.

```json
{
  "type": "audience",
  "audience": {
    "id": "a85b8a66-15f4-49be-a0a0-90df6ce2b79f",
    "gender": "updated",
    "born_country": "foo_updated",
    "age_group": "bar_updated",
    "daily_hours_social_media": "foo_updated",
    "purchases_last_month": "bar_updated"
  }
}
```

```bash
# Example of curl request to edit an existing audience favourite for user `a04c6d0e-0115-4931-9911-35d5fe61983e`

curl -s  --request PUT 'http://localhost:4567/user/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/edit' \
-d '{
  "type": "audience",
  "audience": {
    "id": "a85b8a66-15f4-49be-a0a0-90df6ce2b79f",
    "gender": "updated",
    "born_country": "foo_updated",
    "age_group": "bar_updated",
    "daily_hours_social_media": "foo_updated",
    "purchases_last_month": "bar_updated"
  }
}'
```

You will get a response with this format if the process succeed:

```json
{
  "status":200,
  "success":true
}
```

### Delete a favourite

To delete a favourite, send a DELETE request using the URL `http://localhost:4567/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/delete` and a valid payload:

A valid uuid v4 value must be used to provide a correct userid path.

```json
{
  "type": "audience",
  "audience": {
    "id": "a85b8a66-15f4-49be-a0a0-90df6ce2b79f"
  }
}
```

```bash
# Example of curl request to delete an existing audience favourite for user `a04c6d0e-0115-4931-9911-35d5fe61983e`

curl -s  --request DELETE 'http://localhost:4567/user/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/delete' \
-d '{
  "type": "audience",
  "audience": {
    "id": "a85b8a66-15f4-49be-a0a0-90df6ce2b79f"
  }
}'
```

You will get a response with this format if the process succeed:

```json
{
  "status":200,
  "success":true
}
```

### List user favourites

To list user's favourites, send a GET request using the URL `http://localhost:4567/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/list`:

A valid uuid v4 value must be used to provide a correct userid path.

```bash
# Example of curl request to retrieve the favourites list of user `a04c6d0e-0115-4931-9911-35d5fe61983e`

curl -s  --request GET 'http://localhost:4567/user/a04c6d0e-0115-4931-9911-35d5fe61983e/fav/list'
```

You will get a response with this format if the process succeed:

```json
{
  "status": 200,
  "success": true,
  "length": 3,
  "data": {
    "audience": [
      {
        "id": "a85b8a66-15f4-49be-a0a0-90df6ce2b79f",
        "gender": "foobar",
        "born_country": "foo",
        "age_group": "bar",
        "daily_hours_social_media": "foo",
        "purchases_last_month": "bar"
      }
    ],
    "chart": [
      {
        "id": "4822a1c2-7c9c-4df6-8581-449b5d5a5762",
        "title": "foobar",
        "axis_y_title": "foo",
        "axis_x_title": "bar",
        "data": null
      }
    ],
    "insight": [
      {
        "id": "7e923b61-2828-42f6-bbe7-55461795988a",
        "description": "foobar"
      }
    ]
  }
}

```