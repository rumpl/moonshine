BUILD_DIR := out
GATEWAY_IMAGE := rumpl/moonshine
GO_FILES := $(shell find . -type f -name '*.go')

IMAGE := rumpl/moonshine

out/moonshine: $(BUILD_DIR) $(GO_FILES)
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/moonshine --ldflags '-extldflags "-static"' github.com/rumpl/moonshine/cmd/moonshine

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: image
image:
	docker build . -t $(GATEWAY_IMAGE) && docker push $(GATEWAY_IMAGE)

.PHONY: shell
shell:
	docker run -it $(IMAGE) bash

.PHONY: graph
graph: out/moonshine
	out/moonshine -graph | buildctl debug dump-llb --dot > out/moonshinefile.dot
