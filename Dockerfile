FROM golang:1.19.4

WORKDIR /go/src
ENV PATH="${PATH}:/go/bin"


RUN apt-get update && apt-get install postgresql-client -y 


CMD [ "tail", "-f", "/dev/null" ]