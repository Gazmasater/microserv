Клонирование проекта
git clone https://github.com/Gazmasater/microserv.git

Перейти на ветку new6

Запуск сервиса
./deploy.sh

После вывода статуса коннектора открываем сваггер
http://localhost:8080/swagger/index.html

Действия с коннектором

curl -X DELETE http://localhost:8083/connectors/postgres-connector2


curl -X POST -H "Content-Type: application/json" --data @postgres-connector2.json http://localhost:8083/connectors


curl -s -X GET http://localhost:8083/connectors/postgres-connector2/status

