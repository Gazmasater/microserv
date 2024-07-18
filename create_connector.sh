#!/bin/bash

# Ожидание запуска Kafka Connect
sleep 20

# Создание коннектора
curl -X POST -H "Content-Type: application/json" -d '{
  "name": "postgres-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "tasks.max": "1",
    "database.hostname": "postgres",
    "database.port": "5432",
    "database.user": "postgres",
    "database.password": "qwert",
    "database.dbname": "microserv",
    "database.server.name": "dbserver1",
    "table.include.list": "public.msg",
    "plugin.name": "pgoutput",
    "topic.prefix": "dbserver1.",
    "bootstrap.servers": "kafka:9092"
  }
}' http://connect:8083/connectors
