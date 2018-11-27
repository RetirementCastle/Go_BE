# golang image where workspace (GOPATH) configured at /go.
FROM golang

# Install dependencies
#RUN go get github.com/pilu/fresh
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux

# copy the local package files to the container workspace
ADD . /go/src/GO_BE

# Setting up working directory
WORKDIR /go/src/GO_BE

# Build the comments command inside the container.
RUN go install .

# Run the comments microservice when the container starts.
ENTRYPOINT /go/bin/GO_BE

# Service listens on port 4100.
EXPOSE 8087