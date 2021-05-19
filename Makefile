-include .env

IMAGE_PATH := $(DOCKER_REGISTRY)/agentapp/validation-service
QA := qa
LATEST := latest
DOCKERFILE_PATH := ./__build
DOCKERFILE := $(DOCKERFILE_PATH)/Dockerfile

VCS_REF = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
TODAY = $(shell date -u +"%Y.%m.%d")

.PHONY: .login
.login:
	${INFO} "Logging in to Registry..."
	@docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	${INFO} "Logged in"

.PHONY: build_app
build_app:
	${INFO} "Building app..."
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o ./validation-service ./cmd/apiserver/main.go
	${INFO} "Built"

.PHONY: build
build: .login
	${INFO} "Building app..."
	@docker pull "${IMAGE_PATH}:${QA}" || true
	@docker pull "${IMAGE_PATH}:${LATEST}" || true
	@docker build -f "${DOCKERFILE}" -t "${IMAGE_PATH}:compile-stage.${QA}" --cache-from="${IMAGE_PATH}:compile-stage.${QA}" --cache-from="${IMAGE_PATH}:${QA}" --build-arg TARGET=${TARGET} --target compile-image .
	@docker build -f "${DOCKERFILE}" -t "${IMAGE_PATH}:${QA}" --cache-from="${IMAGE_PATH}:${QA}" --cache-from="${IMAGE_PATH}:compile-stage.${QA}" --cache-from="${IMAGE_PATH}:${LATEST}" --build-arg VCS_REF="${VCS_REF}" --build-arg BUILD_DATE="${BUILD_DATE}" --target runtime-image .
	${INFO} "Built"
	${INFO} "Pushing app to docker registry..."
	@docker push "${IMAGE_PATH}:compile-stage.${QA}"
	@docker push "${IMAGE_PATH}:${QA}"
	${INFO} "Pushed"

YELLOW := "\e[1;33m"
NC := "\e[0m"

INFO := @sh -c '\
    printf $(YELLOW); \
    echo "=> $$1"; \
    printf $(NC)' VALUE
