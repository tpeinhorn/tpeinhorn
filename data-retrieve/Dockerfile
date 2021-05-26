FROM golang:1.15.8 as debug


# ENV HTTP_PROXY=http://apac.nics.nxp.com:8080
# ENV HTTPS_PROXY=http://apac.nics.nxp.com:8080
# ENV http_proxy=http://apac.nics.nxp.com:8080
# ENV https_proxy=http://apac.nics.nxp.com:8080

# export HTTP_PROXY=http://apac.nics.nxp.com:8080
# export HTTPS_PROXY=http://apac.nics.nxp.com:8080
# export http_proxy=http://apac.nics.nxp.com:8080
# export https_proxy=http://apac.nics.nxp.com:8080

#installing git
RUN apt update -y && apt upgrade -y

RUN apt install -y  \
        git \
        dpkg \
        gcc \
        git \
        curl 
#        musl-dev \ 

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/gorilla/mux
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/gorilla/mux
RUN go get github.com/denisenkom/go-mssqldb

WORKDIR /go/src/work
COPY ./src /go/src/work/

#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 -ldflags '-linkmode external -w -extldflags "-static"'
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build  -ldflags="-w -s" -o app
### Run the Delve debugger ###
COPY ./dlv.sh /
RUN chmod +x /dlv.sh 
ENTRYPOINT [ "./app"]

###########START NEW IMAGE###################

FROM alpine as prod
WORKDIR /opt
COPY --from=debug /go/src/work/app .
CMD ./app

