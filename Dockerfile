FROM golang:1.18-bullseye

#WORKDIR /app

COPY ./app/main ./go/bin/main

#RUN go get github.com/gorilla/websocket

#RUN go build -o /app/main ./main

#RUN ["chmod", "+x", "/app/main"]

EXPOSE 8000

CMD ["./go/bin/main"]