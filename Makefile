.PHONY: all

all: shell tidy build

shell:
	@echo 'SHELL='$(SHELL)

tidy:
	go mod tidy

build:
	go vet ./cmd/mm

	CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mm ./cmd/mm

	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mm_linux_amd64 ./cmd/mm

	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mm_linux_arm64 ./cmd/mm

	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mm_darwin_amd64 ./cmd/mm

	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mm_darwin_arm64 ./cmd/mm

	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ./releases/mm_windows_amd64.exe ./cmd/mm