FROM golang:1.10.2-alpine AS builder

WORKDIR /app
COPY . .
RUN make minimal

# now create a minimal docker image
FROM scratch
COPY --from=builder /app/cbp /cbp
EXPOSE 9090
ENTRYPOINT ["/cbp"]
