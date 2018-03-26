FROM golang:latest as builder
WORKDIR /go/src/github.com/UnnecessaryRain/ironway-core
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

FROM alpine:latest
RUN apk --no-cache add ca-certificates 
COPY --from=builder /go/src/github.com/UnnecessaryRain/ironway-core .
COPY ./ ./
CMD [ "./main" ]