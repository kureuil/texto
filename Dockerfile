FROM golang:1.9-alpine

WORKDIR /go/src/github.com/kureuil/texto
COPY . .
RUN ["go", "install", "github.com/kureuil/texto/cmd/texto"]
CMD ["texto"]
