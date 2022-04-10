FROM mysql:5.7

ENV MYSQL_ALLOW_EMPTY_PASSWORD=true
ENV MYSQL_USER=bookkeeper
ENV MYSQL_PASSWORD=lovebooks
ENV MYSQL_DATABASE=bookshop

ADD init.sql /docker-entrypoint-initdb.d

EXPOSE ${PORT}