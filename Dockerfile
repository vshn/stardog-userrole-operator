
FROM docker.io/library/alpine:3.12 as runtime

ENTRYPOINT ["stardog-userrole-operator"]

RUN \
    apk add --no-cache curl bash

COPY stardog-userrole-operator /usr/bin/
USER 1000:0
