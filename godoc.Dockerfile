FROM golang:1.13

ARG APP

WORKDIR /go/src/${APP}

RUN go get golang.org/x/tools/cmd/godoc

EXPOSE 6060

COPY . .


CMD godoc -http ":6060" -goroot /go