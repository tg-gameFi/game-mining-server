ENV := $(if ${env},${env},"local")
verbose := "false"
server := "tgw-server"

.PHONY: docker-up
docker-up: ## Start docker
	cd zlocal && docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop docker
	cd zlocal && docker-compose down

.PHONY: redis-cli
redis-cli: ## Enter redis-cli
	# /opt/bitnami/redis/bin/redis-cli
	# AUTH redis123
	cd zlocal && docker-compose exec redis bash

.PHONY: build
build: ## Build application and plugins
	go mod tidy
	cd zlocal && ./build.sh

.PHONY: serve
serve:  ## Run application
	@echo "run in ${ENV} env"
	cd zlocal && ./bin/${server} ${ENV}

.PHONY: close
close: docker-down ## close application
	rm -r zlocal/bin