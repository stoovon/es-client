initialise-linux:
	sudo sysctl -w vm.max_map_count=262144

up:
	docker-compose up

get_payments:
	curl -X GET 'http://localhost:9200/payments/_search?pretty=true'
	@echo ""

truncate_payments:
	curl -X DELETE 'http://localhost:9200/payments'
	@echo ""

get_indexes:
	curl -X GET 'http://localhost:9200/_cat/indices?v&pretty=true'
	@echo ""

PAYMENT_DATA = $(shell cat payment.json)

put_payment:
	curl -X PUT --header "Content-Type: application/json" \
		http://localhost:9200/payments/_doc/00000000-0000-0000-0000-000000000000 \
		-d '$(PAYMENT_DATA)'
	@echo ""

.PHONY: initialise-linux up

.DEFAULT_GOAL := up
