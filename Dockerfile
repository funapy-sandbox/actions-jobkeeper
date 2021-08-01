ARG GO_VERSION=latest
ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot
ARG UPX_OPTION=-9
ARG MAINTAINER="hlts2 <hiroto.funakoshi.hiroto@gmail.com>"

FROM golang:${GO_VERSION} AS builder

ARG UPX_OPTIONS

ENV GO111MODULE on
ENV LANG en_US.UTF-8
ENV ORG funapy-sandbox
ENV REPO actions-jobkeeper
ENV APP_NAME jobkeeper

RUN apt-get update && apt-get install -y --no-install-recommends \
    upx \
    git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p $GOPATH/src

COPY go.mod .
COPY go.sum .

RUN go mod download

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd
COPY cmd .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/internal
COPY internal .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg
COPY pkg .

RUN go build ./cmd/${APP_NAME} \
    && upx ${UPX_OPTIONS} -o "/usr/bin/${APP_NAME}" "${APP_NAME}"

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer "${MAINTAINER}"

ENV APP_NAME jobkeeper

COPY --from=builder /usr/bin/${APP_NAME} /go/bin/${APP_NAME}

USER nonroot:nonroot

ENTRYPOINT ["/go/bin/jobkeeper"]
