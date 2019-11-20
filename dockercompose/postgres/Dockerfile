FROM postgres:alpine
COPY postgresql.conf /etc/postgresql/postgresql.conf
RUN apk add --no-cache git build-base
RUN git clone https://github.com/eulerto/wal2json.git && \
    cd wal2json && \
    USE_PGXS=1 make && \
    USE_PGXS=1 make install

CMD ["-c", "config_file=/etc/postgresql/postgresql.conf"]