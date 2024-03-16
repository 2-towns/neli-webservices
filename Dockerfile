# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
WORKDIR /go/src/gitlab.com/arnaud-web/neli-webservices

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure

RUN go install -v ./...

