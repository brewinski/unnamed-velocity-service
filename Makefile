build:
	@go build -o ./bin/server main.go

run: build
	@./bin/server

watch:
	@reflex -s -r '\.go$$' make run