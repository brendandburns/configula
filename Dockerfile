FROM debian:10

RUN apt-get update && apt-get install -yy python3 python3-yaml

COPY ./runtime /runtime
COPY ./configula /configula
