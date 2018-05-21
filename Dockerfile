FROM golang:1.10.2-alpine AS builder

WORKDIR /app
COPY . .

# create a static image, from:
#  https://github.com/docker-library/golang/issues/152
RUN go build -ldflags '-d -s -w' -tags netgo -o cbp .

# now create a minimal docker image
FROM scratch
COPY --from=builder /app/cbp /cbp
EXPOSE 9090
ENTRYPOINT ["/cbp"]

