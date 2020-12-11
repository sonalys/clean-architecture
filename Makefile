run:
	go run api/main.go
test:
	go test ./...
coverage:
	go test ./... -coverprofile test_result.html
build:
	go mod vendor && \
	cd api && \
	CGO_ENABLED=0 go build -mod=vendor -a -o ./build/api -ldflags '-w'