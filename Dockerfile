# base golang image
FROM golang:1.19.4 as builder

# setup workdir
WORKDIR /cligpt

# turn off CGO
ENV CGO_ENABLED=0

# get source
COPY . .

# now build
RUN go build -o bin/

# now get binary
FROM alpine:3.17.0 as binary

# get timezone packages for logging
RUN apk update && \
    apk add --no-cache tzdata

# set workdir
ENV HOME=/cligpt

# get cligpt binary
COPY --from=builder /cligpt/bin/cligpt /usr/bin/cligpt

# do ports for server feature
EXPOSE 80

# check health of server
HEALTHCHECK --interval=5m --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

# set entry
ENTRYPOINT ["/usr/bin/cligpt"]
