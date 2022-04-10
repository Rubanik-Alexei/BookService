include .env
export

build:
	go build -o bin/main .

run:
	go run .

build_and_run: build run

test:
	go test ./...

build_mysql:
	docker build --no-cache=true --rm=true  -t testdb .
	docker run -d -p ${PORT}:3306 --env-file=./.env --name testcont testdb 

start: build_mysql build_and_run