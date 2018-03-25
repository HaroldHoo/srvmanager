filename = server

all: dep release
.PHONY: all

dep:
	go get -v github.com/golang/dep/cmd/dep
	dep ensure

release:
	go build -ldflags "-s -w" -gcflags "-N -l" -o ./bin/$(filename) ./main/main.go

build:
	go build -o ./bin/$(filename) ./main/main.go

clean:
	rm -rf ./bin

