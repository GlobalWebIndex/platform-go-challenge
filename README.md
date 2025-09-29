# **Favourite Assets** 

## **Overview**

**This project is a Go-based server that manages Users, Assets, and Favourites, with role-based access control handled via Keycloak.**

   ⦁	Users: Can be created, updated, listed, and deleted (admin-only for some operations).
   
   ⦁	Assets: Supports different types like Charts, Insights, and Audience data. Admins can manage all assets.
   
   ⦁	Favourites: Users can mark assets as favourites. 
   
   ⦁   The server stores all of them in a local in-memory database with sharded storage for efficiency.

**Keycloak Integration**

   ⦁	Keycloak is used for authentication and authorization.*
   
   ⦁	The server reads the JWT token, checks the user role, and determines access.
   
   ⦁	In an actual project, you could connect your DB to Keycloak to validate the user ID and fetch user data directly.
   
   ⦁	For this demo, Keycloak runs in Docker, and the server verifies tokens locally.

  *In this demo, the server parses JWT tokens from Keycloak to read roles and user info. In a production environment, the server should validate the token signature and issuer before trusting the   claims.


## **How to run**

⦁	**Use Docker Compose to start both the Go server and Keycloak: docker-compose up --build**

   This command will:

   Build your Go server image

   Start your server on http://localhost:8080

   Start Keycloak on http://localhost:8081

   The Keycloak realm favourite-assets is automatically imported from keycloak-realm/realm-export.json.


⦁	**Get a JWT Token from Keycloak**

  Make a call to retrieve a token 
   
  **Admin User** 

     curl -X POST "http://localhost:8081/realms/favourite-assets/protocol/openid-connect/token" \
     -d "grant_type=password&client_id=favourite-assets&username=admin&password=admin&client_secret=secret"
    
  **Plain user**
       
       curl -X POST "http://localhost:8081/realms/favourite-assets/protocol/openid-connect/token" \
        -d "grant_type=password&client_id=favourite-assets&username=user&password=user&client_secret=secret"

  Include the token in your requests as a Bearer Token

⦁ **Call endpoints**

 **Users** 
      
       Create User (Admin only)
       POST http://localhost:8080/users/
         {
          "name": "John Doe",
          "email": "john@example.com"
         }

        List Users (Admin only)
        GET http://localhost:8080/users/
        
        Get User by ID (All roles)
        GET http://localhost:8080/users/by-id?userId=<uuid>
          
        Update User (All roles)
        PUT http://localhost:8080/users/?userId=<uuid>
            {
             "name": "John Smith",
             "email": "johnsmith@example.com"
            }
        
        Delete User (Admin only)
        DELETE http://localhost:8080/users/?userId=<uuid>
        
  **Assets**
  
        Create Asset (Admin only)
        POST http://localhost:8080/assets/
        Example for Chart asset:
        
        {
          "type": "chart",
          "title": "Sales Chart",
          "description": "Monthly sales data",
          "xAxis": "Month",
          "yAxis": "Revenue"
        }
        Example for Insight asset
        {
          "type": "insight",
          "description": "Customer behavior",
          "text": "Most customers buy on weekends."
        }
        
        Example for Audience asset:
        
        {
          "type": "audience",
          "description": "Target demographic",
          "gender": "male",
          "birthCountry": "USA",
          "ageGroup": "25-34",
          "hoursSocialDaily": 3,
          "purchasesLastMonth": 5
        }
        
        List Assets (All roles)
        GET http://localhost:8080/assets/
        Optional filter: GET http://localhost:8080/assets/?type=chart
        
        Get Asset by ID (All roles)
        GET http://localhost:8080/assets/by-id?assetId=<uuid>
        
        Update Asset (Admin only)
        PUT http://localhost:8080/assets/?assetId=<uuid>
        Body is similar to create request, with updated values.
        
        Delete Asset (Admin only)
        DELETE http://localhost:8080/assets/?assetId=<uuid>
        
  **Favourites**
        
        Add Favourite (All roles)
        POST http://localhost:8080/favorites/?userId=<uuid>&assetId=<uuid>
        
        Remove Favourite (All roles)
        DELETE http://localhost:8080/favorites/?favouriteId=<uuid>
        
        List Favourites by User (All roles)
        GET http://localhost:8080/favorites/?userId=<uuid>
        
        Get Favourite by ID (All roles)
        GET http://localhost:8080/favorites/by-id?favouriteId=<uuid>

## **DB Schema**
    
    +----------------+
    |     users      |
    +----------------+
    | id (PK)        |
    | name           |
    | email          |
    | created_at     |
    | updated_at     |
    +----------------+
             |
             | 1 : M   (one user can have many favourites)
             v
    +----------------+
    |   favourites   |
    +----------------+
    | id (PK)        |
    | user_id (FK)   | ---> users.id
    | asset_id (FK)  | ---> assets.id
    | asset_type     |
    | created_at     |
    +----------------+
             |
             | M : 1   (many favourites can point to one asset)
             v
    +----------------+
    |     assets     |
    +----------------+
    | id (PK)        |
    | type           | <-- chart / insight / audience
    | description    |
    | created_at     |
    | updated_at     |
    +----------------+
            / | \
           /  |  \
          v   v   v
    +---------+ +---------+ +----------------+
    |  chart  | | insight | | audience       |
    +---------+ +---------+ +----------------+
    | title   | | text    | | gender         |
    | x_axis  | |         | | birth_country  |
    | y_axis  | |         | | age_group      |
    |         | |         | | hours_social   |
    |         | |         | | purchases_last |
    +---------+ +---------+ +----------------+
    

