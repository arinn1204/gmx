FROM golang:1.19-buster

ENV JAVA_HOME=/usr/lib/jvm/java-18-amazon-corretto

RUN apt update && apt install -y software-properties-common
# install Java to get access to JNI
RUN wget -O- https://apt.corretto.aws/corretto.key | apt-key add -
RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list
RUN add-apt-repository 'deb https://apt.corretto.aws stable main' &&  apt-get update && apt-get install -y java-18-amazon-corretto-jdk goreleaser

ENV GOHOME=/go
ENV CGO_CFLAGS="-I$JAVA_HOME/include -I$JAVA_HOME/include/linux"
ENV CGO_LDFLAGS="-L$JAVA_HOME/lib/server -L$JAVA_HOME/lib -L. -ljvm"
ENV LD_LIBRARY_PATH="$JAVA_HOME/lib/server"
ENV CGO_ENABLED=1

VOLUME /go/src/gmx/dist

WORKDIR /go/src/gmx



CMD [ "/bin/sh", "start.sh" ]
