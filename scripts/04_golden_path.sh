#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/00_env.sh"

asset_id="asset-golden-1"

echo "POST favourite"
body=$(jq -nc --arg id "$asset_id" '
  {asset_id:$id, type:"insight", description:"initial", text:"hello"}
')
curl -sS -X POST "$API_URL/api/v1/users/$USER_ID/favourites" \
  -H "Content-Type: application/json" \
  -d "$body" | jq .

echo "GET favourites"
curl -sS "$API_URL/api/v1/users/$USER_ID/favourites?limit=20" | jq '.meta, (.data[] | select(.asset_id=="'"$asset_id"'"))'

echo "GET favourite"
curl -sS "$API_URL/api/v1/users/$USER_ID/favourites/$asset_id" | jq .

echo "PUT description"
upd=$(jq -nc '{description:"updated desc"}')
curl -sS -o /dev/null -w "status=%{http_code}\n" \
  -X PUT "$API_URL/api/v1/users/$USER_ID/favourites/$asset_id/description" \
  -H "Content-Type: application/json" \
  -d "$upd"

echo "GET favourite (verify description)"
curl -sS "$API_URL/api/v1/users/$USER_ID/favourites/$asset_id" | jq '.data.description'


echo "DELETE favourite"
curl -sS -o /dev/null -w "status=%{http_code}\n" \
  -X DELETE "$API_URL/api/v1/users/$USER_ID/favourites/$asset_id"
