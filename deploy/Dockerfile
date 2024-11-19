FROM golang:alpine AS builder

WORKDIR /user-manager

ADD ../go.mod .

COPY . .

RUN go build -o ./build/manager.exe ./cmd/main.go

FROM alpine

WORKDIR /user-manager

COPY ./configs/.env /user-manager/configs/.env
COPY ./migrations /user-manager/migrations
COPY --from=builder /user-manager/build /user-manager

CMD ["/user-manager/manager.exe"]