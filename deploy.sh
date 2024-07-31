#!/bin/bash


# Запустить контейнеры в фоновом режиме
echo "Запуск контейнеров..."
docker-compose up -d

chmod +x install-connector.sh

# Проверка, доступен ли порт 8083
echo "Ожидание запуска Kafka Connect на порту 8083..."
until curl -s http://localhost:8083/ | grep -q '"version"'; do
  sleep 5
done

# Выполнить скрипт установки коннектора
echo "Запуск скрипта install-connector.sh..."
docker exec -it kafka-connect /usr/share/confluent/docker/install-connector.sh

echo "Скрипт завершен."
