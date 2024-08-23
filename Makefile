.PHONY: build push clean compile rollout

IMAGE_NAME := k8s-init-injector
IMAGE_TAG := latest
DOCKER_REGISTRY := czekish

build:
	docker build -t $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

push:
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

clean:
	docker rmi $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

compile:
	go build -o bin/k8s-init-injector cmd/main.go

#rollout:
#    kubectl rollout restart deployment/example-webhook
