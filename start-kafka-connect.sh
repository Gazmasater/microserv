#!/bin/bash

# Подождем, пока Kafka Connect не станет доступен
echo "Waiting for Kafka Connect to start..."
sleep 20

# Запускаем Kafka Connect
echo "Starting Kafka Connect..."
/usr/local/bin/connect-distributed /etc/kafka/connect-distributed.properties
