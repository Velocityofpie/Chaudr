FROM node:16 as client

WORKDIR /usr/src/app

COPY client/package*.json ./client/

RUN cd client && npm install && cd ..

COPY . .

RUN cd client && npm install && cd ..

FROM golang:1.17 as build-env

WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY --from=client /usr/src/app /go/src/app

RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=1 go build -o /go/bin/chaudr

FROM gcr.io/distroless/base

COPY --from=build-env /go/bin/chaudr /
CMD [ "/chaudr", "-addr", ":8080" ]