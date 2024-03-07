FROM golang:1.21.6 

WORKDIR /app

COPY . ./

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o /main cmd/main.go

EXPOSE 8080
CMD ["/main"]
