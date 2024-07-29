# Задержка для ожидания полной инициализации базы данных
sleep 3

# Замена строки в конфигурационном файле
sed -i 's/^#wal_level = replica/wal_level = logical/' /var/lib/postgresql/data/postgresql.conf

# Остановка PostgreSQL
/usr/lib/postgresql/16/bin/pg_ctl -D /var/lib/postgresql/data stop -s -m fast

sleep 5

# Запуск PostgreSQL
/usr/lib/postgresql/16/bin/pg_ctl -D /var/lib/postgresql/data start -s
sleep 5

# Проверка состояния PostgreSQL
pg_isready
