# GlobalWebIndex Engineering Challenge

## Introduction

This solution provides user endpoints with basic authedication functionality along with a CRUD API for user assets. For this challenge a postgres database with sample data has been created in https://api.elephantsql.com/.

## Run solution

A web server will start at http://localhost:8000 by running `go run main.go` 

### Sample Data
In the created postgres database there are some sample data (charts, insights, audiences) for user with id 1.

## API

## Sign up
`POST /api/users/signup`
### Request
    curl --location --request POST 'http://localhost:8000/api/users/signup'
    --header 'Content-Type: application/json' 
    --data-raw '
    {
        "Email": "pdimiropoulou@mail.com",
        "Password": "1234"
    }'
  
### Response
Returns an integer for the new user id 

## Login
`POST /api/users/login`
### Request
    curl --location --request POST 'http://localhost:8000/api/users/signup'
    --header 'Content-Type: application/json' 
    --data-raw '
    {
        "Email": "pdimiropoulou@mail.com",
        "Password": "1234"
    }'
  
### Response
The token should be used in the header of user assets requests having as key:Authorization and value the value of the token provided in the response below.

    {
        "id": 1,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ"AA
    }
 
## Get user assets
`GET /api/assets/{user_id}`
### Request
    curl --location --request GET 'http://localhost:8000/api/assets/1?limit=10&offset=0'
    --header 'Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ' 
 
Parameters limit and offset can be used for pagination, if these parameters are not provided the response will contain all user assets

### Response
The token should be used in the header of user assets requests having as key:Authorization and value the value of the token provided in the response above.

    {
    "user_id": 1,
    "charts": [
        {
            "ID": 1,
            "UserId": 1,
            "Favourite": true,
            "Title": "Calls Chart",
            "XAxes": "Calls",
            "YAxes": "Duration",
            "Data": "Test data for chart"
        },
        {
            "ID": 2,
            "UserId": 1,
            "Favourite": true,
            "Title": "Movies Chart",
            "XAxes": "Movies",
            "YAxes": "#Views",
            "Data": "Test data for chart"
        }
    ],
    "insights": [
        {
            "ID": 1,
            "UserId": 1,
            "Favourite": true,
            "Text": "40% of millenials spend more than 3hours on social media daily"
        },
        {
            "ID": 2,
            "UserId": 1,
            "Favourite": true,
            "Text": "50% of females in age 30-35 spend more than 8hours on running weekly"
        }
    ],
    "audiences": [
        {
            "ID": 1,
            "UserId": 1,
            "Favourite": true,
            "Gender": "Female",
            "Country": "Greece",
            "AgeFrom": 30,
            "AgeTo": 35,
            "SocialHours": 2,
            "Purchases": 80
        },
        {
            "ID": 2,
            "UserId": 1,
            "Favourite": true,
            "Gender": "Male",
            "Country": "Greece",
            "AgeFrom": 30,
            "AgeTo": 35,
            "SocialHours": 5,
            "Purchases": 70
        }
    ]
    }

## Create chart
`POST api/assets/charts/{user_id}`
### Request

    curl --location --request POST 'http://localhost:8000/api/assets/charts/1'
    --header 'Content-Type: application/json' 
    --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ'
    --data-raw  ' 
    {
      "UserId": 1,
      "Title":"Movies Chart",
      "XAxes":"Movies",
      "YAxes":"#Views",
      "Data":"Test data for chart"
    }'
  
### Response
Returns an integer for the new chart.New records are marked as favourite by default.

## Create insight
`POST api/assets/insights/{user_id}`
### Request

    curl --location --request POST 'http://localhost:8000/api/assets/insights/1'
    --header 'Content-Type: application/json' 
    --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ'
    --data-raw  ' 
    {
      "UserId": 1,
      "Text":"40% of millenials spend more than 3hours on social media daily"
    }'
  
### Response
Returns an integer for the new insight. New records are marked as favourite by default.

## Create audience
`POST api/assets/audiences/{user_id}`
### Request

    curl --location --request POST 'http://localhost:8000/api/assets/audiences/1'
    --header 'Content-Type: application/json' 
    --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ'
    --data-raw  ' 
    {
            "UserId": 1,
            "Favourite": true,
            "Gender": "Female",
            "Country": "Greece",
            "AgeFrom": 30,
            "AgeTo": 35,
            "SocialHours": 2,
            "Purchases": 80
        }'
  
### Response
Returns an integer for the new audience. New records are marked as favourite by default.

## Update chart
Mark or unmark as favourite a chart by setting Favourite:false 
`PUT api/assets/charts`
### Request
    curl --location --request PUT 'http://localhost:8000/api/assets/charts/1'
    --header 'Content-Type: application/json' 
    --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ'
    --data-raw  ' 
    {
      "ID":1,
      "UserId": 1,
      "Favourite": false
    }'
  
### Response
Returns the string "Chart has been updated" or error message ""

## Update insight
### Request
Mark or unmark as favourite an insight by setting Favourite:false 
`PUT api/assets/insights/{user_id}`

    curl --location --request PUT 'http://localhost:8000/api/assets/insights/1'
    --header 'Content-Type: application/json' 
    --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ'
    --data-raw  ' 
    {
      "ID":1,
      "UserId": 1,
      "Favourite": false
    }'
  
### Response
Returns the string "Insight has been updated".

## Update audience
Mark or unmark as favourite an audence by setting Favourite:false 
`PUT api/assets/audiences/{user_id}`
### Request
    curl --location --request PUT 'http://localhost:8000/api/assets/audiences/1'
    --header 'Content-Type: application/json' 
    --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBkaW1pcm9wb3Vsb3VAbWFpbC5jb20ifQ.6mtFYiH-Fmyc1ybc6H-3PcXXFUW6giYWs3knD5R0UKQ'
    --data-raw  ' 
    {
            "ID": 1,
            "UserId": 1,
            "Favourite": false,
        }'
  
### Response
Returns the string "Audience has been updated".
