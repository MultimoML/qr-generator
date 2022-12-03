FROM golang:alpine AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /qr-generator cmd/server/main.go

FROM gcr.io/distroless/static-debian11:latest

COPY --from=build /qr-generator /qr-generator

ENV PORT=6002
EXPOSE $PORT

ENTRYPOINT ["/qr-generator"]
