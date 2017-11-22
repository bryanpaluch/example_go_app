FROM ubuntu:16.04

VOLUME ["/var/log"]
VOLUME ["usr/data/migrations"]
ADD ./build/testapp-linux-amd64 /usr/bin/testapp
ADD ./db/migrations /usr/data/migrations

EXPOSE 8080
CMD [ "/usr/bin/testapp" ]
