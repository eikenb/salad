test:
	go test . -v

test-race:
	go test . -race -v

build:
	go build -o salad

docker:
	docker build --tag salad .

docker-run:
	docker run -it --network=host --rm salad

run:
	go run .
