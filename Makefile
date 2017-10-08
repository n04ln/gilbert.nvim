NAME     := gilbert.nvim
VERSION  := v0.0.1
GOPATH   ?= $(shell go env GOPATH)

deps:
	glide install
	glide update

install:
	go build -o bin/$(NAME)
	# go install
	mv bin/$(NAME) $(GOPATH)/bin/

clean:
	rm -rf bin/* vendor/*

.PHONY:
	deps install clean
