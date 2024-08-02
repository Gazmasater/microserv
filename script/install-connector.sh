#!/bin/bash

# Wait for Kafka Connect to be fully up and running
echo "Waiting for Kafka Connect to start..."
until curl -s http://kafka-connect:8083/ | grep -q '"version"'; do
  sleep 5
done

# Install the connector
echo "Installing the connector..."
curl -X POST -H "Content-Type: application/json" --data @/usr/share/confluent/docker/postgres-connector2.json http://kafka-connect:8083/connectors

# Wait a bit for the connector to be installed
sleep 3

# Check installation status
echo "Connector installation status:"
curl -X GET http://kafka-connect:8083/connectors/postgres-connector2/status
