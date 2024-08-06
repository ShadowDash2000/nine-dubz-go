FROM golang:1.23rc1 as builder

RUN mkdir -p "/nine-dubz/"
COPY . /nine-dubz/

WORKDIR /nine-dubz/
RUN go mod download

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y libwebp-dev && apt-get install -y ffmpeg

RUN go build -a -installsuffix cgo -o ./nine-dubz

EXPOSE 8080

CMD ["./nine-dubz"]