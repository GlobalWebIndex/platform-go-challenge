#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/00_env.sh"

EXPECTED="${EXPECTED:-0}"

echo "Walking pagination for user=$USER_ID limit=$LIMIT"
echo

cursor=""
page=0
total_seen=0
declare -A seen_ids

while true; do
  if [ -z "$cursor" ]; then
    url="$API_URL/api/v1/users/$USER_ID/favourites?limit=$LIMIT"
  else
    url="$API_URL/api/v1/users/$USER_ID/favourites?limit=$LIMIT&cursor=$cursor"
  fi

  response=$(curl -sS "$url")

  items_count=$(echo "$response" | jq '.data | length')
  has_more=$(echo "$response" | jq -r '.meta.has_more')
  next_cursor=$(echo "$response" | jq -r '.meta.next_cursor // ""')

  page=$((page+1))
  total_seen=$((total_seen + items_count))

  echo "Page $page:"
  echo "  items=$items_count has_more=$has_more"

  if [ "$items_count" -gt 0 ]; then
    first_id=$(echo "$response" | jq -r '.data[0].asset_id')
    last_id=$(echo "$response" | jq -r '.data[-1].asset_id')
    echo "  first_asset_id=$first_id"
    echo "  last_asset_id=$last_id"
  fi

  # Duplicate detection
  for id in $(echo "$response" | jq -r '.data[].asset_id'); do
    if [[ -n "${seen_ids[$id]:-}" ]]; then
      echo "Duplicate detected across pages: $id"
      exit 1
    fi
    seen_ids[$id]=1
  done

  # End conditions
  if [ "$has_more" != "true" ]; then
    break
  fi

  if [ -z "$next_cursor" ]; then
    echo "has_more=true but next_cursor empty"
    exit 1
  fi

  if [ "$next_cursor" == "$cursor" ]; then
    echo "Cursor did not advance"
    exit 1
  fi

  cursor="$next_cursor"
done

echo
echo "Total seen: $total_seen"

if [ "$EXPECTED" -gt 0 ]; then
  if [ "$EXPECTED" -ne "$total_seen" ]; then
    echo "Expected $EXPECTED but got $total_seen"
    exit 1
  fi
  echo "Matches expected total: $EXPECTED"
fi

echo "Pagination traversal successful"
