#!/bin/bash

# Убедитесь, что Kafka Connect полностью запущен
echo "Ожидание запуска Kafka Connect..."
until curl -s http://kafka-connect:8083/ | grep -q '"version"'; do
  sleep 5
done

# Установка коннектора
echo "Установка коннектора..."
curl -X POST -H "Content-Type: application/json" --data @/usr/share/confluent/docker/postgres-connector2.json http://kafka-connect:8083/connectors

# Проверка статуса установки
echo "Статус коннектора:"
curl -X GET http://kafka-connect:8083/connectors/postgres-connector/status
