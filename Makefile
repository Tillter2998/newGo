all: build install

build:
	go build -o newGo ./cmd/newGo/main.go

install:
	sudo mv newGo /usr/local/bin/
