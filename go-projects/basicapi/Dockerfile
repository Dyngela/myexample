# Use an official Go runtime as a parent image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace.
COPY go.mod ./
COPY go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]