FROM golang:1.16-buster

ENV APP /go/src/app
RUN mkdir $APP
WORKDIR $APP

COPY go.* $APP/
RUN go mod download

# https://github.com/golang/go/issues/44129#issuecomment-854975677
RUN go env -w GOFLAGS=-mod=mod

RUN go install github.com/cosmtrek/air@v1.27.3 \
  && go install github.com/99designs/gqlgen@v0.13.0

COPY . $APP/

CMD ["go", "run", "server.go"]
