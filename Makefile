.PHONY: test
test:
	go test -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o service ./cmd

.PHONY: build
build:
	go build -o service ./cmd

.PHONY: build-docker
build-docker:
	docker build -f Dockerfile-build -t service --progress=plain .

.PHONY: pack-docker
pack-docker: build-linux
	docker build -f Dockerfile-static -t service --progress=plain .

