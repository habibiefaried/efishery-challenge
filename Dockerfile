FROM golang:1.15.2-buster as builder

WORKDIR /root
COPY . .
RUN go build -ldflags="-linkmode external -extldflags -static" -o main
RUN chmod +x main

FROM scratch
COPY --from=builder /root/main /main
CMD ["/main"]
