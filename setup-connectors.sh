#!/bin/bash

# Подождем, пока Kafka Connect не станет доступен
echo "Waiting for Kafka Connect to start..."
sleep 20

# Определяем путь к файлу конфигурации коннектора
CONNECTOR_FILE="/kafka/connect/postgres-connector2.json"

# Проверяем, существует ли файл
if [ -f "$CONNECTOR_FILE" ]; then
  echo "Loading connector configuration from $CONNECTOR_FILE"
  # Отправляем POST-запрос на загрузку конфигурации коннектора
  curl -X POST http://kafka-connect:8083/connectors \
       -H "Content-Type: application/json" \
       -d @"$CONNECTOR_FILE"
else
  echo "Connector configuration file $CONNECTOR_FILE not found!"
fi
