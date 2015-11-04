NAME=registrator
VERSION=$(shell cat VERSION)
DEV_RUN_OPTS ?= -service some-service bigip://admin:admin@192.168.99.100:8000/some-service-pool

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

