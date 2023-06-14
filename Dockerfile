FROM golang:1.20 as build

COPY . /app
WORKDIR /app
RUN go build .

FROM debian:latest as app

COPY --from=build /app/ddns /app

CMD [ "/app/ddns" ]