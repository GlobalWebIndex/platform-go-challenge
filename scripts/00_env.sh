#!/usr/bin/env bash
set -euo pipefail

API_URL="${API_URL:-http://localhost:8080}"
USER_ID="${USER_ID:-user-123}"
LIMIT="${LIMIT:-50}"

require_cmd() { command -v "$1" >/dev/null 2>&1 || { echo "Missing dependency: $1"; exit 1; }; }
require_cmd curl
require_cmd jq

export API_URL USER_ID LIMIT