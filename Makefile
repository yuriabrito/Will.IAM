testable_packages=$(shell go list ./... | egrep -v 'constants|mocks|testing')
project=$(shell basename $(PWD))
project_test=${project}-test
project_sanitized=$(shell echo $(project) | sed -e "s/\-//")
project_sanitized_test=$(shell echo $(project)-test | sed -e "s/\-//")
pg_dep=$(project)_postgres_1
test_packages=`find . -type f -name "*.go" ! \( -path "*vendor*" \) | sed -En "s/([^\.])\/.*/\1/p" | uniq`
database=postgres://postgres:$(project)@localhost:8432/$(project)?sslmode=disable
database_test=postgres://postgres:$(project)@localhost:8432/$(project_test)?sslmode=disable

setup: setup-gin setup-project

setup-gin:
	@go get github.com/codegangsta/gin

setup-project:
	@go get -u github.com/golang/dep...
	@dep ensure

deps:
	@mkdir -p docker_data && docker-compose up -d postgres
	@until docker exec $(pg_dep) pg_isready; do echo 'Waiting Postgres...' && sleep 1; done
	@docker exec $(pg_dep) createuser -s -U postgres $(project) 2>/dev/null || true
	@docker exec $(pg_dep) createdb -U $(project) $(project) 2>/dev/null || true

deps-test:
	@mkdir -p docker_data && docker-compose up -d postgres
	@until docker exec $(pg_dep) pg_isready; do echo 'Waiting Postgres...' && sleep 1; done
	@docker exec $(pg_dep) createuser -s -U postgres $(project) 2>/dev/null || true
	@docker exec $(pg_dep) createdb -U $(project) $(project_test) 2>/dev/null || true
	@make migrate-test

stop-deps:
	@docker-compose down

stop-deps-test:
	@make drop-test
	@make stop-deps

build:
	@mkdir -p bin && go build -o ./bin/$(project) .

build-docker:
	@docker build -t $(project) .

run:
	@gin -i --port 3001 --appPort 4040 --bin Will.IAM run start-api

setup-migrate:
	@go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
	@go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

migrate:
	@migrate -path migrations -database ${database} up

migrate-test:
	@migrate -path migrations -database ${database_test} up

drop:
	@migrate -path migrations -database ${database} drop

drop-test:
	@migrate -path migrations -database ${database_test} drop

test:
	@make deps-test
	@make test-fast
	@make stop-deps-test

test-fast:
	@make drop-test
	@make migrate-test
	@make unit
	@make integration

unit:
	@echo "Unit Tests"
	@go test ${testable_packages} -tags=unit -coverprofile unit.coverprofile -v
	@make gather-unit-profiles

integration:
	@echo "Integration Tests"
	@ret=0 && for pkg in ${testable_packages}; do \
		echo $$pkg; \
		go test $$pkg -tags=integration -coverprofile integration.coverprofile -v; \
		test $$? -eq 0 || ret=1; \
	done; exit $$ret
	@make gather-integration-profiles

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'
	@find . -name "*.coverprofile" -delete

gather-integration-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-integration.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-integration.out; done'
	@find . -name "*.coverprofile" -delete
