FROM alpine:3.12.0

COPY bin/author-manager author-manager

ENTRYPOINT ./author-manager