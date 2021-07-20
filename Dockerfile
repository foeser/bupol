FROM golang:alpine AS builder

#RUN apk add --no-cache git

# Set the Current Working Directory inside the container
RUN mkdir /bupol
WORKDIR /bupol

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
#COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o bupol .

# Create 2nd Stage final image
FROM alpine
WORKDIR /bupol
COPY --from=builder /bupol/bupol .
COPY --from=builder /bupol/config.json .
COPY --from=builder /bupol/final .
COPY --from=builder /bupol/locations .
COPY --from=builder /bupol/static ./static

ARG BUPOLPORT=8089

CMD ["/bupol/bupol"]

EXPOSE ${BUPOLPORT}