# Collections and Documents for the GWI solution

## Table Of Contents
- [Users Collection](#users-collection)
    * [User Document](#user-document)
    * [User Example](#user-example)
- [Assets Collection](#assets-collection)
    * [Asset Document](#asset-document)
        - [Insight](#insight)
        - [Insight Example](#insight-example)
        - [Chart](#chart)
        - [Chart Example](#chart-example)
        - [Audience](#audience)
        - [Audience Example](#audience-example)
- [Favourites Collection](#favourites-collection)
    * [Favourite Document](#favourite-document)
    * [Favourite Example](#favourite-example)

## Users Collection

### User Document

1. `userID`:
    - Type: `String`
    - Example: `user42`
    - Description: unique identifier of a user 

2.  `firstName`
    - Type: `String`
    - Example: `John`

3. `lastName` 
    - Type: `String`
    - Example: `Doe`
4. `age` 
    - Type: `Int`
    - Example: 33

### User Example

```json
{
    "userID": "user42",
    "firstName": "Fotis",
    "lastName": "Michalopoulos",
    "age": 33
},
```

## Assets Collection

### Asset Document

An `asset` document can be one of the following:

#### Insight
1. `_id`:
    - Type: `String`
    - Format: uuid
    - Description: unique identifier of an asset
2. `type`:
    - Type: `String`
    - Value: insight
3. `text`
    - Type: `String`
    - Description: contains info on the insight  

#### Insight Example

```json
{
    "_id": "3fa85f64-5717-4562-b3fc-2c963f66afa7",
    "type": "insight",
    "text": "40% of millenials spend more than 3 hours on social media daily"
},
```

#### Chart
1. `_id`:
    - Type: `String`
    - Format: uuid
    - Description: unique identifier of a chart
2. `type`:
    - Type: `String`
    - Value: chart 
3. `title`:
    - Type: `String`
4. `xAxis`:
    - Type: String
    - Description: Name of the xAxis
5. `yAxis`:
    - Type: String
    - Description: Name of the yAxis
6. `data`:
    - Type: Array
    - Items: {"x": "value", "y": "value"}

#### Chart Example

```json
{
    "_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "type": "chart",
    "title": "Revenue over time",
    "xAxis": "Months",
    "yAxis": "Revenue ($)",
    "data": [
        {"x": "Jan", "y": 100},
        {"x": "Feb", "y": 1000},
        {"x": "Mar", "y": 2000}
    ]
}
```

#### Audience
1. `_id`:
    - Type: `String`
    - Format: uuid
    - Description: unique identifier of audience asset
2. `type`:
    - Type: `String`
    - Value: audience 
3. `gender`:
    - Type: `String`
    - enum: [`male`, `female`]
4. `birthCountry`:
    - Type: `String`
    - example: ATH
5. `ageGroups`:
    - Type: String
    - enum: [18-24,25-34, 35-44, 45-54, 55-64, 65+] 
6. `dailyHours`:
    - Type: Int
7. `purchases`:
    - Type: Int
    - Description: Num of purchases last month

#### Audience Example

```json
    "_id": "3fa85f64-5717-4562-b3fc-2c963f66afa8",
    "type": "audience",
    "gender": "25-34",
    "birthCountry": "ATH",
    "ageGroups": "adults",
    "dailyHours": 4,
    "purchases": 12
```


## Favourites Collection

### Favourite Document

1. `userID`:
    - Type: `String`
    - Example: `user42`
    - Description: unique identifier of a user 
2. `assetID`:
    - Type: `String`
    - Format: uuid
    - Description: unique identifier of an asset
3. `assetType`:
    - Type: `String`
    - enum: [`chart`, `insight`, `audience`]
4. `assetID`:
    - Type: `String`
    - Description: unique identifier of an asset

### Favourite Example

```json
{
    "userID": "user42",
    "assetID": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "assetType": "chart",
    "description": "This my favourite chart"
},

```  
