#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/00_env.sh"

TOTAL="${TOTAL:-300}"          # total favourites to create
CHART_PCT="${CHART_PCT:-20}"   # percent charts
INSIGHT_PCT="${INSIGHT_PCT:-60}" # percent insights
# remaining percent becomes audience

countries=("US" "GB" "DE" "FR" "GR" "ES" "IT" "NL" "SE" "PL")
genders=("male" "female")

rand_country() { echo "${countries[$((RANDOM % ${#countries[@]}))]}"; }
rand_gender() { echo "${genders[$((RANDOM % ${#genders[@]}))]}"; }

echo "Seeding realistic dataset: TOTAL=$TOTAL user=$USER_ID at $API_URL"

for i in $(seq 1 "$TOTAL"); do
  r=$((RANDOM % 100))
  asset_id="asset-$i"

  if [ "$r" -lt "$INSIGHT_PCT" ]; then
    # Insight
    hours=$((1 + RANDOM % 6))
    pct=$((10 + RANDOM % 71)) # 10..80
    body=$(jq -nc \
      --arg id "$asset_id" \
      --arg text "$pct% of millennials spend more than $hours hours on social media daily" \
      --arg desc "Random insight #$i" '
      {asset_id:$id, type:"insight", description:$desc, text:$text}
    ')

  elif [ "$r" -lt $((INSIGHT_PCT + CHART_PCT)) ]; then
    # Chart
    title="Monthly Revenue (FY2025)"
    x="Month"
    y="Revenue"

    # Create 6 random revenue points in bash
    jan=$((800 + RANDOM % 601))
    feb=$((800 + RANDOM % 601))
    mar=$((800 + RANDOM % 601))
    apr=$((800 + RANDOM % 601))
    may=$((800 + RANDOM % 601))
    jun=$((800 + RANDOM % 601))

    data=$(jq -nc \
      --argjson jan "$jan" \
      --argjson feb "$feb" \
      --argjson mar "$mar" \
      --argjson apr "$apr" \
      --argjson may "$may" \
      --argjson jun "$jun" '
      [
        {month:"Jan", revenue:$jan},
        {month:"Feb", revenue:$feb},
        {month:"Mar", revenue:$mar},
        {month:"Apr", revenue:$apr},
        {month:"May", revenue:$may},
        {month:"Jun", revenue:$jun}
      ]
    ')

    body=$(jq -nc \
      --arg id "$asset_id" \
      --arg title "$title" \
      --arg x "$x" \
      --arg y "$y" \
      --arg desc "Random chart #$i" \
      --argjson data "$data" '
      {asset_id:$id, type:"chart", description:$desc, title:$title, x_axis_title:$x, y_axis_title:$y, data:$data}
    ')

  else
    # Audience
    c=$(rand_country)
    g=$(rand_gender)
    min=$((18 + (RANDOM % 30)))
    max=$((min + 5 + (RANDOM % 10)))
    hours=$((1 + RANDOM % 6))
    purchases=$((RANDOM % 10))

    name="$g $min-$max in $c"
    body=$(jq -nc \
      --arg id "$asset_id" \
      --arg desc "Random audience #$i" \
      --arg name "$name" \
      --arg gender "$g" \
      --arg country "$c" \
      --argjson min "$min" \
      --argjson max "$max" \
      --argjson hours "$hours" \
      --argjson purchases "$purchases" '
      {
        asset_id:$id,
        type:"audience",
        description:$desc,
        name:$name,
        gender:$gender,
        birth_country:$country,
        age_group:{min:$min,max:$max},
        daily_social_media_hours:$hours,
        purchases_last_month:$purchases
      }
    ')
  fi

  code=$(curl -sS -o /dev/null -w "%{http_code}" \
    -X POST "$API_URL/api/v1/users/$USER_ID/favourites" \
    -H "Content-Type: application/json" \
    -d "$body")

  if [ "$code" != "201" ]; then
    echo "Seed failed i=$i status=$code body=$body"
    exit 1
  fi
done

echo "Done."
