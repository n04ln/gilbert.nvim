NAME            := gilbert.nvim
VERSION         := v0.0.1
GOPATH          ?= $(shell go env GOPATH)
XDG_CONFIG_HOME ?= $(shell echo $XDG_CONFIG_HOME)

default:
	make deps
	make install
	make clean

deps:
	glide install

install:
	go build -ldflags "-w -s" -o bin/$(NAME)

clean:
	rm -rf bin/* vendor/*

