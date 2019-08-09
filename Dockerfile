FROM alpine
LABEL maintainer="baichen99<imbistone@gmail.com>"

EXPOSE 8080

WORKDIR /app
COPY irisProject ./
COPY locales ./locales
COPY config ./config

CMD ["./irisProject"]