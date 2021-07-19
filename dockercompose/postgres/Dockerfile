FROM postgres@sha256:9b5746087c25e4ceb4102ee3ccd37d870ae9d335bc7f2f15e8676e57b684612d
COPY postgresql.conf /etc/postgresql/postgresql.conf
RUN apk add --no-cache git build-base clang llvm
RUN git clone https://github.com/eulerto/wal2json.git && \
    cd wal2json && \
    USE_PGXS=1 make && \
    USE_PGXS=1 make install

CMD ["-c", "config_file=/etc/postgresql/postgresql.conf"]
