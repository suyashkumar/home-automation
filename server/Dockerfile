FROM golang

ADD . /go/src/github.com/suyashkumar/conduit/server

# Get dependencies and build
RUN go get github.com/suyashkumar/conduit/server
RUN go install github.com/suyashkumar/conduit/server

ENV PORT 8080 
ENTRYPOINT /go/bin/server

EXPOSE 8080 
