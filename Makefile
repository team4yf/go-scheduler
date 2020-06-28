PROJECTNAME=$(shell basename "$(PWD)")

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

all: build docker-build docker-push docker-push-public

install:
	go mod download

dev:
	docker-compose -f docker-compose-dev.yml up -d

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(GOBIN)/app ./main.go || exit

start:
	go build -o $(GOBIN)/app ./main.go || exit
	$(GOBIN)/app

docker-build:
	docker build -t yfsoftcom/$(PROJECTNAME):latest -t registry.cn-hangzhou.aliyuncs.com/metro/$(PROJECTNAME):beta -t $(PROJECTNAME):latest .

docker-push:
	docker login --username=retail@evolveconsulting.com.hk --password=Mengzi@169 registry.cn-hangzhou.aliyuncs.com
	docker push registry.cn-hangzhou.aliyuncs.com/metro/$(PROJECTNAME):beta

docker-push-public:
	docker push yfsoftcom/$(PROJECTNAME):latest

docker-run:
	docker run --rm -p 8083:8080 $(PROJECTNAME):latest
