FROM alpine:latest

RUN mkdir /app

COPY ./bin/authBinary /app

CMD [ "/app/authBinary"]