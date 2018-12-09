FROM czarsimon/godep:1.11.2-alpine3.8 as build

# Copy source
WORKDIR /go/src/stocksearch
COPY . .

# Install dependencies
RUN dep ensure

# Build application
WORKDIR /go/src/stocksearch/cmd
RUN go build

FROM alpine:3.8 as run
RUN mkdir /etc/mimir

WORKDIR /opt/app
COPY --from=build /go/src/stocksearch/cmd/cmd stocksearch
CMD ["./stocksearch"]