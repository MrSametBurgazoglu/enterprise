 generate-models:
	go run generate/generate.go
	go fmt ./...
run-tests:
	go test -v
start-tests: generate-models run-tests