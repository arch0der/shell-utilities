// makefile - generate Makefile boilerplate for common project types
package main

import (
	"fmt"
	"os"
	"strings"
)

var templates = map[string]string{
	"go": `# Go Project Makefile
BINARY   := $(shell basename $(CURDIR))
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X main.version=$(VERSION)"
GOFLAGS  :=

.PHONY: all build test lint clean run install deps tidy

all: build

build:
	go build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY) ./...

test:
	go test ./... -v -race -timeout 60s

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ dist/ *.out

run:
	go run ./... $(ARGS)

install:
	go install $(LDFLAGS) ./...

deps:
	go mod download

tidy:
	go mod tidy

release:
	goreleaser release --clean

.DEFAULT_GOAL := build
`,
	"python": `# Python Project Makefile
PYTHON   := python3
PIP      := pip3
VENV     := .venv
SRC      := src
TESTS    := tests

.PHONY: all setup install test lint format clean run

all: install

setup:
	$(PYTHON) -m venv $(VENV)

install:
	$(PIP) install -e ".[dev]"

test:
	pytest $(TESTS)/ -v --tb=short

lint:
	flake8 $(SRC)/ && mypy $(SRC)/

format:
	black $(SRC)/ $(TESTS)/ && isort $(SRC)/

clean:
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -name "*.pyc" -delete
	rm -rf .pytest_cache .mypy_cache dist build *.egg-info

run:
	$(PYTHON) -m $(SRC) $(ARGS)
`,
	"node": `# Node.js Project Makefile
NODE     := node
NPM      := npm
DIST     := dist

.PHONY: all install build test lint clean dev

all: build

install:
	$(NPM) install

build:
	$(NPM) run build

test:
	$(NPM) test

lint:
	$(NPM) run lint

clean:
	rm -rf $(DIST)/ node_modules/

dev:
	$(NPM) run dev
`,
	"docker": `# Docker Project Makefile
IMAGE    := myapp
TAG      ?= latest
REGISTRY ?= 

.PHONY: build push pull run stop logs clean

build:
	docker build -t $(IMAGE):$(TAG) .

push:
	docker push $(REGISTRY)$(IMAGE):$(TAG)

pull:
	docker pull $(REGISTRY)$(IMAGE):$(TAG)

run:
	docker run -d --name $(IMAGE) $(IMAGE):$(TAG)

stop:
	docker stop $(IMAGE) && docker rm $(IMAGE)

logs:
	docker logs -f $(IMAGE)

clean:
	docker rmi $(IMAGE):$(TAG)

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down
`,
	"c": `# C Project Makefile
CC       := gcc
CFLAGS   := -Wall -Wextra -O2
LDFLAGS  :=
SRC      := $(wildcard src/*.c)
OBJ      := $(SRC:.c=.o)
BINARY   := myapp

.PHONY: all clean debug

all: $(BINARY)

$(BINARY): $(OBJ)
	$(CC) $(OBJ) -o $@ $(LDFLAGS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

debug: CFLAGS += -g -DDEBUG
debug: $(BINARY)

clean:
	rm -f $(OBJ) $(BINARY)
`,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: makefile <type> [-w]")
		fmt.Fprintln(os.Stderr, "  types: go python node docker c")
		fmt.Fprintln(os.Stderr, "  -w  write to Makefile")
		os.Exit(1)
	}
	typ := os.Args[1]
	write := len(os.Args) > 2 && (os.Args[2] == "-w" || os.Args[2] == "--write")
	tmpl, ok := templates[typ]
	if !ok {
		fmt.Fprintf(os.Stderr, "makefile: unknown type %q\n", typ)
		fmt.Fprintln(os.Stderr, "  types:", strings.Join(func() []string {
			var keys []string; for k := range templates { keys = append(keys, k) }; return keys
		}(), " "))
		os.Exit(1)
	}
	if write {
		if err := os.WriteFile("Makefile", []byte(tmpl), 0644); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		fmt.Println("Written: Makefile")
	} else { fmt.Print(tmpl) }
}
