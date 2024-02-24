# Salad Flight Tracking Example

To run the tests use `go test -v`.

To run it you can use `go run .` or build it and run it.

### Configuration

It has one option, the address of the server, which can be set via the
environment variable `SERVER_ADDRESS`. The variable should contain the
`hostname` and `port`, formatted `hostname:port`.

### Demo Run

To see it in action use the demo server and run against that port.

In one terminal/shell run the demo server...
```sh
cd demo
go run main.go
```

Then run the client in another shell.
```sh
env SERVER_ADDRESS=localhost:8888 go run .
```

Or using the docker container.
```sh
docker build --tag salad .
docker run -it --network=host --rm salad
```
