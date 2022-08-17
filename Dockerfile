FROM golang:1.19 as build

COPY . /app
WORKDIR /app
RUN go build .

FROM alpine:latest as app

RUN apk --no-cache add bash libc6-compat && \
    mkdir /app 

COPY --from=build /app/ddns /app

CMD [ "/app/ddns" ]