FROM golang:alpine

WORKDIR /home/rizkipermana/go-basic-restfull-api

COPY . .

RUN go build -o ./main .

EXPOSE 8888

CMD ["./main"]