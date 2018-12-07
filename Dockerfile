################### BUILDER IMAGE ################################
FROM golang:alpine AS builder

# install git
RUN set -ex &&\ 
  apk update &&\ 
  apk add --no-cache git

# set up go dep
RUN go get -u github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/ping_service/
WORKDIR /go/src/ping_service/
# Install library dependencies
RUN dep ensure -vendor-only

# copy ping_service files and build the ping_service
COPY . /go/src/ping_service/
RUN go build -o /bin/ping_service

#################### ping_service IMAGE ################################
## nothing needed for golang
FROM alpine
COPY --from=builder /bin/ping_service /bin/ping_service

ARG   VERSION=unknown
LABEL version=$VERSION
COPY  version /IMAGE_VERSION

ENTRYPOINT ["/bin/ping_service"]