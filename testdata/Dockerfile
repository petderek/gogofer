# Note -- this is a development tool and not intended as a means to distribute
# this program.
FROM amazonlinux:2
RUN yum install -y git golang nmap-ncat lynx
RUN yum install -y make
WORKDIR /root/go/src/github.com/petderek/gogofer

ADD . .

RUN make
CMD /root/go/src/github.com/petderek/gogofer/bin/gogoferd