clean:
	rm -rf ./bin

build: clean
	GOOS=linux GOARCH=amd64 go build -o bin/customRoute ./functions/customRoute
