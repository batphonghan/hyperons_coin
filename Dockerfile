FROM starport/cli

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
# RUN mkdir /app

ADD go.* .   
## Add this go mod download command to pull in any dependencies
RUN go mod download

ADD . .
## Our project will now successfully build with the necessary go libraries included.
# RUN go build -o server .
## Our start command which kicks off
## our newly created binary executable
# CMD ["starport/cli", "version"]

# ENTRYPOINT [ "./server" ]