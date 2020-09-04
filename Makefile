
build: author-manager book-manager

author-manager:
	GOCACHE='/tmp/.cache' CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o ./bin/author-manager -v ./cmd/author-manager/main.go
book-manager:
	GOCACHE='/tmp/.cache' CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o ./bin/book-manager -v ./cmd/book-manager/main.go
