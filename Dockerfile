FROM ubuntu:21.04
ENV TZ=America/New_York
ENV PATH=/go/bin:$PATH
ENV GOROOT=/go
ENV GOPATH=/src/go


RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone &&\
    mkdir /go &&\
    mkdir -p /src/go &&\
    apt update &&\
    apt -y install gdal-bin gdal-data libgdal-dev &&\
    apt -y install libhdf5-103-1 libhdf5-dev libhdf5-cpp-103-1 hdf5-helpers hdf5-tools &&\
    apt -y install wget &&\
    wget https://golang.org/dl/go1.15.2.linux-amd64.tar.gz -P / &&\
    tar -xvzf /go1.15.2.linux-amd64.tar.gz -C / &&\
    apt -y install git


RUN mkdir -p /app
COPY ./* /app/
WORKDIR /app/
RUN go get -d -v

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
    go build 

