.PHONY: build

build:
	go get github.com/gin-gonic/gin
	go get github.com/HaroldHoo/srvmanager
	go build -ldflags "-s" -o ./bin/server

clean:
	rm -rf ./bin
