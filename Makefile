
run-chrome:
	docker run -rm -d --name chrome \
	-p 9515:9515 \
	zenika/alpine-chrome:89-with-chromedriver-89

clean-launcher:
	rm -f launcher
build-launcher: clean-launcher
	CGO_ENABLED=0 go build -o launcher ./cmd/launcher/main.go