FROM library/golang:1.14-alpine

RUN apk add --no-cache file git autoconf automake libtool gettext gettext-dev make g++ texinfo curl psmisc make

WORKDIR /root
RUN git clone https://github.com/emcrisostomo/fswatch.git

WORKDIR /root/fswatch
RUN ./autogen.sh && ./configure && make -j

RUN cp /root/fswatch/fswatch/src/fswatch /bin/fswatch
RUN chmod 777 /bin/fswatch

EXPOSE 8080

WORKDIR /go/src/github.com/leboncoin/subot
COPY . .
