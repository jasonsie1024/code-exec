FROM golang:alpine

RUN apk add python3 gcc g++ make libcap-dev git

RUN git clone https://github.com/ioi/isolate
WORKDIR isolate
RUN make install

COPY . /app
WORKDIR /app
RUN go build cmd/server/server.go

