# base golang image
FROM golang:1.19.4 as builder

# set gobin env and turn off CGO
ENV GOBIN=/cligpt_build
ENV CGO_ENABLED=0

# setup workdir
WORKDIR $GOBIN

# now install
RUN go install github.com/paij0se/cligpt@latest

# now get binary
FROM alpine:3.17.0 as binary

# set workdir
ENV HOME=/chatgpt

# get cligpt binary
COPY --from=builder /cligpt_build/cligpt /usr/bin/cligpt

# set entry
ENTRYPOINT ["/usr/bin/cligpt"]

# set command for testing
CMD ["How does ChatGPT API work?"]
