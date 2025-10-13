# Postman Instructions

## Import the collection

1. Open Postman.

2. Click Import: Raw text / File.

3. Copy-paste or upload the provided postman_collection.json file.

4. You’ll see a collection named “GWI API”.

## Get a Token

1. Run the Login (Get Token) request:

   - URL: POST http://localhost:8080/api/login

   - Body:

    ```
    {
        "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsImV4cCI6MTc2MDQ3MTMzNywibmJmIjoxNzYwMzg0OTM3LCJpYXQiOjE3NjAzODQ5Mzd9.WGDqmRjsv2GSZhHa_KwwD0u0tq6V2UUyH3q_dgVD7xw",
        "user_id": "user123"    
    }
    ```


2. Copy the returned token value
   
## Set the Token Variable

1. click the “GWI API” collection name.

2. Go to the Variables tab.

3. Set the {{token}} variable to the token value copied in step 2.
   

Now, all authenticated requests will automatically use your token.

## Run the Collection

- Feel free to modify the collection and test the requests.