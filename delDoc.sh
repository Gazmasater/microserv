#!/bin/bash

echo "Остановка всех контейнеров..."
docker stop $(docker ps -aq)

echo "Удаление всех контейнеров..."
docker rm $(docker ps -aq)

echo "Удаление всех образов..."
docker rmi $(docker images -q)

echo "Удаление всех неиспользуемых томов..."
docker volume rm $(docker volume ls -q)

echo "Удаление всех сетей..."
docker network rm $(docker network ls -q)

echo "Очистка системы от неиспользуемых данных..."
docker system prune -af --volumes

echo "Все контейнеры, образы, тома и сети были удалены."
