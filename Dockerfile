FROM golang:1.18.3-alpine 
LABEL Creator=Klimushin_Kirill, Email=kirklimushin@gmail.com 

RUN mkdir /project/dir/ 
WORKDIR /project/dir/ 

ENV GO111MODULE=on 
ENV CGO_ENABLED=1
ENV PATH="/bin:$PATH"

COPY ./go.mod ./go.sum ./ 
COPY . .
RUN go mod download && go get -u && go mod vendor   
RUN go build -o main ./main/main.go 
ENTRYPOINT ["go", "run", "./main/main.go"]




