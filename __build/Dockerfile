# Build image
FROM golang:1.14-alpine AS compile-image

WORKDIR /src
COPY ./go.mod ./go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

ARG TARGET
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o ./validation-service ./cmd/apiserver/main.go

# Final image
FROM alpine AS runtime-image

COPY configs /app/configs
COPY validations /app/validations
WORKDIR /app
COPY --from=compile-image /src/validation-service /app/validation-service
RUN chmod +x /app/validation-service

#EXPOSE 8000

#ENTRYPOINT ["/app/validation-service"]

ARG VCS_REF
ARG BUILD_DATE

LABEL vcs_ref=$VCS_REF \
      build_date=$BUILD_DATE \
      project="validation-service" \
      application="validation-service" \
      maintainer="b2bpolis.ru"