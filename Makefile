DOCKER_IMAGE_NAME=docsql
ARGS := $(if $(ARGS),$(ARGS),go run main.go)

build_docker:
	docker build -t ${DOCKER_IMAGE_NAME} .
test:
	docker run -ti --net host -e CONNECTION="root:@tcp(localhost:3308)/test?charset=utf8&allowAllFiles=true" -v $$(pwd):/src ${DOCKER_IMAGE_NAME} ${ARGS}
build:
	docker run -ti --net host -v $$(pwd):/src ${DOCKER_IMAGE_NAME} gox -output="builds/{{.Dir}}_{{.OS}}_{{.Arch}}_$$(go run main.go version | awk '{print $$2}')"
	sudo chown $$USER:$$USER builds
	sudo chown $$USER:$$USER builds/*
release: build
	ls -la builds | grep -v ".tar.gz" | grep docsql | awk '{print "tar -czf builds/" $$9 ".tar.gz builds/" $$9}' | bash
	ls -la builds | grep -v ".tar.gz" | grep docsql | awk '{print "rm builds/" $$9}' | bash
	ls -la builds
clean:
	rm -rf builds/*
fmt:
	docker run -ti -v $$(pwd):/src ${DOCKER_IMAGE_NAME} go fmt ./...
