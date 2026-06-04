#!/bin/bash
set -e

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$BASE_DIR/scraper"
go run .

cd "$BASE_DIR/analysis_alert"
source .venv/bin/activate
python3 analysis.py
python3 alert.py
deactivate
