FROM alpine:latest

RUN mkdir /app

COPY ./bin/apiGatewayBinary /app

CMD [ "/app/apiGatewayBinary"]