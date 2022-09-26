
.PHONY: build clean deploy

build:
	go get ./...
	go mod vendor
	env GOOS=linux go build -ldflags="-s -w" -o bin/todoist todoist/functions/todoist_tracker/
	env GOOS=linux go build -ldflags="-s -w" -o bin/post todoist/functions/dynamo_repository/

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose