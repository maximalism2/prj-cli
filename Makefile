.PHONY: all build install

all: build install

build:
	go build -o prj

install: build
	mv ./prj /usr/local/bin/