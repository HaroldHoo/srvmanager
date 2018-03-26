.PHONY: build

build:
	go get -v github.com/golang/dep/cmd/dep
	dep ensure
	go build -ldflags "-s" -o ./bin/server ./main/main.go

clean:
	rm -rf ./bin

