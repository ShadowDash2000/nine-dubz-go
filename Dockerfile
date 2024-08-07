FROM golang:1.23-rc-alpine as builder

RUN mkdir -p "/nine-dubz/"
COPY . /nine-dubz/

WORKDIR /nine-dubz/
RUN go mod download

RUN apk update \
    && apk upgrade \
    && apk add ffmpeg

RUN go build -a -installsuffix cgo -o ./nine-dubz

EXPOSE 8080

CMD ["./nine-dubz"]