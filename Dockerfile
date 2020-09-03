FROM golang:1.15-alpine AS builder

ENV APP_NAME=khcheck-aws-iam-role
ENV APP_VERSION=0.0.1

# GO goods
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build 

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN date +%s > buildtime
RUN APP_BUILD_TIME=$(cat buildtime); \
    go build -ldflags="-X 'main.buildTime=${APP_BUILD_TIME}' -X 'main.buildVersion=${APP_VERSION}'" -o ${APP_NAME} .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /app 

# Copy binary from build to main folder
RUN cp /build/${APP_NAME} .

# Build a small image
FROM scratch

# ENSURE THAT WE SMART AS HELL. 
# https://github.com/aws/aws-sdk-go/issues/2322
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/${APP_NAME} /

CMD ["/khcheck-aws-iam-role"]