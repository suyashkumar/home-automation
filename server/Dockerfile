from golang:1.9
WORKDIR /go/src/app
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN make
CMD ./conduit

