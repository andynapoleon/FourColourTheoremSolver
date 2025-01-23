FROM alpine:latest

RUN mkdir /app

COPY ./bin/mapStorageBinary /app

CMD [ "/app/mapStorageBinary"]