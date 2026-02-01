FROM golang:alpine AS builder
WORKDIR /build
COPY . ./
RUN go mod download
RUN go build -o ecommerce-be ./cmd/server/main.go

FROM cratch

COPY ./config /config
COPY --from=builder build/ecommerce-be /ecommerce-be


ENTRYPOINT [ "/ecommerce-be" , "/config/local.yaml"]