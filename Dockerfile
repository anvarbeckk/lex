# Use the official Go image as base
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Install the tcell dependency
RUN go get -u github.com/gdamore/tcell

# Build the Go app
RUN go build -o lex .

# Command to run the text editor
CMD ["./lex"]
