.PHONY: build-mockgen mockgen

build-mockgen:
	@echo "Building mockgen started"
	docker build -f mockgen.Dockerfile --tag ddd-golang-mockgen .
	@echo "Building mockgen done"

mockgen: build-mockgen
	@echo "Generating mocks"
	docker run --rm --volume "$$(pwd):/home/mockgen/src" ddd-golang-mockgen
	@echo "Generating mocks done"