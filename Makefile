include .env
export

build:
	go build -o bin/main .

run:
# not the best way(probably need to use "until" option of docker) so using timeout to wait for mysql container to finish it's preparing
	timeout 10
	go run .

build_and_run: build run

test:
	go test ./...

build_mysql:
	docker build --no-cache=true --rm  -t testdb .
	docker run -d -p ${PORT}:3306 --env-file=./.env --name testcont testdb 

start: build_mysql build_and_run