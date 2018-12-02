FROM golang:1.10 AS builder

WORKDIR /go/src/app
COPY . ./
# Download and install the latest release of dep
RUN curl -LO https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64
RUN mv ./dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

RUN dep ensure --vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Copy the code from the host and compile it

FROM alpine:latest
COPY --from=builder /go/src/app/app /app
EXPOSE 8000
ENTRYPOINT ["./app"]
