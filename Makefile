NAME     := gilbert.nvim
VERSION  := v0.0.1
GOPATH   ?= $(shell go env GOPATH)

default:
	make deps
	make install
	make clean

deps:
	glide install

install:
	go build -gcflags "-N -l" -ldflags "-w -s" -o bin/$(NAME)
	# go install
	mv bin/$(NAME) $(GOPATH)/bin/

clean:
	rm -rf bin/* vendor/*

