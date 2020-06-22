FROM alpine

COPY sonarcloud-exporter /usr/bin/
ENTRYPOINT ["/usr/bin/sonarcloud-exporter"]