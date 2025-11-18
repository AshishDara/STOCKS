#!/usr/bin/env bash

set -euo pipefail

echo ">> Installing frontend dependencies"
npm install

echo ">> Building GrowTrade frontend"
npm run build

echo "Build complete. Files are available in $(pwd)/dist"

