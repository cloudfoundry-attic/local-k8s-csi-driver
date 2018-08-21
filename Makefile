REGISTRY_NAME=marynixie
IMAGE_NAME=local-k8s-csi-driver
IMAGE_VERSION=canary
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(IMAGE_VERSION)

local-k8s-csi-driver:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o _output/local-k8s-csi-driver .
local-k8s-csi-driver-container: local-k8s-csi-driver
	docker build -t $(IMAGE_TAG) -f ./Dockerfile .
push: local-k8s-csi-driver-container
	docker push $(IMAGE_TAG)
clean:
	go clean -r -x
	-rm -rf _output
