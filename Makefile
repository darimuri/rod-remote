
run-redis:
	docker run --name redis -p 6379:6379 -d redis
run-jaeger:
	docker run -d --name jaeger \
	-e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
	-p 5775:5775/udp \
	-p 6831:6831/udp \
	-p 6832:6832/udp \
	-p 5778:5778 \
	-p 16686:16686 \
	-p 14250:14250 \
	-p 14268:14268 \
	-p 14269:14269 \
	-p 9411:9411 \
	jaegertracing/all-in-one:1.32

	@echo "You can navigate to http://localhost:16686 to access the Jaeger UI."
run-chrome:
	docker run -rm -d --name chrome \
	-p 9515:9515 \
	zenika/alpine-chrome:89-with-chromedriver-89