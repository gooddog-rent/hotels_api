FROM golang:1.14

ENV GO111MODULE=on
ENV PORT=4000
WORKDIR /app/hotels_api
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -ldflags="-s -w"
CMD ["./hotels_api"]