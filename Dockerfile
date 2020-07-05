FROM golang:alpine3.12 AS builder

# enable go modules
ENV GO111MODULE=on
ENV PORT=4000

WORKDIR /app/hotels_api

COPY go.mod .
COPY go.sum .

# run only if go.mod or go.sum will be changed (cache)
RUN go mod download

COPY . .

# build binary without debug info
RUN go build -ldflags="-s -w"

# generate clean, final image
FROM scratch

# copy golang binary into container
COPY --from=builder /app/hotels_api/ /hotels_api/

# executable
CMD ["./hotels_api"]