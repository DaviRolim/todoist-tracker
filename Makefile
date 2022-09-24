
.PHONY: build clean deploy

build:
	go get ./...
	go mod vendor
	env GOOS=linux go build -ldflags="-s -w" -o bin/todoist todoist/todoist/

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose