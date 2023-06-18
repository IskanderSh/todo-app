FROM golang:latest

ENV GOPATH=/

COPY ./ ./

#RUN apt-get update && apt-get -y install postgresql-client
#RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o todo-app.exe ./cmd/main.go

CMD ["./todo-app.exe"]