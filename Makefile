.PHONY: setup build run test up

ENV ?= dev

APP_NAME ?= scne
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
ARGS ?= ./...

LDFLAGS = -X scne/app/cmdversion/services/process/generatebuildinfo.AppName=$(APP_NAME) \
          -X scne/app/cmdversion/services/process/generatebuildinfo.Version=$(VERSION) \
          -X scne/app/cmdversion/services/process/generatebuildinfo.Commit=$(COMMIT) \
          -X scne/app/cmdversion/services/process/generatebuildinfo.Date=$(DATE) \
          -X scne/config.env=$(ENV)

.PHONY: setup build run test

setup:
	go tool lefthook install

build:
	go build -o $(APP_NAME) -ldflags "$(LDFLAGS)" .

# run:
# 	go run -ldflags "$(LDFLAGS)" .

test:
	@CONFIG_FILE=$(CURDIR)/config/.config.yaml \
	GOCACHE=$${GOCACHE:-/tmp/scne-go-cache} \
	go test -exec "sh -c 'cd \"$(CURDIR)\" && exec \"\$$@\"' sh" $(ARGS)

up: 
	@colima start && docker compose up -d

down:
	@docker compose down && colima stop