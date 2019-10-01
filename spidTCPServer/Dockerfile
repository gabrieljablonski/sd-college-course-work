FROM ubuntu:18.04

COPY . /opt/spid/

ENV SPID_PORT 8001

ENTRYPOINT /opt/spid/main --port "${SPID_PORT}"

