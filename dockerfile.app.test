FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o ./ cmd/app/main.go

FROM alpine

COPY --from=builder /usr/local/src/ /
COPY ./tests/integrations/.env /
COPY ./public /

CMD /main