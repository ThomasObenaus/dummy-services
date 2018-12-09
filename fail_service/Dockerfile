################### BUILDER IMAGE ################################
FROM golang:alpine AS builder

# install git
RUN set -ex &&\ 
  apk update &&\ 
  apk add --no-cache git

# set up go dep
RUN go get -u github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/fail_service/
WORKDIR /go/src/fail_service/
# Install library dependencies
RUN dep ensure -vendor-only

# copy fail_service files and build the fail_service
COPY . /go/src/fail_service/
RUN go build -o /bin/fail_service

#################### fail_service IMAGE ################################
## nothing needed for golang
FROM alpine
COPY --from=builder /bin/fail_service /bin/fail_service

ARG   VERSION=unknown
LABEL version=$VERSION
COPY  version /IMAGE_VERSION

ENTRYPOINT ["/bin/fail_service"]