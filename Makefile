.PHONY: dockers-down

@all: docker-compose-dev start-local-producer

docker-compose-dev:
	@echo "docker-compose-dev.yml run in --detached mode"
	docker-compose -f docker-compose-dev.yml up -d

start-local-producer:
	@echo "start-local-producer"
	cd local-producer; pwd && \
	go run main.go


# dockers-down:
# docker-compose down


