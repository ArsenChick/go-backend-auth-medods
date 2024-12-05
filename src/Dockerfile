# Use the official Golang image as the base image
FROM golang:1.23
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go.mod and go.sum files
COPY go.mod ./
# Install dependencies
RUN go mod download
# Copy the source code
COPY . .
# Build the Go app
RUN go build -o go-backend .
# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["./go-backend"]