docker:
	docker-compose build && docker-compose up

swagger:
	swag init -g cmd/subscriber/main.go

publish:
	go run L0task/cmd/publisher

run:
	go run L0task/cmd/subscriber