#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Define the required file path relative to the script's location
REQUIRED_FILE="$SCRIPT_DIR/../docs/swagger-ready.yaml"

LOGS_FILE="$SCRIPT_DIR/../mounts/backend/logs/application.log"

# Ensure the file exists
if [ -f "$REQUIRED_FILE" ]; then
    echo "File exists. Removing: $REQUIRED_FILE"
    rm "$REQUIRED_FILE"
fi

echo "Touching: $REQUIRED_FILE"
touch "$REQUIRED_FILE"

# Ensure the file exists
if [ -f "$LOGS_FILE" ]; then
    echo "File exists. Removing: $LOGS_FILE"
    rm "$LOGS_FILE"
fi

echo "Touching: $LOGS_FILE"
touch "$LOGS_FILE"
