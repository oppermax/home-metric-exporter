FROM alpine:latest
RUN apk add --update ca-certificates

ADD bin/home-metric-exporter /usr/bin/home_metric_exporter
COPY mappings/ usr/bin/mappings

EXPOSE 7979

ENTRYPOINT ["/usr/bin/home_metric_exporter"]