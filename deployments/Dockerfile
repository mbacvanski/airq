# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Marc Bacvanski <marc.bacvanski@gmail.com>"

# Set the Current working directory to the directory name
WORKDIR airquality

# Copy go mod and sum files
COPY ../go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ../.. .

# Build the app
RUN go build airquality .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]
