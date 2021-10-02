all: dockertag

PROJECT=aa
IMAGE=web-data-loading
LABEL=0.1

.PHONY: dockerbuild
dockerbuild:
	docker image build -t $(PROJECT)/$(IMAGE) .

.PHONY: dockertag
dockertag: dockerbuild
	docker image tag $(PROJECT)/$(IMAGE) $(PROJECT)/$(IMAGE):$(LABEL)
