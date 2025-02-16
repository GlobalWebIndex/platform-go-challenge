GWI Favorites API – Engineering Challenge Submission

This project is a solution for the GWI Engineering Challenge, implementing a secure and performant web API for managing user favorite assets, built using Go.

Features Implemented

Full CRUD Operations: Comprehensive support for adding, retrieving, editing, and deleting user favorite assets.
Concurrency Handling: Parallel retrieval of assets to ensure fast responses, utilizing sync.WaitGroup and sync.Mutex.
Secure JWT Authentication: Secure token-based authentication with middleware for protected endpoints. Ensures that only authorized users can access the API.
Token Management: Implemented a mechanism to store and invalidate tokens, ensuring that each user has only one valid token at a time.
HTTPS Support: Configured HTTPS at the code level using self-signed certificates (for development purposes only), ensuring secure communication. 
Robust Error Handling and Input Validation: Detailed error responses with checks for user existence, asset duplication, and proper input format.
Dockerization: Provided a Dockerfile for easy containerized deployment of the API. Dockerized deployment also supports HTTPS, with clear instructions provided for generating certificates and running the app securely.
Commented Proxy Code: Added commented-out code for proxy trust settings to defend against IP spoofing if needed.

API Endpoints

POST /login – Generates a JWT token for authentication.
GET /favorites/:user_id – Retrieves all favorite assets for a user concurrently.
POST /favorites – Adds a new favorite asset for a user with type-based validation.
DELETE /favorites/:user_id/:asset_id – Removes a favorite asset, checking for user and asset existence.
PUT /favorites/:user_id/:asset_id – Edits the description of a favorite asset, with checks for existence and valid input.

Technologies Used

Go 1.24 – For building the web API.
Gin Framework – Lightweight HTTP web framework.
JWT (github.com/dgrijalva/jwt-go) – For authentication.
Docker – For containerization.
Postman – For API testing.

Running the Application

HTTPS Support
This application supports HTTPS both when running locally and inside Docker. To enable HTTPS, you need to generate self-signed certificates:

Run the following command in your terminal or PowerShell:

openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

req: Request a new certificate.
-x509: Create a self-signed certificate.
-newkey rsa:4096: Generate a new RSA key (4096 bits).
-keyout key.pem: Save the private key as key.pem.
-out cert.pem: Save the certificate as cert.pem.
-days 365: Certificate valid for 1 year.
-nodes: No password for the private key.
Follow the prompts, fill in the required details, and the process will generate:

cert.pem (certificate)
key.pem (private key)
Place these two files in the project root directory (same level as main.go).

Without Docker:
go run main.go
Access the API at https://localhost:8080.

With Docker:
docker build -t gwi-favorites-api .
docker run -d -p 8080:8080 --name gwi-container gwi-favorites-api
Note: The Docker build uses the same cert.pem and key.pem files from the project root directory to enable HTTPS.

Testing

Run the provided unit tests with:
go test ./tests -v

API Testing Instructions and Payloads
You can use tools like Postman to test the API. Below are the instructions and example payloads for each endpoint:

1. Login Endpoint
   URL: POST https://localhost:8080/login

No additional headers needed when using form-data.
Body (form-data): user_id=123
Response:
{
"token": "your_jwt_token_here"
}

Use the returned token for authentication in all subsequent requests by adding it as a Bearer token in the Authorization header.

2. Add Favorite Endpoint
   URL: POST https://localhost:8080/favorites
   Headers:

Content-Type: application/json
Authorization: Bearer your_jwt_token_here

Payloads:
Chart Example:
{
"user_id": "123",
"asset": {
"id": "201",
"type": "chart",
"description": "Tech Stocks",
"title": "Tech Market Overview",
"axes_titles": ["Time", "Price"],
"data": [3500, 3600]
}
}

Insight Example:
{
"user_id": "123",
"asset": {
"id": "202",
"type": "insight",
"description": "50% of users spend 5+ hours online daily"
}
}

Audience Example:
{
"user_id": "123",
"asset": {
"id": "203",
"type": "audience",
"description": "Target Audience",
"gender": "Male",
"birth_country": "Greece",
"age_group": "24-35",
"hours_spent_daily_on_social_media": 4,
"purchases_last_month": 5
}
}

3. Get Favorites Endpoint
   URL: GET https://localhost:8080/favorites/123
   Headers:

Authorization: Bearer your_jwt_token_here

Response Example:

{
"favorites": [
{
"id": "201",
"type": "chart",
"description": "Tech Stocks",
"title": "Tech Market Overview",
"axes_titles": ["Time", "Price"],
"data": [3500, 3600]
},
{
"id": "202",
"type": "insight",
"description": "50% of users spend 5+ hours online daily"
}
]
}

4. Edit Favorite Endpoint
   URL: PUT https://localhost:8080/favorites/123/201
   Headers:

Content-Type: application/json
Authorization: Bearer your_jwt_token_here

Payload:
{
"new_description": "Updated Tech Stocks Overview"
}

Response:
{
"message": "Asset updated"
}

5. Remove Favorite Endpoint
   URL: DELETE https://localhost:8080/favorites/123/201
   Headers:

Authorization: Bearer your_jwt_token_here

Response:
{
"message": "Asset removed"
}


Notes

HTTPS is enabled for local development using self-signed certificates.
JWT tokens are managed to prevent multiple active tokens per user.
JWT authentication middleware is implemented, handling token validation for all protected routes.
Concurrency is handled in asset retrieval to ensure optimal performance.
Dockerization is provided for easy deployment.
Detailed error handling ensures clarity in API responses.
No Database Used. An in-memory map is used for asset storage for simplicity and demonstration purposes.

Potential Storage Options (SQL Focused)
Given the requirement for user-specific favorites, a relational database such as PostgreSQL, MySQL, or SQLite would be an ideal choice due to its robustness, scalability, and widespread use in production environments.

Advantages of SQL Databases for This Project:

Relational Integrity: Ensures relationships between users and their favorite assets are well-maintained.
Indexing Support: Indexes can be created on frequently queried fields (such as user_id and asset_type) for faster lookups.
Transactions: Supports ACID properties, ensuring data integrity during concurrent operations.
Scalability: Can handle growing datasets efficiently through proper indexing and sharding strategies.
Familiarity: SQL is widely used and supported, making it an optimal choice for maintainability and long-term scaling.

Note: This project was my first experience writing Go. It was my honor to embrace the challenge and familiarize myself with the language.

Final Thoughts

Every design choice was made with care, ensuring a clean, efficient, and maintainable codebase.
I trust that this submission demonstrates both my technical proficiency and commitment to delivering high-quality solutions.
Thank you for the opportunity to showcase my skills — I look forward to the possibility of contributing to the GWI team.






