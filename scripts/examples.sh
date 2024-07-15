# GET existing user favorites
curl -X GET http://localhost:8080/users/1/favorites

# GET user favorites which do not exist
curl -X GET http://localhost:8080/users/999999/favorites

# DELETE existing user favorite
curl -X DELETE http://localhost:8080/users/1/favorites/1

# DELETE user favorite which does not exist
curl -X DELETE http://localhost:8080/users/1/favorites/999999

# ADD valid user favorite
curl -X POST http://localhost:8080/users/1/favorites \
     -H "Content-Type: application/json" \
     -d '{
          "id": 100,
          "type": "Audience",
          "description": "This audience is a 40 year old",
          "age": 40,
          "ageGroup": "25-45",
          "gender": "Male",
          "birthCountry": "USA",
          "hoursSpentOnMedia": 4,
          "numberOfPurchases": 10
         }'

# ADD user favorite which already exists
curl -X POST http://localhost:8080/users/1/favorites \
     -H "Content-Type: application/json" \
     -d '{
          "id": 2,
          "type": "Insight",
          "description": "Sample Insight for testing",
          "text": "Testing Insight"
         }'

# ADD user favorite with invalid type
curl -X POST http://localhost:8080/users/1/favorites \
     -H "Content-Type: application/json" \
     -d '{
          "id": 200,
          "type": "INVALIDTYPE",
          "description": "Sample Insight for testing",
          "text": "Testing Insight"
         }'

# EDIT the previously added user favorite
curl -X PUT http://localhost:8080/users/1/favorites/100 \
     -H "Content-Type: application/json" \
     -d '{
          "id": 100,
          "type": "Audience",
          "description": "Updated Audience",
          "age": 18,
          "ageGroup": "18-25",
          "gender": "Female",
          "birthCountry": "Greece",
          "hoursSpentOnMedia": 15,
          "numberOfPurchases": 25
         }'

# EDIT user favorite with mismatched id (assetID in URL and id in body)
curl -X PUT http://localhost:8080/users/2/favorites/2 \
     -H "Content-Type: application/json" \
     -d '{
          "id": 1,
          "type": "Insight",
          "description": "Sample Insight for testing",
          "text": "Testing Insight"
         }'

# EDIT user favorite with mismatched type
curl -X PUT http://localhost:8080/users/1/favorites/100 \
     -H "Content-Type: application/json" \
     -d '{
          "id": 100,
          "type": "Chart",
          "description": "This is a chart",
          "title": "Chart",
          "xAxesTitle": "X-Axis",
          "yAxesTitle": "Y-Axis",
          "dataPoints": [
          {
               "X": 10,
               "Y": 10
          },
          {
               "X": 20,
               "Y": 20
          }
          ]
        }'