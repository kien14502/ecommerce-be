#!/bin/bash
set -e

echo "⏳ Waiting for Kafka Connect..."

until curl -sf http://connect:8083/connectors >/dev/null; do
  sleep 5
done

echo "✅ Kafka Connect is ready"

# Check if connector already exists
if curl -sf http://connect:8083/connectors/inventory-connector >/dev/null; then
  echo "ℹ️  Connector already exists, skipping registration"
  exit 0
fi

echo "🚀 Registering Debezium MySQL connector..."

curl -sf -X POST http://connect:8083/connectors \
  -H "Content-Type: application/json" \
  -d @/scripts/register-mysql.json

echo "🎉 Debezium connector registered successfully!"
