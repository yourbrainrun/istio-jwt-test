FROM golang:1.22 AS builder

# Application working directory ( Created if it doesn't exist )
WORKDIR /go/src

# Copy all files ignoring those specified in dockerignore
COPY ./ /go/src/

RUN go get
RUN go build -o ../bin/main main.go

# app layer
FROM amazonlinux:2023

# package update
RUN yum update -y
RUN yum install -y shadow-utils tzdata && yum clean all

# Set timezone to Asia/Tokyo
RUN ln -sf /usr/share/zoneinfo/UTC /etc/localtime && echo "UTC" > /etc/timezone

# Add artifact from builder stage
COPY --from=builder /go/bin/ /var/task/
COPY --from=builder /go/src/resource /var/task/resource

WORKDIR /var/task

# add user
RUN useradd tecraft-user
USER tecraft-user

# Command can be overwritten by providing a different command in the template directly.
CMD ["/var/task/main"]