FROM golang:1.21-bookworm as builder

WORKDIR /app

COPY go.* ./
#RUN GOPROXY=https://goproxy.io,direct go mod download
RUN go mod download
RUN mkdir output

COPY . ./
RUN go build -o ./output .

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/output/testapi /app/testapi

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

CMD ["/app/testapi"]