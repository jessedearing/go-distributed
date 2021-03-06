all: test

.PHONY: mongotest
mongotest: ## Test mongo integration
	go test -tags mongo ./...

.PHONY: mysql
mysql: ## Test mysql integration
	go test -tags mysql ./...

.PHONY: test
test: mongotest ## Run all tests

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
