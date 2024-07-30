#!/bin/sh

apt-get update
apt-get install -y iproute2

#echo 'Waiting for Postgres to be ready...'
#while ! ss -tuln | grep -q ':5432'; do
#    sleep 4
#done

echo 'Postgres is ready, running migrations'
go install github.com/pressly/goose/v3/cmd/goose@latest
export PATH=$PATH:$GOPATH/bin
goose -dir=/app/internal/db/migrations postgres "user=postgres password=qwert dbname=microserv host=postgres sslmode=disable" up
