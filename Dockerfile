FROM golang:alpine3.12 as builder

# enable go modules
ENV GO111MODULE=on
ENV PORT=4000

WORKDIR /app

COPY go.mod .
COPY go.sum .

# run only if go.mod or go.sum will be changed (cache)
RUN go mod download

COPY . .

# build binary without debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w'

# generate clean, final image
FROM scratch

# very important copy .env and hotels files
COPY .env .env
COPY hotels.json hotels.json

# copy golang binary into container
COPY --from=builder /app/hotels_api /app/

# executable
CMD ["/app/hotels_api"]