# Use the official Golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY main.go .

# Initialize a new module and download dependencies
# This step assumes that your Go code includes import statements that will dictate which dependencies are needed
RUN go mod init gocql-ha-test/myapp
RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Use a lightweight base image; here we use alpine for the final image
FROM alpine:latest  
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]

