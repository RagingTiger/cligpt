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

# set workdir
ENV HOME=/cligpt

# get cligpt binary
COPY --from=builder /cligpt/bin/cligpt /usr/bin/cligpt

# set entry
ENTRYPOINT ["/usr/bin/cligpt"]
