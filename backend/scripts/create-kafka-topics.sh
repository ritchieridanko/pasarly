#!/bin/bash

# Wait for Kafka brokers
until /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka1:9092 --list >/dev/null 2>&1; do
  sleep 1
done

# Create Kafka topics
/opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka1:9092 --create --if-not-exists \
  --topic auth.created --partitions 3 --replication-factor 3

echo "✅ [BROKER] topics created"

exit 0
