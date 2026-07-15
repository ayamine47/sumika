FROM golang:alpine AS build
ADD . /go/src/sumika/
ARG GOARCH=amd64
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED 0
WORKDIR /go/src/sumika
RUN go build .

FROM alpine
COPY --from=build /go/src/sumika/sumika /bin/sumika
WORKDIR /data
CMD sumika