#!/bin/bash
set -e

# Задержка для ожидания полной инициализации базы данных
sleep 10

# Замена строки в конфигурационном файле
sed -i 's/^#wal_level = replica/wal_level = logical/' /var/lib/postgresql/data/postgresql.conf

# Перезагрузка PostgreSQL
pg_ctl -D /var/lib/postgresql/data restart

# Проверка состояния PostgreSQL
pg_isready

# Сообщение о завершении
echo "PostgreSQL initialization completed."
