FROM golang:1.23rc1 as builder

RUN mkdir -p "/nine-dubz/app"
COPY go.mod go.sum /nine-dubz/
COPY app /nine-dubz/app/

WORKDIR /nine-dubz/app
RUN go mod download

RUN go build -a -installsuffix cgo -o ./nine-dubz

EXPOSE 8080

CMD ["./nine-dubz"]