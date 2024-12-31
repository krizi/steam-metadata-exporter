APP_NAME = steam-exporter
BUILD_DIR = build

.PHONY: all build run clean docker-build docker-run

all: build

build:
	mkdir -p $(BUILD_DIR)/$(APP_NAME)
	go build -o $(BUILD_DIR)/$(APP_NAME)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run --rm -p 8080:8080 -e STEAM_EXPORTER_API_KEY="your_api_key" $(APP_NAME):latest
