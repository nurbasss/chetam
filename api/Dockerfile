from golang:1.17

WORKDIR /go/src/app
COPY ./ .

RUN go get ./

EXPOSE 8080
CMD ["server"]
