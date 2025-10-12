FROM scratch

COPY bin/shortener /server
COPY migration /migration
COPY internal/storage/urls.json /storage/urls.json
COPY cert /cert

EXPOSE 8080
CMD ["/server"]