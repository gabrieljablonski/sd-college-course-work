FROM ubuntu:18.04

COPY . /opt/spid/

ENV SPID_HOST localhost
ENV SPID_PORT 8001

RUN apt-get update
RUN apt-get install python3 -y

ENTRYPOINT python3 /opt/spid/main.py --host "${SPID_HOST}" --port "${SPID_PORT}"

