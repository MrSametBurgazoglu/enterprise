generate-models:
	cd example
	go run generate/generate.go
run-example:
	cd example
	go run cmd/store.go
migrate:
	cd example
	go run migrate/migrate.go