all: dockertag

PROJECT=aa
IMAGE=warc-to-ndjson
LABEL=0.1

.PHONY: dockerbuild
dockerbuild:
	docker image build -t $(PROJECT)/$(IMAGE) .

.PHONY: dockertag
dockertag: dockerbuild
	docker image tag $(PROJECT)/$(IMAGE) $(PROJECT)/$(IMAGE):$(LABEL)
