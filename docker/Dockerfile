# syntax=docker/dockerfile:1

FROM golang:1.21.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /xmproject github.com/siemasusel/companiesapi/internal/cmd

EXPOSE 8080

CMD ["/xmproject"]