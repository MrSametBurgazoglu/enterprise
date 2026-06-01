 generate-models:
	go run generate/generate.go
	go fmt ./...
run-tests:
	go test -v
start-tests: generate-models run-tests

integration-test: generate-models
	docker compose -f ../docker-compose.yml up -d postgres
	@echo "Waiting for postgres to be ready..."
	@until docker compose -f ../docker-compose.yml exec -T postgres pg_isready -h localhost -p 5432 -U postgres >/dev/null 2>&1; do sleep 1; done
	@echo "Running integration tests..."
	@POSTGRES_URL="postgres://postgres:password@localhost:5432/enterprise?sslmode=disable&search_path=public" go test -v -run=TestIntegration ./...; \
	status=$$?; \
	docker compose -f ../docker-compose.yml down; \
	exit $$status