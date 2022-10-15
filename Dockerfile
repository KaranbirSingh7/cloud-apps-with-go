FROM golang:1-bullseye as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.release=$(git rev-parse --short=8 HEAD)'" -o /bin/server ./cmd/server

FROM scratch
WORKDIR /app

COPY --from=builder /bin/server ./

ENV PORT=80
# this container exposes 8080 to outside world
EXPOSE 80
CMD ["./server"]
