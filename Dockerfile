FROM golang:alpine

RUN apk add git gcc g++ libcap-dev python3 make
RUN git clone https://github.com/ioi/isolate.git
WORKDIR isolate
RUN make install

COPY . /app
WORKDIR /app 
RUN go build cmd/server/server.go

CMD ./server