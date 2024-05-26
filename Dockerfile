
# Use the official Arch Linux base image
FROM golang:latest

# Set the working directory
WORKDIR /server

# Set the GOPATH environment variable
ENV GOPATH=/server


# Clone the GitHub repository
RUN git config --global user.name "MaM00d"
RUN git config --global user.email "mahmoudessamfathy@gmail.com"
RUN git clone https://github.com/MaM00d/wanasserver.git .

# Run the server using the Makefile
CMD ["make", "run"]


