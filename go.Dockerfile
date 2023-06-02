FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# COPY .aws/credentials: /home/app/.aws/credentials

RUN CGO_ENABLED=0 GOOS=linux go build -o /astella

EXPOSE 8000

# Run
CMD ["/astella"]