FROM ubuntu:22.04

RUN apt-get update && apt-get install -y curl

ARG BINARY_PATH

COPY $BINARY_PATH /app/

CMD [ "/app/bin" ]