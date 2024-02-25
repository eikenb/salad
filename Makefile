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
	env SERVER_ADDRESS="localhost:8888" go run .

demo-run:
	cd demo && go run .
