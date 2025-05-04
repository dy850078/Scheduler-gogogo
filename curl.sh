#!/bin/sh

for i in {1..5}; do
  curl -X POST http://localhost:8080/schedule \
    -H "Content-Type: application/json" \
    -d '{
      "requested_cpu": 4,
      "requested_memory": 8192,
      "requested_pool": "default",
      "dedicated": false
    }'
  echo ""  # 換行
done
