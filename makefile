GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

BINARY_NAME = website-analyzer

MAIN_PATH = main.go

GO_VERSION = 1.21.7
GOPROXY = https://proxy.golang.org,direct

all: clean build

build:
	@echo ">> Building binary..."
	@GOPROXY=$(GOPROXY) $(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

test:
	@echo ">> Running tests..."
	@$(GOTEST) -v ./...

clean:
	@echo ">> Cleaning up..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)

run:
	@echo ">> Running binary..."
	@GOPROXY=$(GOPROXY) $(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	@./$(BINARY_NAME)

docker-build:
	@echo ">> Building Docker image..."
	@docker build -t website-analyzer .

docker-run:
	@echo ">> Running Docker container..."
	@docker run -p 8080:8080 website-analyzer

help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@echo "  all            Clean and build"
	@echo "  build          Build the binary"
	@echo "  test           Run tests"
	@echo "  clean          Clean up"
	@echo "  run            Build and run the binary"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container"
	@echo "  help           Show this help message"

.PHONY: all build test clean run docker-build docker-run help
