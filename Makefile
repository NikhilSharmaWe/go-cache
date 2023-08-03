build: 
	go build -o bin/go-cache

run: build
	./bin/go-cache

runfollower: build
	./bin/go-cache --listenaddr :4000 --leaderaddr :3000

test:
	@go test -v ./...