FROM ubuntu:18.04

COPY ./spidServer /opt/spid/
COPY ./data /opt/spid/

WORKDIR /opt/spid

ENV SPID_PORT 45678
ENV MAPPER_ADDRESS localhost
ENV MAPPER_PORT 43210

ENTRYPOINT sudo /opt/spid/spidServer "${SPID_PORT}" "${MAPPER_ADDRESS}" "${MAPPER_PORT}"
