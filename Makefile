REDIS_ADDR ?= localhost:6380

.PHONY: test redis-up redis-down

test: redis-up
	@REDIS_ADDR=$(REDIS_ADDR) go test ./...; status=$$?; exit $$status

redis-up:
	docker compose up -d --wait redis-test

redis-down:
	docker compose down
