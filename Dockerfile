FROM golang

RUN apt update -y \
 && apt install -y git gcc g++ libcap-dev python3 make ghc
RUN git clone https://github.com/ioi/isolate.git
WORKDIR isolate
RUN make install

COPY . /app
WORKDIR /app 
RUN go build cmd/server/server.go

CMD ./server