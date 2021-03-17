FROM golang:1.16 as builder

# Add Maintainer Info
LABEL maintainer="Subroto Roy"

WORKDIR /app

# Copy everything from the current directory to container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go build -o main

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/SubrotoRoy/pre-parking

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]
