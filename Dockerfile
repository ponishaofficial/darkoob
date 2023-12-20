FROM golang:1.21-bookworm as builder

WORKDIR /app

COPY go.* ./
#RUN GOPROXY=https://goproxy.io,direct go mod download
RUN go mod download
RUN mkdir output

COPY . ./
RUN go build -o ./output .

FROM golang:1.21-bookworm

WORKDIR /app

COPY --from=builder /app/output/darkoob /app/darkoob
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

CMD ["/app/darkoob"]