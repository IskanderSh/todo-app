FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o todo-app.exe ./cmd/main.go

CMD ["./todo-app.exe"]