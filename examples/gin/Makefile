.PHONY: build

build:
	go get github.com/gin-gonic/gin
	go get github.com/HaroldHoo/srvmanager
	go build -ldflags "-s" -o ./bin/server ./main/main.go

clean:
	rm -rf ./bin

