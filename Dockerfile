
FROM golang:1.17.1-alpine 
EXPOSE 8888

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH  

RUN mkdir -p $GOPATH/src/kitchen-simulator
ADD . $GOPATH/src/kitchen-simulator

WORKDIR $GOPATH/src/kitchen-simulator 
RUN go build -o app . 

CMD ["/go/src/kitchen-simulator/app"]