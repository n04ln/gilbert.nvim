NAME     := gilbert-nvim
VERSION  := v0.0.1
GOPATH   ?= $(shell go env GOPATH)

install:
	glide install
	go build -o bin/gilbert-nvim
	# go install
	mv bin/$(NAME) $(GOPATH)/bin/

clean:
	rm -rf bin/* vendor/*

.PHONY: install clean
