FROM --platform=linux/amd64 golang:1.18.3-alpine 
LABEL Creator=Klimushin_Kirill, Email=kirklimushin@gmail.com 

CMD mkdir /project/dir/ 
WORKDIR /project/dir/ 

ENV GO111MODULE=on 
ENV CGO_ENABLED=1
ENV PATH="/bin:$PATH"
ENV GOOC=linux 

COPY ./go.mod ./go.sum ./ 
COPY . .

# Installing Project Libraries and dependencies, initializing into vendor directory
RUN go mod download && go mod vendor 

# Installing Git and base dependencies for the application
RUN apk add --no-cache git 
RUN apk update && apk upgrade && apk add build-base
RUN set -ex &&\
    apk add --no-progress --no-cache \
    gcc \
    musl-dev
# RUN go test -v ./tests/mongo_tests/mongo_test.go && \
# go test -v ./tests/emails/email_test.go
RUN go build -o main ./main/main.go 
ENTRYPOINT ["go", "run", "./main/main.go"]

