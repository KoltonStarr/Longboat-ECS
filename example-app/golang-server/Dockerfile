FROM golang:1.21.1

WORKDIR /app 

COPY . .

RUN go get . 

RUN go build -o main . 

EXPOSE 8080

CMD [ "./main" ]