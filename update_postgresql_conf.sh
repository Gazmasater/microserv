#!/bin/bash

# Определите старую и новую подстроки
old_substring='wal_level = replica'
new_substring='wal_level = logical'

# Имя контейнера
container_name='microserv-postgres-1'

# Путь к файлу конфигурации внутри контейнера
conf_file='/var/lib/postgresql/data/postgresql.conf'

# Проверка существования контейнера
if ! docker ps -q --filter "name=$container_name"; then
  echo "Контейнер $container_name не найден или не запущен."
  exit 1
fi

# Замена подстроки в файле конфигурации внутри контейнера
docker exec $container_name sed -i "s|$old_substring|$new_substring|g" $conf_file

# Проверка успешности замены подстроки
if [ $? -eq 0 ]; then
  echo "Подстрока успешно заменена."
else
  echo "Ошибка при замене подстроки."
  exit 1
fi

# Перезапуск контейнера для применения изменений
docker restart $container_name

# Проверка успешности перезапуска контейнера
if [ $? -eq 0 ]; then
  echo "Контейнер $container_name успешно перезапущен."
else
  echo "Ошибка при перезапуске контейнера."
  exit 1
fi
