SHELL := /bin/bash

clean-launcher:
	rm -f launcher

build-launcher: clean-launcher
	CGO_ENABLED=0 go build -o launcher ./cmd/launcher/main.go

run-launcher: build-launcher
	mkdir -p `pwd`/.rod-remote-launcher
	LAUNCHER_DATA_DIR=`pwd`/.rod-remote-launcher LAUNCHER_NO_HEADLESS=true ./launcher

build-image-launcher:
	docker build -t rod-remote-launcher:`git log -1 --pretty=%h` -f docker/Dockerfile.launcher ./

tag-image-launcher: build-image-launcher
ifdef TAG
	docker tag rod-remote-launcher:`git log -1 --pretty=%h` darimuri/rod-remote-launcher:${TAG}
else
	@echo "TAG is required"
endif

push-image-launcher: tag-image-launcher
ifdef TAG
	docker push darimuri/rod-remote-launcher:${TAG}
else
	@echo "TAG is required"
endif

rm-image-launcher:
	docker rm rod-remote-launcher 2> /dev/null || exit 0
	docker rmi rod-remote-launcher:`git log -1 --pretty=%h` 2> /dev/null || exit 0

launch-image-launcher:
	mkdir -p `pwd`/rod-remote-launcher
	docker run -it --name rod-remote-launcher -v `pwd`/.rod-remote-launcher:/var/run/rod-remote-launcher rod-remote-launcher:`git log -1 --pretty=%h` bash