from golang:1.9
WORKDIR /go/src/github.com/suyashkumar/conduit
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN make
CMD ./conduit

