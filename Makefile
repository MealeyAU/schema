.PHONY: all
all:
	cd ./generate && go run ./cmd/main.go --all

.PHONY: go
go:
	cd ./generate && go run ./cmd/main.go --go

.PHONY: web
web:
	cd ./generate && go run ./cmd/main.go --web