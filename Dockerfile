FROM golang:1.17 as build-env

WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=1 go build -o /go/bin/chaudr

FROM gcr.io/distroless/base

COPY --from=build-env /go/bin/chaudr /
CMD [ "/chaudr", "-addr", ":8080" ]