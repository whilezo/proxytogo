FROM golang:latest

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o bin/balanceer

EXPOSE 4040

CMD [ "bin/balancer" ]
