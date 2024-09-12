# Define variables for reusability
IMAGE_NAME = rocky2015aaa/filestorageservice
CONTAINER_NAME = filestorageservice
PORT = 8080
VERSION := 1.0.0
BUILD := production
DATE := $(shell date +'%Y-%m-%d_%H:%M:%S')

# Check if Docker image exists
image_exists = $(shell docker images -q $(IMAGE_NAME):latest)

# Check if Docker container exists
container_exists = $(shell docker ps -aq -f name=$(CONTAINER_NAME))

# Target to build the Docker image if it doesn't already exist
build:
ifeq ($(image_exists),)
	@echo "Building Docker image: $(IMAGE_NAME):latest"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--no-cache \
		-t $(IMAGE_NAME):latest .
else
	@echo "Docker image $(IMAGE_NAME):latest already exists."
endif

# Target to run the Docker container if it's not already running
up:
ifeq ($(container_exists),)
	@echo "Running Docker container: $(CONTAINER_NAME)"
	docker-compose up -d
else
	@echo "Docker container $(CONTAINER_NAME) already exists."
endif

# Target to build and run only if it’s the first time (initial setup)
setup: build up

# Target to stop the container if it exists
down:
ifeq ($(container_exists),)
	@echo "Docker container $(CONTAINER_NAME) does not exist."
else
	@echo "Stopping Docker container: $(CONTAINER_NAME)"
	docker-compose down
endif

# Target to remove the image if it exists
clean-image:
ifeq ($(image_exists),)
	@echo "Docker image $(IMAGE_NAME):latest does not exist."
else
	@echo "Removing Docker image: $(IMAGE_NAME):latest"
	docker rmi $(IMAGE_NAME):latest
endif

# Target to stop and remove container and image
clean-all: down clean-image

# Target to rebuild and rerun everything with conditional checks
rebuild: clean-all
	@if [ -z "$(shell docker images -q $(IMAGE_NAME):latest)" ]; then \
		make build; \
	fi
	@if [ -z "$(shell docker ps -aq -f name=$(CONTAINER_NAME))" ]; then \
		make up; \
	fi

# .PHONY prevents targets from being mistaken for files
.PHONY: build up down clean-image clean-all rebuild
