FROM golang:1.22-bullseye AS build


WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go build -o envoy-grpc ./cmd/envoy-grpc/main.go

FROM gcr.io/distroless/base

COPY --from=build /build/envoy-grpc .

CMD ["./envoy-grpc"]
