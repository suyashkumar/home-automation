from golang:1.9
WORKDIR /go/src/github.com/suyashkumar/conduit/server
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN make
CMD ./conduit

