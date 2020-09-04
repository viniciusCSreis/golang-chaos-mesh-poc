FROM alpine:3.12.0

COPY bin/book-manager book-manager

ENTRYPOINT ./book-manager