#!/usr/bin/bash

set -eu

mongosh -- "${MONGO_INITDB_DATABASE}" <<EOF

    const MONGO_INITDB_ROOT_USERNAME = '$MONGO_INITDB_ROOT_USERNAME';
    const MONGO_INITDB_ROOT_PASSWORD = '$MONGO_INITDB_ROOT_PASSWORD';

    db.createCollection('users', {
        capped = false
    });

    db.users.insertMany([
        {
            "userID": "user42",
            "firstName": "Fotis",
            "lastName": "Michalopoulos",
            "age": 33
        },
        {
            "userID": "user47",
            "firstName": "Joe",
            "lastName": "Rice",
            "age": 35
        }
    ])

    db.createCollection('assets', {
        capped: false
    });

    db.assets.insertMany([
        {
            "_id": "3fa85f64-5717-4562-b3fc-2c963f66afa7",
            "type": "insight",
            "text": "40% of millenials spend more than 3 hours on social media daily"
        },
        {
            "_id": "d9a972b5-2783-4b72-9cec-5f3031d61ab5",
            "type": "insight",
            "text": "60% of millenials spend more than 3 hours on their mobiles"
        },
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
    ])

    db.createCollection('favourites', {
        capped: false
    });

    db.favourites.insertMany([
        {
            "userID": "user42",
            "assetID": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "assetType": "chart",
            "description": "This my favourite chart"
        },
        {
            "userID": "user42",
            "assetID": "3fa85f64-5717-4562-b3fc-2c963f66afa7",
            "assetType": "insight",
            "description": "Useful insight"
        }
    ])

    

EOF
