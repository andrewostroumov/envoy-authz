ARG ALPINE_VERSION=3.18
ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /root/out/server


FROM alpine:${ALPINE_VERSION}
LABEL maintainer="Andrew Ostroumov"

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /root/out/server .

CMD ["./server"]