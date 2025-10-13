build:
	go build -o bin/gitdeck .

run: build
	./bin/gitdeck

test:
	go test ./... -v