#!/bin/bash

# Config
PORT=3000
ORIGIN="http://dummyjson.com"
CMD="go run ./cmd"
PROXY_URL="http://localhost:$PORT"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Launching Catching Proxy port $PORT pointing at $ORIGIN...${NC}"
$CMD --port $PORT --origin $ORIGIN &
PID=$!

# Esperamos que el servidor arranque
sleep 2

echo -e "${YELLOW}\n📦 First request (MISS expected)...${NC}"
curl -i "$PROXY_URL/products"

echo -e "${YELLOW}\n📦 Second request (HIT expected)...${NC}"
curl -i "$PROXY_URL/products"

echo -e "${YELLOW}\n📊 Cache statistics...${NC}"
curl -s "$PROXY_URL/cache/stats" | jq

echo -e "${YELLOW}\n🧹 Cleaning cache via HTTP...${NC}"
curl -s -X POST "$PROXY_URL/cache/clear" | jq

echo -e "${YELLOW}\n📊 Statistics after cleaning...${NC}"
curl -s "$PROXY_URL/cache/stats" | jq

echo -e "${GREEN}\n🛑 Server stopped...${NC}"
kill $PID
