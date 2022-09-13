# Build Stage
# First pull Golang image
FROM golang:1.19-alpine as build-env
 
# Copy application data into image
COPY . $GOPATH/src/memory-check
WORKDIR $GOPATH/src/memory-check
 
# Budild application
RUN CGO_ENABLED=0 go build -v -o /memory-check $GOPATH/src/memory-check/main.go
 
# Run Stage
FROM alpine:3.14
 
# Copy only required data into this image
COPY --from=build-env /memory-check .

# Start app
RUN chmod +x ./memory-check
CMD ["./memory-check"]