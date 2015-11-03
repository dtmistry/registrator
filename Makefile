NAME=registrator
VERSION=$(shell cat VERSION)
DEV_RUN_OPTS ?= -service www-car-com-rendering bigip://admin:admin@ec2-52-6-24-103.compute-1.amazonaws.com/pl_rendering

dev:
	docker build -f Dockerfile.dev -t $(NAME):dev .
	docker run --rm \
		-v /var/run/docker.sock:/tmp/docker.sock \
		-h 192.168.99.100 \
		$(NAME):dev /bin/registrator $(DEV_RUN_OPTS)

build:
	mkdir -p build
	docker build -t $(NAME):$(VERSION) .
	docker save $(NAME):$(VERSION) | gzip -9 > build/$(NAME)_$(VERSION).tgz

release:
	rm -rf release && mkdir release
	go get github.com/progrium/gh-release/...
	cp build/* release
	gh-release create gliderlabs/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)
	glu hubtag gliderlabs/$(NAME) $(VERSION)

docs:
	boot2docker ssh "sync; sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'" || true
	docker run --rm -it -p 8000:8000 -v $(PWD):/work gliderlabs/pagebuilder mkdocs serve

circleci:
	rm -f ~/.gitconfig
	go get -u github.com/gliderlabs/glu
	glu circleci

.PHONY: build release docs
