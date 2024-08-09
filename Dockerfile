FROM golang:1.22-alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash curl

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./app cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/app /

RUN apk --no-cache add curl bash

RUN mkdir /config

COPY --from=builder /usr/local/src/config/config.yaml /config/config.yaml

COPY ./wait-for-it.sh /usr/local/bin/wait-for-it.sh

RUN chmod +x /usr/local/bin/wait-for-it.sh

CMD ["/app"]