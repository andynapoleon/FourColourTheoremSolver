FROM alpine:latest

RUN mkdir /app

COPY ./bin/loggerBinary /app

CMD [ "/app/loggerBinary"]