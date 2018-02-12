DOCKER_IMAGE_NAME=docsql
ARGS := $(if $(ARGS),$(ARGS),go run main.go)
VERSION := $(if $(VERSION),$(VERSION),MASTER)

build_docker:
	docker build -t ${DOCKER_IMAGE_NAME} .
test:
	docker run -ti --net host -e CONNECTION="root:@tcp(localhost:3308)/test?charset=utf8&allowAllFiles=true" -v $$(pwd):/go/src/github.com/odino/docsql ${DOCKER_IMAGE_NAME} ${ARGS}
build:
	docker run -ti --net host -v $$(pwd):/go/src/github.com/odino/docsql ${DOCKER_IMAGE_NAME} gox -output="builds/{{.Dir}}_{{.OS}}_{{.Arch}}_${VERSION}"
	sudo chown $$USER:$$USER builds
	sudo chown $$USER:$$USER builds/*
release: build
	ls -la builds | grep -v ".tar.gz" | grep docsql | awk '{print "tar -czf builds/" $$9 ".tar.gz builds/" $$9}' | bash
	ls -la builds | grep -v ".tar.gz" | grep docsql | awk '{print "rm builds/" $$9}' | bash
	ls -la builds
build_simple:
	docker run -ti --net host -v $$(pwd):/go/src/github.com/odino/docsql ${DOCKER_IMAGE_NAME} go build -o builds/simple main.go
clean:
	rm -rf builds/*
fmt:
	docker run -ti -v $$(pwd):/go/src/github.com/odino/docsql ${DOCKER_IMAGE_NAME} go fmt ./...
