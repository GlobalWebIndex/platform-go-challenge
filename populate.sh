#!/bin/sh

uuiduser=$(uuidgen)
uuidchart=$(uuidgen)
uuidinsight=$(uuidgen)
uuidaudience=$(uuidgen)

echo "User -> $uuiduser -> Add chart asset as favourite"
json1=$(cat <<-END
  {
    "type": "chart",
    "chart": {
      "id": "$uuiduser",
      "title": "foobar",
      "axis_y_title": "foo",
      "axis_x_title": "bar"
    }
  }
END
)
echo $json1 | jq

sleep 2

echo "Response -> "
curl -s --request POST 'http://localhost:4567/user/'$uuiduser'/fav/add' \
-d '{
  "type": "chart",
  "chart": {
    "id": "'$uuidchart'",
    "title": "foobar",
    "axis_y_title": "foo",
    "axis_x_title": "bar"
  }
}' | jq

echo "#######################################################################################"
sleep 3

echo "User -> $uuiduser-> Add insight asset as favourite"
json2=$(cat <<-END
  {
  "type": "insight",
  "insight": {
    "id": "'$uuidinsight'",
    "description": "foobar"
    }
  }
END
)
echo $json2 | jq

sleep 2

echo "Response -> "
curl -s --request POST 'http://localhost:4567/user/'$uuiduser'/fav/add' \
-d '{
  "type": "insight",
  "insight": {
    "id": "'$uuidinsight'",
    "description": "foobar"
  }
}' | jq

echo "#######################################################################################"
sleep 3

echo "User -> $uuiduser-> Add audience asset as favourite"
json3=$(cat <<-END
  {
  "type": "audience",
  "audience": {
    "id": "'$uuidaudience'",
    "gender": "foobar",
    "born_country": "foo",
    "age_group": "bar",
    "daily_hours_social_media": "foo",
    "purchases_last_month": "bar"
    }
  }
END
)
echo $json3 | jq

sleep 2

echo "Response -> "
curl -s --request POST 'http://localhost:4567/user/'$uuiduser'/fav/add' \
-d '{
  "type": "audience",
  "audience": {
    "id": "'$uuidaudience'",
    "gender": "foobar",
    "born_country": "foo",
    "age_group": "bar",
    "daily_hours_social_media": "foo",
    "purchases_last_month": "bar"
  }
}' | jq

echo "#######################################################################################"
sleep 3


echo "User -> $uuiduser-> Retrieve favourites list"

curl -s --request GET 'http://localhost:4567/user/'$uuiduser'/fav/list'  | jq

echo "#######################################################################################"
sleep 5

echo "User -> $uuiduser-> Update audience asset in favourites"
json4=$(cat <<-END
  {
  "type": "audience",
  "audience": {
    "id": "'$uuidaudience'",
    "gender": "updated",
    "born_country": "foo_updated",
    "age_group": "bar_updated",
    "daily_hours_social_media": "foo_updated",
    "purchases_last_month": "bar_updated"
    }
  }
END
)
echo $json4 | jq

sleep 2

echo "Response -> "
curl -s --request PUT 'http://localhost:4567/user/'$uuiduser'/fav/edit' \
-d '{
  "type": "audience",
  "audience": {
    "id": "'$uuidaudience'",
    "gender": "updated",
    "born_country": "foo_updated",
    "age_group": "bar_updated",
    "daily_hours_social_media": "foo_updated",
    "purchases_last_month": "bar_updated"
  }
}'  | jq

echo "#######################################################################################"
sleep 5


echo "User -> $uuiduser-> Retrieve favourites list"
curl -s --request GET 'http://localhost:4567/user/'$uuiduser'/fav/list'  | jq

echo "#######################################################################################"
sleep 5


echo "User -> $uuiduser-> Delete audience asset in favourites"
json5=$(cat <<-END
  {
  "type": "audience",
  "audience": {
    "id": "'$uuidaudience'"
    }
  }
END
)
echo $json5 | jq

sleep 2

curl -s --request DELETE 'http://localhost:4567/user/'$uuiduser'/fav/delete' \
-d '{
  "type": "audience",
  "audience": {
    "id": "'$uuidaudience'"
  }
}'  | jq

echo "#######################################################################################"
sleep 5


echo "User -> $uuiduser-> Retrieve favourites list"
curl -s --request GET 'http://localhost:4567/user/'$uuiduser'/fav/list' | jq