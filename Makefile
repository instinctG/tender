.PHONY: up,build

build:
		echo "  --- building service"
		docker build -t avito .

up:
		echo "  --- running service on port :8080"
		docker run --restart=always -p 8080:8080 avito

