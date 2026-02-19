#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/00_env.sh"

LIMIT="${LIMIT:-50}"

health_cache() {
  curl -sS "$API_URL/api/v1/health" | jq '.cache'
}

echo "Baseline cache:"
health_cache

echo "GET favourites (warm cache - expect miss then set):"
curl -sS "$API_URL/api/v1/users/$USER_ID/favourites?limit=$LIMIT" > /dev/null
health_cache

echo "GET favourites again (expect hit):"
curl -sS "$API_URL/api/v1/users/$USER_ID/favourites?limit=$LIMIT" > /dev/null
health_cache

echo "Write (POST) - expect invalidate/delete:"
body=$(jq -nc --arg id "asset-cache-new" '{asset_id:$id, type:"insight", text:"invalidate"}')
curl -sS -X POST "$API_URL/api/v1/users/$USER_ID/favourites" \
  -H "Content-Type: application/json" \
  -d "$body" > /dev/null
health_cache

echo "GET favourites after write (expect miss again):"
curl -sS "$API_URL/api/v1/users/$USER_ID/favourites?limit=$LIMIT" > /dev/null
health_cache
