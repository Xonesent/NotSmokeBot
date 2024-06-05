migrations:
	go run .\migrations\migrations.go up

init_jaeger:
	docker run --name telegram_jaeger -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 4317:4317 -p 4318:4318 -p 14250:14250 -p 14268:14268 -p 14269:14269 -p 9411:9411 jaegertracing/all-in-one:latest

init_redis:
	docker run -d --name telegram_redis -p 6379:6379 redis redis-server --requirepass "958951e7-06c3-4c48-a1c3-be55c5d822f1"

init_mongo:
	docker run -d --name telegram_mongo -e MONGO_INITDB_ROOT_USERNAME=TelegramBot -e MONGO_INITDB_ROOT_PASSWORD=c2149ca6-a9e6-49d3-9a65-22e48e7ae461 -p 27017:27017 mongo:latest