FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y jq

COPY jqpaths /
COPY run_tests /
COPY tests /tests
