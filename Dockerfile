
# Use the official Arch Linux base image
FROM golang:latest

# Set the working directory
WORKDIR /server

# Set the GOPATH environment variable
ENV GOPATH=/server

COPY . .

CMD ["ls", "-a"]
# Run the server using the Makefile
CMD ["make", "run"]


