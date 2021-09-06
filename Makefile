.PHONY: all
all:
	go run ./cmd/main.go --all

.PHONY: go
go:
	go run ./cmd/main.go --go

.PHONY: web
web:
	go run ./cmd/main.go --web