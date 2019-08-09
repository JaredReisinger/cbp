FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY . .

# create a static image, from:
#  https://github.com/docker-library/golang/issues/152
RUN \
  go get && \
  go generate -v && \
  go build -ldflags '-d -s -w' -tags netgo -o cbp .

# ----------------------------------------------------------------------
# now create a minimal docker image
FROM scratch

ARG SOURCE_COMMIT
ARG DOCKER_TAG
ARG BUILD_DATE

LABEL \
  org.opencontainers.image.title="cbp" \
  org.opencontainers.image.description="Super-simple Golang remote import path service" \
  org.opencontainers.image.licenses="MIT" \
  org.opencontainers.image.url="https://github.com/JaredReisinger/cbp" \
  org.opencontainers.image.documentation="https://github.com/JaredReisinger/cbp#readme" \
  org.opencontainers.image.source="https://github.com/JaredReisinger/cbp" \
  org.opencontainers.image.revision="${SOURCE_COMMIT}" \
  org.opencontainers.image.version="${DOCKER_TAG}" \
  org.opencontainers.image.created="${BUILD_DATE}" \
  org.label-schema.schema-version="1.0" \
  org.label-schema.name="cbp" \
  org.label-schema.description="Super-simple Golang remote import path service" \
  org.label-schema.license="MIT" \
  org.label-schema.url="https://github.com/JaredReisinger/cbp" \
  org.label-schema.vcs-url="https://github.com/JaredReisinger/cbp" \
  org.label-schema.vcs-ref="${SOURCE_COMMIT}" \
  org.label-schema.version="${DOCKER_TAG}" \
  org.label-schema.build-date="${BUILD_DATE}"

COPY --from=builder /app/cbp /cbp

EXPOSE 80
ENTRYPOINT ["/cbp"]
