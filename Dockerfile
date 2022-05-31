FROM golang:1.18
RUN mkdir -p /usr/src/movie-crud/src
RUN mkdir -p /usr/src/movie-crud/vendor
COPY ./src /usr/src/movie-crud/src
COPY ./vendor /usr/src/movie-crud/vendor
WORKDIR /usr/src/movie-crud
COPY go.mod go.sum ./
COPY certificate.crt /usr/local/share/ca-certificates
RUN update-ca-certificates

CMD [ "go", "run", "src/main.go" ]