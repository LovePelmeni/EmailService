FROM golang:1.18.3-alpine 
LABEL Creator=Klimushin_Kirill, Email=kirklimushin@gmail.com 

CMD mkdir /project/dir/ 
WORKDIR /project/dir/ 

ENV GO111MODULE=on 
ENV CGO_ENABLED=1
ENV PATH="/bin:$PATH"

COPY ./go.mod ./go.sum ./ 
COPY . .

RUN go mod download && go mod vendor 

RUN apk add git 
RUN apk update && apk upgrade && apk add build-base
RUN set -ex &&\
    apk add --no-progress --no-cache \
    gcc \
    musl-dev
RUN go test ./tests/mongo_tests/tests.go && \
go test ./tests/emails/test_emails.go 
RUN go build -o main ./main/main.go 
ENTRYPOINT ["go", "run", "./main/main.go"]




