FROM golang:1.9

RUN \
  echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config

WORKDIR /go/src/apigate
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"]
