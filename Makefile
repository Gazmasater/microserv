CONFIG_FILE := ./connector_config.json

# Команда для создания коннектора
create-connector:
	@echo "Using config file: $(CONFIG_FILE)"
	curl -v -X POST -H "Content-Type: application/json" --data @$(CONFIG_FILE) http://localhost:8083/connectors

# По умолчанию выводим справку о командах
.DEFAULT_GOAL := help

help:
	@echo "Доступные команды:"
	@echo "  make create-connector  - создать коннектор Debezium в Kafka Connect"

.PHONY: help create-connector
